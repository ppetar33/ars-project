package main

type Config struct {
	Entries map[string]string `json:"entries"`
}

type Service struct {
	Id   string             `json:"id"`
	Data map[string]*Config `json:"data"`
}
