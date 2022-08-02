//go:build mage

//go:mock mage

package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	app             = "tec-doc"
	mockDestination = "./internal/tec-doc/mocks"
	sourceFiles     = []string{
		"internal/tec-doc/web/externalserver/server.go",
		"internal/tec-doc/web/internalserver/server.go",
		"internal/tec-doc/store/postgres/store.go",
	}
)

//goland:noinspection GoBoolExpressions
func init() {
	//validation source file
	for _, sc := range sourceFiles {
		if !strings.HasSuffix(sc, ".go") {
			log.Fatalln("error:", sc, "file isn't go file")
		}
	}
	if runtime.GOOS == "windows" {
		app += ".exe"
	}
}

// Runs go mod download and then installs the binary.
func Build() error {
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}
	return sh.Run("go", "build", "-o", app, "./cmd/tec-doc/main.go")
}

//Generate mock
func Mock() (err error) {
	for _, sc := range sourceFiles {
		dirPath, fileName := filepath.Split(sc)
		dirPath = "mock_" + filepath.Base(dirPath)
		fileName = "mock_" + fileName
		destination := filepath.Clean(strings.Join([]string{mockDestination, dirPath, fileName}, string(filepath.Separator)))
		if err = sh.Run("mockgen", "-source", sc, "-destination", destination); err != nil {
			return err
		}
		fmt.Println("destination:", destination)
	}
	return nil
}