package functions

import (
	"log"
	"path/filepath"
	"runtime"
)

func GetFile() string {
	_, file, _, _ := runtime.Caller(1)
	return filepath.Base(file)
}

func LogError(file string, err error) {
	log.Print("Error in ")
	log.Print(file)
	log.Println(":")
	log.Println(err)
}
