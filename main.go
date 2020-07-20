package main

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type SearchValueForZip struct {
	findDir  string
	cacheDir string
	fileZip  string
	fileCSV  string
	value    string
}

func (s SearchValueForZip) String() string {
	if s.findDir == "" {
		panic("Не установлена директория поиска")
	}
	if s.cacheDir == "" {
		panic("Не установлен директория кеша")
	}
	if err := s.dellCache(); err != nil {
		panic(fmt.Sprintf("Какая-то ошибка удаления: %v\n", err))
	}
	if err := filepath.Walk(s.findDir, s.searchZip); err != nil {
		panic(fmt.Sprintf("Какая-то ошибка поиска: %v\n", err))
	}
	if err := s.unzip(); err != nil {
		panic(fmt.Sprintf("Какая-то ошибка разархивирования: %v\n", err))
	}
	if err := filepath.Walk(s.cacheDir, s.findCSV); err != nil {
		panic(fmt.Sprintf("Какая-то ошибка поиска: %v\n", err))
	}
	if err := s.dellCache(); err != nil {
		panic(fmt.Sprintf("Какая-то ошибка удаления: %v\n", err))
	}
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

func (s *SearchValueForZip) unzip() error {
	reader, err := zip.OpenReader(s.fileZip)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(s.cacheDir, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(s.cacheDir, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()
		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}
	return nil
}

func (s *SearchValueForZip) dellCache() error {
	if _, err := os.Stat(s.cacheDir); err == nil {
		return os.RemoveAll(s.cacheDir)
	}
	return nil

}

func (s *SearchValueForZip) findCSV(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err // Если по какой-то причине мы получили ошибку, проигнорируем эту итерацию
	}
	if info.IsDir() {
		return nil // Проигнорируем директории
	}
	if strings.Contains(info.Name(), ".txt") {
		return s.parseCSV(path)
	}
	return nil
}

func (s *SearchValueForZip) parseCSV(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	r := csv.NewReader(f)
	data, err := r.ReadAll()
	if err != nil {
		return err
	}
	i := 1
	for _, row := range data {
		if len(row) > 1 && i == 5 {
			s.value = row[2]
			return nil
		}
		i++
	}
	return nil
}

func main() {
	searchValueForZip := SearchValueForZip{findDir: ".", cacheDir: "./cache"}
	fmt.Print(searchValueForZip)

}
