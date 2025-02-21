package handler

import (
	"context"
	"d1-server/entity"
	"d1-server/infra"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	QUOTE_API_URL                     = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	QUOTE_TIMEOUT_DEFAULT_DB          = 10
	QUOTE_TIMEOUT_DEFAULT_API_REQUEST = 200
)

func QuoteHandler(w http.ResponseWriter, r *http.Request) {
	quote, err := getDataApiHandler(r.Context())
	if err != nil {
		if err == context.DeadlineExceeded {
			http.Error(w, "", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	err = saveQuoteHandler(quote, r.Context())
	if err != nil {
		if err == context.DeadlineExceeded {
			http.Error(w, "", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	quoteBid := entity.QuoteBid{
		Bid: quote.USDBRL.Bid,
	}

	jsonData, err := json.Marshal(quoteBid)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	log.Println("[INFO] Requisição processada com sucesso", quoteBid)
	w.Write([]byte(jsonData))
}

func getDataApiHandler(requestContext context.Context) (*entity.QuoteToUSDBRL, error) {
	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		debug = false
	}

	var quote entity.QuoteToUSDBRL
	timeout, err := time.ParseDuration(os.Getenv("QUOTE_TIMEOUT_REQUEST"))
	if err != nil {
		timeout = QUOTE_TIMEOUT_DEFAULT_API_REQUEST * time.Millisecond
	}

	log.Println("[INFO] API request timeout: ", timeout)
	apiCtx, apiCancel := context.WithTimeout(context.Background(), timeout)
	defer apiCancel()

	req, err := http.NewRequestWithContext(apiCtx, "GET", QUOTE_API_URL, nil)
	if err != nil {
		log.Println("[ERROR] ao criar a requisição:", err)
		return &quote, err
	}

	apiChan := make(chan *http.Response, 1)
	errChan := make(chan error, 1)
	go func(quote *entity.QuoteToUSDBRL) {
		delay, err := time.ParseDuration(os.Getenv("QUOTE_REQUEST_DELAY"))
		if err != nil {
			delay = 0 * time.Second
		}

		if debug {
			log.Println("[DEBUG] getDataApiHandler! #0")
		}

		log.Println("[INFO] delay to API request: ", delay)
		if delay > 0 {
			time.Sleep(delay)
		}

		if debug {
			log.Println("[DEBUG] getDataApiHandler! #1")
		}

		if apiCtx.Err() == context.Canceled {
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if debug {
			log.Println("[DEBUG] getDataApiHandler! #2")
		}

		if err != nil {
			if debug {
				log.Println("[DEBUG] getDataApiHandler! #3")
			}
			errChan <- err
			return
		}

		if debug {
			log.Println("[DEBUG] getDataApiHandler! #4")
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			if debug {
				log.Println("[DEBUG] getDataApiHandler! #5")
			}
			errChan <- err
			return
		}

		if debug {
			log.Println("[DEBUG] getDataApiHandler! #6")
		}
		err = json.Unmarshal([]byte(body), quote)
		if err != nil {
			if debug {
				log.Println("[DEBUG] getDataApiHandler! #7")
			}
			errChan <- err
			return
		}
		if debug {
			log.Println("[DEBUG] getDataApiHandler! #8")
		}
		apiChan <- resp
	}(&quote)

	select {
	case <-requestContext.Done():
		err := requestContext.Err()
		if requestContext.Err() == context.Canceled {
			log.Println("[INFO] requisição principal cancelada")
			err = fmt.Errorf("request canceled")
		}
		return &quote, err
	case <-apiCtx.Done():
		if apiCtx.Err() == context.DeadlineExceeded {
			log.Println("[ERRO] timeout na requisição: ", timeout)
		}
		return &quote, apiCtx.Err()
	case err := <-errChan:
		log.Println("[ERRO] na requisição: ", err.Error())
		return &quote, err
	case <-apiChan:
		log.Println("[INFO] requisição a API realizada com sucesso.")
		return &quote, nil
	}
}

func saveQuoteHandler(quote *entity.QuoteToUSDBRL, requestContext context.Context) error {
	timeout, err := time.ParseDuration(os.Getenv("DB_QUOTE_TIMEOUT"))
	if err != nil {
		timeout = QUOTE_TIMEOUT_DEFAULT_DB * time.Millisecond
	}

	dbCtx, dbCancel := context.WithTimeout(context.Background(), timeout)
	defer dbCancel()

	dbChan := make(chan bool, 1)
	errChan := make(chan error, 1)
	go func() {
		err = infra.DbSaveQuote(dbCtx, quote.USDBRL)
		if err != nil {
			errChan <- err
			return
		}
		dbChan <- true
	}()

	select {
	case <-requestContext.Done():
		err := requestContext.Err()
		if requestContext.Err() == context.Canceled {
			log.Println("[INFO] requisição principal cancelada")
			err = fmt.Errorf("request canceled")
		}
		return err
	case <-dbCtx.Done():
		if dbCtx.Err() == context.DeadlineExceeded {
			log.Println("[ERRO] timeout na operação db: ", timeout)
		}
		return dbCtx.Err()
	case err := <-errChan:
		log.Println("[ERRO] ao gravar no banco de dados: ", err)
		return err
	case <-dbChan:
		log.Println("[INFO] registro gravado com sucesso")
		return nil
	}
}
