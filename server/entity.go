package main

type QuoteToUSDBRL struct {
	USDBRL Quote `json:"USDBRL"`
}

type Quote struct {
	Code        string `json:"code"`
	Codein      string `json:"codein"`
	Name        string `json:"name"`
	High        string `json:"high"`
	Low         string `json:"low"`
	VarBid      string `json:"varBid"`
	PctChange   string `json:"pctChange"`
	Bid         string `json:"bid"`
	Ask         string `json:"ask"`
	Timestamp   string `json:"timestamp"`
	Create_date string `json:"create_date"`
}

type QuoteDB struct {
	Id          int    `json:"id"`
	Version     string `json:"version"`
	Json        string `json:"data"`
	Create_date string `json:"create_date"`
}

type QuoteBid struct {
	Bid string `json:"bid"`
}
