package main

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SearchValueForZip struct {
	findDir string
	fileZip string
	value   string
}

func (s SearchValueForZip) String() string {
	if s.findDir == "" {
		panic("Не установлена директория поиска")
	}
	if err := filepath.Walk(s.findDir, s.searchZip); err != nil {
		panic(fmt.Sprintf("Какая-то ошибка поиска: %v\n", err))
	}
	if err := s.unzip(); err != nil {
		panic(fmt.Sprintf("Какая-то ошибка разархивирования: %v\n", err))
	}
	return s.value
}

func (s *SearchValueForZip) searchZip(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}
	if strings.Contains(info.Name(), ".zip") {
		s.fileZip = path
	}
	return nil
}

func (s *SearchValueForZip) unzip() error {
	reader, err := zip.OpenReader(s.fileZip)
	if err != nil {
		return err
	}
	for _, file := range reader.File {
		fileOpen, _ := file.Open()
		if rows, _ := csv.NewReader(fileOpen).ReadAll(); len(rows) == 10 && len(rows[4]) == 10 {
			s.value = rows[4][2]
			return nil
		}
	}
	return nil
}

func main() {
	searchValueForZip := SearchValueForZip{findDir: "."}
	fmt.Print(searchValueForZip)

}
