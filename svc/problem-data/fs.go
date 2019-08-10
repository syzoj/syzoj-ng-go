package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/syzoj/syzoj-ng-go/util"
)

func (app *App) ensurePath(name string) (string, error) {
	if len(name) != 32 {
		return "", fmt.Errorf("Invalid name")
	}
	path := filepath.Join(app.dataPath, name[0:2], name[2:4], name[4:6], name)
	err := os.MkdirAll(path, os.ModeDir|0755)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (app *App) makeTempFile() (*os.File, error) {
	tempPath := filepath.Join(app.dataPath, "temp")
	err := os.MkdirAll(tempPath, os.ModeDir|0755)
	if err != nil {
		return nil, err
	}
	hexName := util.RandomHex(16)
	return os.Create(filepath.Join(tempPath, hexName))
}

func (app *App) makeTempDir() (string, error) {
	hexName := util.RandomHex(16)
	path := filepath.Join(app.dataPath, "temp", hexName)
	err := os.MkdirAll(path, os.ModeDir|0755)
	if err != nil {
		return "", err
	}
	return path, nil
}
