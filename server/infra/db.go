package infra

import (
	"context"
	"d1-server/entity"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func DbSaveQuote(ctx context.Context, quote entity.Quote) error {
	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		debug = false
	}

	delay, err := time.ParseDuration(os.Getenv("DB_QUOTE_DELAY"))
	if err != nil {
		delay = 0 * time.Second
	}

	if debug {
		log.Println("[DEBUG] dbSaveQuote! #0")
	}

	log.Println("[INFO] request DB delay: ", delay)
	if delay > 0 {
		time.Sleep(delay)
	}

	if debug {
		log.Println("[DEBUG] dbSaveQuote! #1")
	}

	db, err := dbGetConectionSqlite()
	if err != nil {
		if debug {
			log.Println("[DEBUG] dbSaveQuote! #1.1")
		}
		return err
	}
	defer db.Close()

	if debug {
		log.Println("[DEBUG] dbSaveQuote! #2")
	}

	jsonData, err := json.Marshal(quote)
	if err != nil {
		return err
	}

	if debug {
		log.Println("[DEBUG] dbSaveQuote! #3")
	}

	var id int64
	dml := `INSERT INTO quotes (version, json, created_date) VALUES (?, ?, ?) RETURNING id`
	err = db.QueryRowContext(ctx, dml, "1", jsonData, time.Now().Format("2006-01-02 15:04:05")).Scan(&id)
	if err != nil {
		if debug {
			log.Println("[DEBUG] dbSaveQuote! #3.1")
		}
		return err
	}

	if debug {
		log.Println("[DEBUG] dbSaveQuote! #4")
	}

	if id == 0 {
		return fmt.Errorf("invalid id")
	}

	if debug {
		log.Println("[DEBUG] dbSaveQuote! #5")
	}

	return nil
}

func dbGetConectionSqlite() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", os.Getenv("DB_SQLITE"))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		defer db.Close()
	}

	return db, err
}
