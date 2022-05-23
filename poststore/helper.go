package poststore

import (
	"fmt"
	"github.com/google/uuid"
)

const (
	all    = "conf/"
	conf   = "conf/%s/%s/"
	confId = "conf/%s/"
)

func generateKey(ver string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(conf, id, ver), id
}

func constructKey(id string, version string) string {
	return fmt.Sprintf(conf, id, version)
}

func constructConfigIdKey(id string) string {
	return fmt.Sprintf(confId, id)
}
