package db

import "github.com/ability-sh/abi-db/source"

type Context interface {
	Recycle()
}

type Database interface {
	Put(c Context, key string, data []byte) error
	Del(c Context, key string) error
	Get(c Context, key string) ([]byte, error)
	GetObject(c Context, key string) (interface{}, error)
	PutObject(c Context, key string, object interface{}) error
	MergeObject(c Context, key string, object interface{}) error
	Query(c Context, prefix string) (source.Cursor, error)
	NewContext() Context
	Recycle()
}
