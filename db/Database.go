package db

import "github.com/ability-sh/abi-db/source"

type Context interface {
	Recycle()
}

type Collection interface {
	Key() string
	Put(c Context, key string, data []byte) error
	Del(c Context, key string) error
	Get(c Context, key string) ([]byte, error)
	GetObject(c Context, key string) (interface{}, error)
	PutObject(c Context, key string, object interface{}) error
	MergeObject(c Context, key string, object interface{}) error
	Exec(c Context, code string) (string, error)
	Query(c Context, prefix string) (source.Cursor, error)
}

type Database interface {
	Collection(key string) Collection
	NewContext() Context
	Recycle()
}
