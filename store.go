package tidstore

import (
	"bufio"
	"encoding/gob"
	"log"
	"os"
)

// Store provides an interface for implementating new stores
type Store interface {
	// Get must create a Storage if it doesn't exists and return it.
	// Note that the Storage must never be nil, even if there is a error
	// creating or reading it. See the 'InStore' attribute in Storage.Options.
	Get(key string) (error, *Storage)

	// Save must save a Storage to the store. Must also set the
	// Storage.Options.IsNew boolean to false before saving it.
	Save(storage *Storage) error
}

const fsPermissions = 0600

// FilesystemStore provides an implementation of the Store interface that
// writes the storages to the filesystem, relative to the BasePath
type FilesystemStore struct {
	BasePath string
}

// NewFilesystemStore creates a new FilesystemStore with the OS' temporary
// directory as the BasePath
func NewFilesystemStore() *FilesystemStore {
	return &FilesystemStore{
		BasePath: os.TempDir(),
	}
}

// Get implements the Get method of the Store interface.
func (store *FilesystemStore) Get(key string) (*Storage, error) {
	path := store.BasePath + "/" + key
	file, err := os.OpenFile(path, os.O_RDONLY, fsPermissions)
	if err != nil {
		return NewStorage(key, DefaultMaxAge), err
	}
	in := bufio.NewReader(file)

	// Decode the storage from the file
	storage := new(Storage)
	dec := gob.NewDecoder(in)
	err = dec.Decode(storage)
	file.Close()
	log.Println(storage)
	return storage, err
}

// Save implements the Save method of the Store interface.
func (store *FilesystemStore) Save(storage *Storage) error {
	storage.Options.InStore = true
	storage.Options.IsNew = false

	path := store.BasePath + "/" + storage.Options.Key
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		fsPermissions)
	if err != nil {
		return err
	}
	out := bufio.NewWriter(file)

	// Encode the storage to GOB and write it to the file
	enc := gob.NewEncoder(out)
	err = enc.Encode(storage)
	out.Flush()
	file.Close()
	return err
}
