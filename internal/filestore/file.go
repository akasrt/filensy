package filestore

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/rs/xid"
)

const dirName = "vault"

var store fileStore

type FileStore interface {
	Save(src io.Reader) (string, uint64, error)
	Get(name string) (*os.File, error)
	Delete(name string) error
}

func InitFileStore(root string) {
	dir := filepath.Join(root, dirName)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Panic("failed to init file directory: ", err)
	}

	store = fileStore{
		dir: dir,
	}
}

func Get() FileStore {
	return &store
}

type fileStore struct {
	dir string
}

func (f *fileStore) Save(src io.Reader) (string, uint64, error) {
	storageKey := generateStorageKey()
	path := filepath.Join(f.dir, storageKey)

	dst, err := os.Create(path)
	if err != nil {
		return "", 0, err
	}
	defer dst.Close()

	n, err := io.Copy(dst, src)
	if err != nil {
		_ = os.Remove(path)
		return "", 0, err
	}

	return storageKey, uint64(n), nil
}

func (f *fileStore) Get(name string) (*os.File, error) {
	path := filepath.Join(f.dir, filepath.Clean(name))
	return os.Open(path)
}

func (f *fileStore) Delete(name string) error {
	path := filepath.Join(f.dir, filepath.Clean(name))
	return os.Remove(path)
}

func generateStorageKey() string {
	return xid.New().String()
}
