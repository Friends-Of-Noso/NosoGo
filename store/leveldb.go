package store

import (
	"fmt"
	"log"
	"reflect"

	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/protobuf/proto"
)

type ProtoStore[T proto.Message] struct {
	db     *leveldb.DB
	prefix string
}

func NewProtoStore[T proto.Message](db *leveldb.DB) *ProtoStore[T] {
	return &ProtoStore[T]{db: db}
}

// Optional: create a sub-store with a prefix (e.g. "blk:", "tx:", etc.)
func (ps *ProtoStore[T]) WithPrefix(prefix string) *ProtoStore[T] {
	return &ProtoStore[T]{db: ps.db, prefix: prefix}
}

func (ps *ProtoStore[T]) makeKey(key string) []byte {
	return []byte(ps.prefix + key)
}

// Put any proto.Message into LevelDB under the given key
func (ps *ProtoStore[T]) Put(key string, msg T) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	return ps.db.Put([]byte(ps.prefix+key), data, nil)
}

// Get retrieves a proto.Message by key and unmarshals it into T
func (ps *ProtoStore[T]) Get(key string) (T, error) {
	var zero T

	data, err := ps.db.Get(ps.makeKey(key), nil)
	if err != nil {
		return zero, err
	}

	// msg := any(new(T)).(T)
	msg := newMessage[T]()
	log.Printf("msg: %v", msg)

	if err := proto.Unmarshal(data, msg); err != nil {
		return zero, fmt.Errorf("unmarshal failed: %w", err)
	}

	return msg, nil
}

// Exists checks if a key exists
func (ps *ProtoStore[T]) Exists(key string) (bool, error) {
	return ps.db.Has(ps.makeKey(key), nil)
}

// Delete removes a key from the DB
func (ps *ProtoStore[T]) Delete(key string) error {
	return ps.db.Delete(ps.makeKey(key), nil)
}

// func newMessage[T proto.Message]() T {
// 	return reflect.New(reflect.TypeOf((*T)(nil)).Elem().Elem()).Interface().(T)
// }

func newMessage[T proto.Message]() T {
	// Assumes T is a pointer type (e.g., *pb.Block)
	typ := reflect.TypeOf((*T)(nil)).Elem() // get *pb.Block
	val := reflect.New(typ.Elem())          // allocate pb.Block → gives *pb.Block
	return val.Interface().(T)              // cast to T (which is *pb.Block)
}
