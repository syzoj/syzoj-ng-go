package testdata

import (
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
