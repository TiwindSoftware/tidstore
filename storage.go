package tidstore

import (
	"time"
)

// DefaultMaxAge is a constant that can be used by stores for giving a
// storage a default max age. Is set to 86400 seconds (one day).
const DefaultMaxAge = 86400

// Options stores the options and the metadata of a specific Storage
type Options struct {
	Key      string
	MaxAge   int
	Creation time.Time
	InStore  bool
	IsNew    bool
}

// Storage stores the underlying in-memory key-value map, and its options.
type Storage struct {
	Values  map[interface{}]interface{}
	Options Options
}

// NewStorage creates a new storage, not in store, with the given
// key and max age.
func NewStorage(key string, maxAge int) *Storage {
	return &Storage{
		Values: make(map[interface{}]interface{}),
		Options: Options{
			Key:      key,
			MaxAge:   maxAge,
			Creation: time.Now(),
			InStore:  false,
			IsNew:    true,
		},
	}
}
