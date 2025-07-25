package store

import (
	"fmt"
	"strings"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"google.golang.org/protobuf/proto"

	pb "github.com/Friends-Of-Noso/NosoGo/protobuf"
)

// Prefixes for different data types
const (
	StatusPrefix             = "status:"
	BlockPrefix              = "block:"
	TransactionPrefix        = "transaction:"
	PendingTransactionPrefix = "pending:"
	PeerInfoPrefix           = "peer:"
)

// ProtoMessage interface for protobuf messages
type ProtoMessage interface {
	proto.Message
}

// Storage represents a generic LevelDB storage
type Storage[T ProtoMessage] struct {
	mu     sync.RWMutex
	db     *leveldb.DB
	prefix string
}

// newStorage creates a new generic storage instance
func newStorage[T ProtoMessage](db *leveldb.DB, prefix string) *Storage[T] {
	return &Storage[T]{
		db:     db,
		prefix: prefix,
	}
}

// Put stores a protobuf message with the given key
func (s *Storage[T]) Put(key string, value T) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := proto.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal protobuf: %w", err)
	}

	fullKey := s.prefix + key
	return s.db.Put([]byte(fullKey), data, nil)
}

// Get retrieves a protobuf message by key
func (s *Storage[T]) Get(key string, value T) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fullKey := s.prefix + key
	data, err := s.db.Get([]byte(fullKey), nil)
	if err != nil {
		return fmt.Errorf("failed to get data: %w", err)
	}

	if err := proto.Unmarshal(data, value); err != nil {
		return fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return nil
}

// Delete removes a key from storage
func (s *Storage[T]) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	fullKey := s.prefix + key
	return s.db.Delete([]byte(fullKey), nil)
}

// Has checks if a key exists in storage
func (s *Storage[T]) Has(key string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fullKey := s.prefix + key
	return s.db.Has([]byte(fullKey), nil)
}

// ListKeys returns all keys with the storage prefix
func (s *Storage[T]) ListKeys() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var keys []string

	iter := s.db.NewIterator(util.BytesPrefix([]byte(s.prefix)), nil)
	defer iter.Release()

	for iter.Next() {
		key := string(iter.Key())
		// Remove prefix from key
		if strings.HasPrefix(key, s.prefix) {
			keys = append(keys, key[len(s.prefix):])
		}
	}

	return keys, iter.Error()
}

// ListValues returns all key-value pairs with the storage prefix
func (s *Storage[T]) ListValues(newInstance func() T) ([]T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []T

	iter := s.db.NewIterator(util.BytesPrefix([]byte(s.prefix)), nil)
	defer iter.Release()

	for iter.Next() {
		key := string(iter.Key())
		// Remove prefix from key
		if strings.HasPrefix(key, s.prefix) {
			cleanKey := key[len(s.prefix):]

			value := newInstance()
			if err := proto.Unmarshal(iter.Value(), value); err != nil {
				return nil, fmt.Errorf("failed to unmarshal value for key %s: %w", cleanKey, err)
			}

			results = append(results, value)
		}
	}

	return results, iter.Error()
}

// Batch operations
type Batch[T ProtoMessage] struct {
	mu     sync.Mutex
	batch  *leveldb.Batch
	prefix string
}

// NewBatch creates a new batch operation
func (s *Storage[T]) NewBatch() *Batch[T] {
	return &Batch[T]{
		batch:  new(leveldb.Batch),
		prefix: s.prefix,
	}
}

// Put adds a put operation to the batch
func (b *Batch[T]) Put(key string, value T) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	data, err := proto.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal protobuf: %w", err)
	}

	fullKey := b.prefix + key
	b.batch.Put([]byte(fullKey), data)
	return nil
}

// Delete adds a delete operation to the batch
func (b *Batch[T]) Delete(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	fullKey := b.prefix + key
	b.batch.Delete([]byte(fullKey))
}

// Write executes the batch
func (b *Batch[T]) Write(db *leveldb.DB) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	return db.Write(b.batch, nil)
}

// StorageManager manages multiple storage types
type StorageManager struct {
	db *leveldb.DB
}

// NewStorageManager creates a new storage manager
func NewStorageManager(dbPath string) (*StorageManager, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open leveldb: %w", err)
	}

	return &StorageManager{db: db}, nil
}

// Close closes the database connection
func (sm *StorageManager) Close() error {
	return sm.db.Close()
}

// GetDB returns the underlying LevelDB instance
func (sm *StorageManager) GetDB() *leveldb.DB {
	return sm.db
}

// Helper functions for creating specific storage instances
func (sm *StorageManager) StatusStorage() *Storage[*pb.Status] {
	return newStorage[*pb.Status](sm.db, StatusPrefix)
}

func (sm *StorageManager) BlockStorage() *Storage[*pb.Block] {
	return newStorage[*pb.Block](sm.db, BlockPrefix)
}

func (sm *StorageManager) TransactionStorage() *Storage[*pb.Transaction] {
	return newStorage[*pb.Transaction](sm.db, TransactionPrefix)
}

func (sm *StorageManager) PendingTransactionStorage() *Storage[*pb.Transaction] {
	return newStorage[*pb.Transaction](sm.db, PendingTransactionPrefix)
}

func (sm *StorageManager) PeerInfoStorage() *Storage[*pb.PeerInfo] {
	return newStorage[*pb.PeerInfo](sm.db, PeerInfoPrefix)
}

// Utility functions for key generation
func (sm *StorageManager) BlockKey(height uint64) string {
	return fmt.Sprintf("%016d", height) // Zero-padded for proper ordering
}

func (sm *StorageManager) TransactionKey(blockHeight uint64, txHash string) string {
	return fmt.Sprintf("%016d:%s", blockHeight, txHash)
}

func (sm *StorageManager) PeerKey(address string, id string) string {
	return fmt.Sprintf("%s:%s", address, id)
}

// Range query helpers
func (s *Storage[T]) GetRange(startKey, endKey string, newInstance func() T) (map[string]T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make(map[string]T)

	start := []byte(s.prefix + startKey)
	end := []byte(s.prefix + endKey)

	iter := s.db.NewIterator(&util.Range{Start: start, Limit: end}, nil)
	defer iter.Release()

	for iter.Next() {
		key := string(iter.Key())
		if strings.HasPrefix(key, s.prefix) {
			cleanKey := key[len(s.prefix):]

			value := newInstance()
			if err := proto.Unmarshal(iter.Value(), value); err != nil {
				return nil, fmt.Errorf("failed to unmarshal value for key %s: %w", cleanKey, err)
			}

			results[cleanKey] = value
		}
	}

	return results, iter.Error()
}

// Count returns the number of items with the storage prefix
func (s *Storage[T]) Count() (uint64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var count uint64 = 0

	iter := s.db.NewIterator(util.BytesPrefix([]byte(s.prefix)), nil)
	defer iter.Release()

	for iter.Next() {
		count++
	}

	return count, iter.Error()
}
