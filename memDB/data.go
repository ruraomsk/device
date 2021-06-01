package memDB

import (
	"fmt"
	"sync"
)

type mDB struct {
	mu   sync.RWMutex
	data map[string]interface{}
}
type Tx struct {
	name     string
	mdb      *mDB
	writable bool
	added    map[string]bool
	deleted  map[string]bool
	updated  map[string]bool
	ReadAll  func() map[string]interface{}
	WriteAll func() error
	AddFn    func(key string, value interface{}) string
	DeleteFn func(key string) string
	UpdateFn func(key string, value interface{}) string
}

func Create() *Tx {
	tx := Tx{writable: true, mdb: &mDB{data: make(map[string]interface{})}, updated: make(map[string]bool), added: make(map[string]bool), deleted: make(map[string]bool)}
	return &tx
}

func (tx *Tx) Set(key string, value interface{}) {
	//tx.Lock()
	//defer tx.Unlock()
	_, is := tx.mdb.data[key]
	tx.mdb.data[key] = value
	if !is {
		tx.added[key] = true
	} else {
		tx.updated[key] = true
	}
}
func (tx *Tx) GetAllKeys() []string {
	//tx.Lock()
	//defer tx.Unlock()
	result := make([]string, 0)
	for key := range tx.mdb.data {
		result = append(result, key)
	}
	return result
}

func (tx *Tx) Delete(key string) {
	//tx.Lock()
	//defer tx.Unlock()
	delete(tx.mdb.data, key)
	tx.deleted[key] = true
}

func (tx *Tx) Get(key string) (interface{}, error) {
	var err error
	_, is := tx.mdb.data[key]
	if !is {
		err = fmt.Errorf("нет такого ключа %s", key)
	}
	return tx.mdb.data[key], err
}

func (tx *Tx) Lock() {
	if tx.writable {
		tx.mdb.mu.Lock()
	} else {
		tx.mdb.mu.RLock()
	}
}

func (tx *Tx) Unlock() {
	if tx.writable {
		tx.mdb.mu.Unlock()
	} else {
		tx.mdb.mu.RUnlock()
	}
}
