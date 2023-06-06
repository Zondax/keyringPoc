package main

import (
	"errors"
)

// fmt.Sprintf("%s.%s", r.Uid, "info")
type db map[string][]byte

func newDb() db {
	return make(db)
}

func (d db) save(id string, v []byte) {
	d[id] = v
}

func (d db) get(id string) ([]byte, error) {
	item, ok := d[id]
	if !ok {
		return nil, errors.New("key not found")
	}
	if len(item) == 0 {
		return nil, errors.New("error key")
	}
	return item, nil
}
