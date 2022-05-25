package poststore

type Config struct {
	Label   map[string]string `json:"label"`
	Entries map[string]string `json:"entries"`
}

type Service struct {
	Id      string   `json:"id"`
	Version string   `json:"version"`
	Data    []Config `json:"data"`
}
