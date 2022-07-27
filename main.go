package main

import (
	"fmt"
	"os"
	"path/filepath"
	"qdb/cmd/qdb"
)

func touchConfigFile() error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	qdbdir := filepath.Join(homedir, ".qdbctl")
	if err := os.MkdirAll(qdbdir, os.ModePerm); err != nil {
		return err
	}

	file, err := os.OpenFile(filepath.Join(qdbdir, "config.yaml"), os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := touchConfigFile(); err != nil {
		fmt.Fprintf(os.Stderr, "Error while checking config file '%s'", err)
		os.Exit(1)
	}
	if err := qdb.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Whoops. There was an error while running qdb '%s'", err)
		os.Exit(1)
	}
}
