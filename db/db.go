package db

import (
	"crypto/sha1"
	"log"
	"sync"

	"github.com/ability-sh/abi-db/source"
	"github.com/ability-sh/abi-lib/dynamic"
	"github.com/ability-sh/abi-lib/json"
)

const wCount = 0x100

type context struct {
	C chan error
}

func (c *context) Recycle() {
	close(c.C)
}

type database struct {
	s  source.Source
	ss []chan func(s source.Source)
	w  sync.WaitGroup
}

func Open(s source.Source) Database {
	v := &database{s: s}
	v.ss = make([]chan func(s source.Source), wCount)
	for i := 0; i < wCount; i++ {
		v.ss[i] = make(chan func(s source.Source), 64)
		v.w.Add(1)

		go func(C chan func(s source.Source)) {

			defer v.w.Done()

			for {

				cmd, ok := <-C

				if !ok {
					break
				}

				if cmd == nil {
					break
				}

				cmd(s)
			}

		}(v.ss[i])
	}
	return v
}

func keyIndex(key string) int {
	m := sha1.New()
	m.Write([]byte(key))
	b := m.Sum(nil)
	log.Println("keyIndex", b)
	return int(uint8(b[0]))
}

func (db *database) Put(c Context, key string, data []byte) error {
	cc := c.(*context)
	db.ss[keyIndex(key)] <- func(s source.Source) {
		cc.C <- s.Put(key, data)
	}
	return <-cc.C
}

func (db *database) Del(c Context, key string) error {
	cc := c.(*context)
	db.ss[keyIndex(key)] <- func(s source.Source) {
		cc.C <- s.Del(key)
	}
	return <-cc.C
}

func (db *database) Get(c Context, key string) ([]byte, error) {
	return db.s.Get(key)
}

func (db *database) GetObject(c Context, key string) (interface{}, error) {
	b, err := db.s.Get(key)
	if err != nil {
		return nil, err
	}
	var r interface{} = nil
	err = json.Unmarshal(b, &r)
	return r, err
}

func (db *database) PutObject(c Context, key string, object interface{}) error {
	b, err := json.Marshal(object)
	if err != nil {
		return err
	}
	return db.Put(c, key, b)
}

func (db *database) MergeObject(c Context, key string, object interface{}) error {
	cc := c.(*context)
	db.ss[keyIndex(key)] <- func(s source.Source) {
		var v interface{} = nil

		b, err := s.Get(key)
		if err != nil {
			if err != source.ErrNoSuchKey {
				cc.C <- err
				return
			}
		} else {
			json.Unmarshal(b, &v)
		}

		if v == nil {
			v = object
		} else {
			dynamic.Each(object, func(key interface{}, value interface{}) bool {
				dynamic.Set(v, dynamic.StringValue(key, ""), value)
				return true
			})
		}
		b, err = json.Marshal(v)
		if err != nil {
			cc.C <- err
			return
		}
		cc.C <- s.Put(key, b)
	}
	return <-cc.C
}

func (db *database) NewContext() Context {
	return &context{C: make(chan error)}
}

func (db *database) Recycle() {
	for i := 0; i < wCount; i++ {
		db.ss[i] <- nil
	}
	db.w.Wait()
	for i := 0; i < wCount; i++ {
		close(db.ss[i])
	}
}
