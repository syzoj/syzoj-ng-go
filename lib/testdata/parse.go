package testdata

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

func ParseTestdata(path string) (*TestdataInfo, error) {
	info, err := os.Stat(filepath.Join(path, "data.yml"))
	if err == nil && !info.IsDir() {
		return ParseDataYml(path)
	}
	return ParseDefault(path)
}

func getFile(path string, name string) (*File, error) {
	f, err := os.Open(filepath.Join(path, name))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return nil, err
	}
	hash := hasher.Sum(nil)
	return &File{
		Name:      name,
		Sha256Sum: hex.EncodeToString(hash),
	}, nil
}

type fileSet map[string]*File

func getFileCached(path string, name string, fs fileSet) (*File, error) {
	if f, ok := fs[name]; ok {
		return f, nil
	}
	f, err := getFile(path, name)
	if err == nil {
		fs[name] = f
		return f, nil
	}
	return nil, err
}
