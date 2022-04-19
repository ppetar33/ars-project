package main

type Config struct {
	entries map[string]string
}

type Service struct {
	data map[string]*Config
}
