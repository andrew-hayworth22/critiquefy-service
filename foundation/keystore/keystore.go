package keystore

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"path/filepath"
)

// Key represents an RSA key pair
type key struct {
	privatePEM string
	publicPEM  string
}

// KeyStore can load RSA key pairs from the filesystem
type KeyStore struct {
	store map[string]key
}

// New constructs a new key store
func New() *KeyStore {
	return &KeyStore{
		store: make(map[string]key),
	}
}

// LoadByFileSystem loads in key pairs by walking through a file system
func (ks *KeyStore) LoadByFileSystem(fsys fs.FS) (int, error) {
	fn := func(filePath string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walkdir failure: %w", err)
		}

		if dirEntry.IsDir() {
			ks.store[dirEntry.Name()] = key{}
			return nil
		}

		if path.Ext(filePath) != ".pem" {
			return nil
		}

		file, err := fsys.Open(filePath)
		if err != nil {
			return fmt.Errorf("opening key file: %w", err)
		}
		defer file.Close()

		pem, err := io.ReadAll(io.LimitReader(file, 1024*1024))
		if err != nil {
			return fmt.Errorf("reading auth key: %w", err)
		}

		parentDirName := filepath.Dir(filePath)
		fileName := filepath.Base(filePath)
		key := ks.store[parentDirName]
		if fileName == "private.pem" {
			key.privatePEM = string(pem)
			ks.store[parentDirName] = key
			return nil
		}

		key.publicPEM = string(pem)
		ks.store[parentDirName] = key
		return nil
	}

	if err := fs.WalkDir(fsys, ".", fn); err != nil {
		return 0, fmt.Errorf("walking directory: %w", err)
	}

	return len(ks.store), nil
}

// PublicKey fetches the public key given a key id
func (ks *KeyStore) PublicKey(kid string) (string, error) {
	key, ok := ks.store[kid]
	if !ok {
		return "", errors.New("kid lookup failed")
	}

	return key.publicPEM, nil
}

// PrivateKey fetches the private key given a key id
func (ks *KeyStore) PrivateKey(kid string) (string, error) {
	key, ok := ks.store[kid]
	if !ok {
		return "", errors.New("kid lookup failed")
	}

	return key.privatePEM, nil
}
