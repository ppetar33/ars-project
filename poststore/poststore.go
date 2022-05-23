package poststore

import (
	"encoding/json"
	"errors"
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

func (ps *PostStore) Get(id string, version string) (*Service, error) {
	kv := ps.cli.KV()
	key := constructKey(id, version)
	data, _, err := kv.Get(key, nil)

	if err != nil || data == nil {
		return nil, errors.New("no data")
	}

	service := &Service{}
	err = json.Unmarshal(data.Value, service)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (ps *PostStore) GetAll() ([]*Service, error) {
	kv := ps.cli.KV()
	data, _, err := kv.List(all, nil)
	if err != nil {
		return nil, err
	}

	posts := []*Service{}
	for _, pair := range data {
		fmt.Println(pair)
		post := &Service{}
		err = json.Unmarshal(pair.Value, post)

		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (ps *PostStore) Delete(id string, version string) (*Service, error) {
	kv := ps.cli.KV()
	key := constructKey(id, version)
	data, _, err := kv.Get(key, nil)

	if err != nil || data == nil {
		return nil, errors.New("no data")
	}

	_, errDelete := kv.Delete(key, nil)
	if errDelete != nil {
		return nil, err
	}

	service := &Service{}
	err = json.Unmarshal(data.Value, service)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (ps *PostStore) Post(post *Service) (*Service, error) {
	kv := ps.cli.KV()

	sid, rid := generateKey(post.Version)
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

func (ps *PostStore) Update(service *Service) (*Service, error) {
	kv := ps.cli.KV()

	data, err := json.Marshal(service)
	if err != nil {
		return nil, err
	}

	c := &api.KVPair{Key: constructKey(service.Id, service.Version), Value: data}
	_, err = kv.Put(c, nil)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (ps *PostStore) FindConfVersions(id string) ([]*Service, error) {
	kv := ps.cli.KV()

	key := constructConfigIdKey(id)
	data, _, err := kv.List(key, nil)
	if err != nil {
		return nil, err
	}

	var configs []*Service

	for _, pair := range data {
		config := &Service{}
		err := json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}

		configs = append(configs, config)
	}

	return configs, nil
}
