package poststore

import (
	"fmt"
	"github.com/google/uuid"
)

const (
	all             = "conf/"
	conf            = "conf/%s/%s"
	confId          = "conf/%s/"
	confIdempotency = "conf/%s/%s/%s"
)

func generateKey(ver string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(conf, id, ver), id
}

func generateKeyWithIdempotency(ver string, idempotency string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(confIdempotency, id, ver, idempotency), id
}

func constructKeyUpdatey(id string, ver string, idempotency string) string {
	return fmt.Sprintf(confIdempotency, id, ver, idempotency)
}

func constructKey(id string, version string) string {
	return fmt.Sprintf(conf, id, version)
}

func constructConfigIdKey(id string) string {
	return fmt.Sprintf(confId, id)
}
