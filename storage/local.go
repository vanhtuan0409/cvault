package storage

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/vanhtuan0409/cvault"
)

type localStorage struct {
	root string
}

func NewLocalStorage(storeUrl string) *localStorage {
	root := strings.TrimPrefix(storeUrl, "local://")
	return &localStorage{root: root}
}

func (s *localStorage) List(ctx context.Context) ([]*VaultItem, error) {
	ret := []*VaultItem{}
	entries, err := os.ReadDir(s.root)
	if err != nil {
		return ret, err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if name := e.Name(); cvault.IsEncryptedName(name) {
			fInfo, _ := e.Info()
			ret = append(ret, &VaultItem{
				Key:          name,
				LastModified: fInfo.ModTime(),
			})
		}
	}
	return ret, nil
}

func (s *localStorage) Get(ctx context.Context, key string) ([]byte, error) {
	p := filepath.Join(s.root, key)
	return ioutil.ReadFile(p)
}

func (s *localStorage) Put(ctx context.Context, key string, data []byte) error {
	if _, err := os.Stat(s.root); os.IsNotExist(err) {
		if err := os.MkdirAll(s.root, 0644); err != nil {
			return err
		}
	}

	p := filepath.Join(s.root, key)
	return ioutil.WriteFile(p, data, 0644)
}

func (s *localStorage) Remove(ctx context.Context, key string) error {
	p := filepath.Join(s.root, key)
	return os.Remove(p)
}
