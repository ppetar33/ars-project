package poststore

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"os"
)

type PostStore struct {
	cli *api.Client
}

func New() (*PostStore, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &PostStore{
		cli: client,
	}, nil
}

func (ps *PostStore) Get(id string) (*Service, error) {
	kv := ps.cli.KV()

	pair, _, err := kv.Get(constructKey(id), nil)
	if err != nil {
		return nil, err
	}

	post := &Service{}
	err = json.Unmarshal(pair.Value, post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (ps *PostStore) GetAll() ([]*Service, error) {
	kv := ps.cli.KV()
	data, _, err := kv.List(all, nil)
	if err != nil {
		return nil, err
	}

	posts := []*Service{}
	for _, pair := range data {
		post := &Service{}
		err = json.Unmarshal(pair.Value, post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (ps *PostStore) Delete(id string) (map[string]string, error) {
	kv := ps.cli.KV()
	_, err := kv.Delete(constructKey(id), nil)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}

func (ps *PostStore) Post(post *Service) (*Service, error) {
	kv := ps.cli.KV()

	sid, rid := generateKey()
	post.Id = rid

	data, err := json.Marshal(post)
	if err != nil {
		return nil, err
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}

	return post, nil
}
