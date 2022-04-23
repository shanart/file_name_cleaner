package main

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func checkAlphaChar(charVariable rune) bool {
	if (charVariable >= 'a' && charVariable <= 'z') ||
		(charVariable >= 'A' && charVariable <= 'Z') ||
		(charVariable >= 'а' && charVariable <= 'я') ||
		(charVariable >= 'А' && charVariable <= 'Я') {
		return true
	} else {
		return false
	}
}

func is_path_exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func get_file_ext(file_name string) string {
	return filepath.Ext(file_name)
}

func sanitize_name(name string, ext *string) string {
	return strings.TrimSpace(strings.Replace(name, *ext, "", 1))
}

func name_cleaner(file fs.FileInfo) string {

	file_ext := get_file_ext(file.Name())
	file_name := sanitize_name(file.Name(), &file_ext)
	var index int

	for i, L := range file_name {
		if checkAlphaChar(rune(L)) {
			break
		}
		index = i
	}

	return file_name[index:] + file_ext
}

func main() {
	path := os.Args[1]
	pathInfo, err := os.Stat(path)
	if err != nil {
		log.Println("Path does not exits")
		return
	}

	if pathInfo.IsDir() {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			if !file.IsDir() {

				abs, err := filepath.Abs(filepath.Join(path, file.Name()))
				if err != nil {
					log.Fatal(err)
				}

				newFileName := name_cleaner(file)
				newFilePath, _ := filepath.Abs(filepath.Join(path, newFileName))

				// TODO: check if path exist
				// if true - alter filename and check egain.
				// while condition become false - it mean that path does not exists
				// and we can save the file with new path
				ex, err := is_path_exists(newFilePath)
				if ex != false {

				}

				err = os.Rename(abs, newFilePath)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	} else {
		file, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		name_cleaner(file)
	}
}
