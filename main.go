package main

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
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
	s.unzip()
	return s.value
}

func (s *SearchValueForZip) searchZip(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err // Если по какой-то причине мы получили ошибку, проигнорируем эту итерацию
	}

	if info.IsDir() {
		return nil // Проигнорируем директории
	}
	if strings.Contains(info.Name(), ".zip") {
		s.fileZip = path
	}
	return nil
}

func (s *SearchValueForZip) unzip() {
	isFound := false
	fileZip, err := zip.OpenReader(s.fileZip)
	if err != nil {
		panic(fmt.Sprintf("Какая-то ошибка zip: %v\n", err))
	}
	defer fileZip.Close()
	for _, file := range fileZip.File {
		if strings.Contains(file.Name, ".txt") {
			if file.CompressedSize64 > 0 {
				isFound = s.parseFile(file)
			}
			if isFound == true {
				return
			}
		}
	}
}

func (s *SearchValueForZip) parseFile(file *zip.File) bool {
	rc, err := file.Open()
	if err != nil {
		panic(fmt.Sprintf("Какая-то ошибка Чтения Файла: %v\n", err))
	}
	_, err = io.CopyN(os.Stdout, rc, 68)
	if err != nil {
		panic(fmt.Sprintf("Какая-то ошибка Копирования Файла: %v\n", err))
	}
	rc.Close()
	content := fmt.Sprint()
	if content != "" {
		fmt.Println(file.Name)
		return s.parseCSV([]byte(content))
	}
	return false
}

func (s *SearchValueForZip) parseCSV(content []byte) bool {
	buf := bytes.NewBuffer(content)
	r := csv.NewReader(buf)
	for i := 1; i <= 2; i++ {
		// Читать данные мы тоже можем построчно, получая срез строк за каждую итерацию
		row, err := r.Read()
		if err != nil && err != io.EOF { // Здесь тоже нужно учитывать конец файла
			panic(fmt.Sprintf("Какая-то ошибка Разбора csv: %v\n", err))
		}
		fmt.Println(row)
	}
	return false
}

func main() {
	searchValueForZip := SearchValueForZip{root: "."}
	fmt.Print(searchValueForZip)

}
