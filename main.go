package main

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SearchValueForZip struct {
	fileZip string
	value   string
	root    string
}

func (s SearchValueForZip) String() string {
	if s.root == "" {
		panic("Не установлен root")
	}
	if err := filepath.Walk(s.root, s.searchZip); err != nil {
		panic(fmt.Sprintf("Какая-то ошибка Хождения: %v\n", err))
	}
	fmt.Print("String " + s.fileZip)
	s.unzip()
	return s.value
}

func (s SearchValueForZip) searchZip(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err // Если по какой-то причине мы получили ошибку, проигнорируем эту итерацию
	}

	if info.IsDir() {
		return nil // Проигнорируем директории
	}
	if strings.Contains(info.Name(), ".zip") {
		s.fileZip = path
		fmt.Print("searchZip " + s.fileZip)
	}
	return nil
}

func (s SearchValueForZip) unzip() {
	fileZip, err := zip.OpenReader(s.fileZip)
	if err != nil {
		panic(fmt.Sprintf("Какая-то ошибка zip: %v\n", err))
	}
	defer fileZip.Close()
	for _, file := range fileZip.File {
		fmt.Printf("Contents of %s:\n", file.Name)
	}
}

func main() {
	searchValueForZip := SearchValueForZip{root: "."}
	fmt.Print(searchValueForZip)

}
