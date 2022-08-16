package db

import (
	"crypto/sha1"
	"fmt"
	"sync"

	"github.com/ability-sh/abi-db/source"
	"github.com/ability-sh/abi-lib/dynamic"
	"github.com/ability-sh/abi-lib/json"
	"github.com/dop251/goja"
)

const wCount = 0x100

type exec func(s source.Source, vm *goja.Runtime)

type context struct {
	C chan error
}

func (c *context) Recycle() {
	close(c.C)
}

type collection struct {
	key string
	s   source.Source
	ss  chan exec
}

type collectionCursor struct {
	n int
	s source.Cursor
}

func (rs *collectionCursor) Next() (string, error) {
	r, e := rs.s.Next()
	if e != nil {
		return "", e
	}
	return r[rs.n:], nil
}

func (rs *collectionCursor) Close() {
	rs.s.Close()
}

type database struct {
	s  source.Source
	ss []chan exec
	w  sync.WaitGroup
}

func Open(s source.Source) Database {
	v := &database{s: s}
	v.ss = make([]chan exec, wCount)
	for i := 0; i < wCount; i++ {
		v.ss[i] = make(chan exec, 64)
		v.w.Add(1)

		go func(C chan exec) {

			defer v.w.Done()

			vm := goja.New()

			vm.Set("get", func(key string) string {
				b, err := s.Get(key)
				if err != nil {
					if err == source.ErrNoSuchKey {
						return ""
					}
					panic(vm.ToValue(err.Error()))
				}
				return string(b)
			})

			vm.Set("put", func(key string, data string) {
				err := s.Put(key, []byte(data))
				if err != nil {
					panic(vm.ToValue(err.Error()))
				}
			})

			vm.Set("del", func(key string) {
				err := s.Del(key)
				if err != nil {
					panic(vm.ToValue(err.Error()))
				}
			})

			for {

				cmd, ok := <-C

				if !ok {
					break
				}

				if cmd == nil {
					break
				}

				cmd(s, vm)
			}

		}(v.ss[i])
	}
	return v
}

func keyIndex(key string) int {
	m := sha1.New()
	m.Write([]byte(key))
	b := m.Sum(nil)
	return int(uint8(b[0]))
}

func (db *database) Collection(key string) Collection {
	return &collection{key: key, ss: db.ss[keyIndex(key)], s: db.s}
}

func (ct *collection) Key() string {
	return ct.key
}

func (ct *collection) Put(c Context, key string, data []byte) error {
	k := fmt.Sprintf("%s%s", ct.key, key)
	cc := c.(*context)
	ct.ss <- func(s source.Source, vm *goja.Runtime) {
		cc.C <- s.Put(k, data)
	}
	return <-cc.C
}

func (ct *collection) Del(c Context, key string) error {
	k := fmt.Sprintf("%s%s", ct.key, key)
	cc := c.(*context)
	ct.ss <- func(s source.Source, vm *goja.Runtime) {
		cc.C <- s.Del(k)
	}
	return <-cc.C
}

func (ct *collection) Get(c Context, key string) ([]byte, error) {
	k := fmt.Sprintf("%s%s", ct.key, key)
	return ct.s.Get(k)
}

func (ct *collection) GetObject(c Context, key string) (interface{}, error) {
	k := fmt.Sprintf("%s%s", ct.key, key)
	b, err := ct.s.Get(k)
	if err != nil {
		return nil, err
	}
	var r interface{} = nil
	err = json.Unmarshal(b, &r)
	return r, err
}

func (ct *collection) PutObject(c Context, key string, object interface{}) error {
	b, err := json.Marshal(object)
	if err != nil {
		return err
	}
	return ct.Put(c, key, b)
}

func (ct *collection) MergeObject(c Context, key string, object interface{}) error {
	k := fmt.Sprintf("%s%s", ct.key, key)
	cc := c.(*context)
	ct.ss <- func(s source.Source, vm *goja.Runtime) {
		var v interface{} = nil

		b, err := s.Get(k)
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
		cc.C <- s.Put(k, b)
	}
	return <-cc.C
}

func (ct *collection) Query(c Context, prefix string) (source.Cursor, error) {
	k := fmt.Sprintf("%s%s", ct.key, prefix)
	rs, err := ct.s.Query(k, "")
	if err != nil {
		return nil, err
	}
	return &collectionCursor{n: len(ct.key), s: rs}, nil
}

type execResult struct {
	data string
}

func (e *execResult) Error() string {
	return e.data
}

func (ct *collection) Exec(c Context, code string) (string, error) {
	k := ct.key
	cc := c.(*context)
	ct.ss <- func(s source.Source, vm *goja.Runtime) {
		vm.Set("collection", k)
		rs, err := vm.RunString(code)
		if err != nil {
			cc.C <- err
			return
		}
		cc.C <- &execResult{data: rs.String()}
	}
	err := <-cc.C
	r, ok := err.(*execResult)
	if ok {
		return r.data, nil
	}
	return "", err
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
