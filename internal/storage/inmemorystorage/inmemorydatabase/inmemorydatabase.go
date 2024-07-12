package inmemorydatabase

import (
	"errors"
	"sync"
)

var (
	// ErrNotFound is returned when the key is not found in the database.
	ErrNotFound = errors.New("key not found")

	// ErrDBClosed is returned when the database is closed.
	ErrDBClosed = errors.New("database closed")
)

// Simple in-memory key-value database.
// store value as byte slice for storing complex data after marshalling.
type InMemoryDatabase struct {
	db   map[string][]byte
	lock sync.RWMutex
}

func New() *InMemoryDatabase {
	return &InMemoryDatabase{
		db: make(map[string][]byte),
	}
}

func (d *InMemoryDatabase) Get(key string) ([]byte, error) {
	if d.db == nil {
		return nil, ErrDBClosed
	}

	d.lock.RLock()
	defer d.lock.RUnlock()

	if value, ok := d.db[key]; ok {
		return value, nil
	}

	return nil, ErrNotFound
}

func (d *InMemoryDatabase) Set(key string, value []byte) error {
	if d.db == nil {
		return ErrDBClosed
	}

	d.lock.Lock()
	defer d.lock.Unlock()

	d.db[key] = value

	return nil
}

func (d *InMemoryDatabase) Delete(key string) error {
	if d.db == nil {
		return ErrDBClosed
	}

	d.lock.Lock()
	defer d.lock.Unlock()

	delete(d.db, key)

	return nil
}

func (d *InMemoryDatabase) Keys() ([]string, error) {
	if d.db == nil {
		return nil, ErrDBClosed
	}

	d.lock.RLock()
	defer d.lock.RUnlock()
	keys := make([]string, 0, len(d.db))
	for key := range d.db {
		keys = append(keys, key)
	}
	return keys, nil
}

func (d *InMemoryDatabase) Close() {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.db = nil
}
