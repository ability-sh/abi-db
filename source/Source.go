package source

import "fmt"

type Cursor interface {
	Next() (string, error)
	Close()
}

type Source interface {
	Put(key string, data []byte) error
	Del(key string) error
	Get(key string) ([]byte, error)
	Query(prefix string) (Cursor, error)
}

var ErrNoSuchKey = fmt.Errorf("no such key")
