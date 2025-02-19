package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	DEBUG                             = true
	QUOTE_API_URL                     = "http://localhost:8085/cotacao"
	QUOTE_TIMEOUT_DEFAULT_API_REQUEST = 300
	QUOTE_FILE                        = "cotacao.txt"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	port := os.Getenv("SERVER_PORT")
	log.Print("Server starting on port " + port)

	http.HandleFunc("/", Handler)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

type Quote struct {
	Bid string `json:"bid"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	quote, err := getDataApiHandler(r.Context())
	if err != nil {
		if err == context.DeadlineExceeded {
			http.Error(w, "", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	err = saveQuoteTxt(quote)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	log.Println("[INFO] requisição processada com sucesso")
	w.Write([]byte(""))
}

func getDataApiHandler(requestContext context.Context) (*Quote, error) {
	timeout, err := time.ParseDuration(os.Getenv("QUOTE_TIMEOUT_REQUEST"))
	if err != nil {
		timeout = QUOTE_TIMEOUT_DEFAULT_API_REQUEST * time.Millisecond
	}

	timer := time.NewTimer(timeout)
	var quote Quote
	apiChan := make(chan *http.Response, 1)
	errChan := make(chan error, 1)
	go func(quote *Quote) {
		delay, err := time.ParseDuration(os.Getenv("QUOTE_REQUEST_DELAY"))
		if err != nil {
			delay = 0 * time.Second
		}

		if DEBUG {
			println("getDataApiHandler() #1")
		}

		if requestContext.Err() == context.Canceled {
			return
		}

		log.Println("[INFO] API request delay: ", delay)
		if delay > 0 {
			time.Sleep(delay)
		}

		if requestContext.Err() == context.Canceled {
			return
		}

		if DEBUG {
			println("getDataApiHandler() #2")
		}

		client := &http.Client{}
		resp, err := client.Get(QUOTE_API_URL)
		if err != nil {
			if DEBUG {
				println("getDataApiHandler() #3")
			}
			errChan <- err
			return
		}

		if resp.StatusCode != http.StatusOK {
			if DEBUG {
				println("getDataApiHandler() #4")
			}
			errChan <- fmt.Errorf("error status code %v", resp.StatusCode)
			return
		}

		if DEBUG {
			println("getDataApiHandler() #5")
		}

		defer resp.Body.Close()

		if requestContext.Err() == context.Canceled {
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			if DEBUG {
				println("getDataApiHandler() #6")
			}
			errChan <- err
			return
		}

		if DEBUG {
			println("getDataApiHandler() #7")
		}

		err = json.Unmarshal([]byte(body), quote)
		if err != nil {
			if DEBUG {
				println("getDataApiHandler() #8")
			}
			errChan <- err
			return
		}

		if DEBUG {
			println("getDataApiHandler() #9")
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

		if requestContext.Err() == context.DeadlineExceeded {
			log.Println("[ERRO] timeout na requisição")
		}
		return &quote, err
	case err := <-errChan:
		log.Println("[ERRO] na requisição: ", err)
		return &quote, err
	case <-apiChan:
		log.Println("[INFO] requisição a API realizada com sucesso.")
		return &quote, nil
	case <-timer.C:
		log.Println("[ERRO] timeout na requisição original")
		return &quote, context.DeadlineExceeded
	}
}

func saveQuoteTxt(quote *Quote) error {
	if DEBUG {
		println("saveQuoteTxt() #1")
	}

	f, err := os.OpenFile(QUOTE_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if DEBUG {
			println("saveQuoteTxt() #2")
		}
		log.Println("[Erro] ao criar/abrir o arquivo: ", err.Error())
		return err
	}

	if DEBUG {
		println("saveQuoteTxt() #3")
	}

	defer f.Close()
	msg := fmt.Sprintf("Dólar:%s\n", quote.Bid)
	_, err = f.WriteString(msg)

	if DEBUG {
		println("saveQuoteTxt() #4")
	}

	if err != nil {
		if DEBUG {
			println("saveQuoteTxt() #5")
		}
		log.Println("[Erro] ao escrever no arquivo: ", err.Error())
		return err
	}
	return nil
}
