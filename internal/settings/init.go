package settings

import (
	"fmt"
	"log"
	"path/filepath"
)

func init() {
	err := DefinePaths()
	if err != nil {
		log.Fatal(err)
	}
}

func DefinePaths() (err error) {
	Path, err = filepath.Abs(Path)
	if err != nil {
		return fmt.Errorf("unable to get absolute path: %w", err)
	}

	return nil
}
