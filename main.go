package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
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

// Clean filename of all special characters and numbers
// from the file name begining until it find first letter.
// And return cleaner file_name
func name_cleaner(file fs.FileInfo) string {
	file_ext := get_file_ext(file.Name())
	file_name := sanitize_name(file.Name(), &file_ext)
	var index int

	for i, l := range file_name {
		index = i
		if unicode.IsLetter(l) {
			break
		}
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
				newFilePath := name_cleaner(file)
				if newFilePath == file.Name() {
					continue
				}

				abs, err := filepath.Abs(filepath.Join(path, newFilePath))
				if err != nil {
					log.Fatal(err)
				}

				_, err = os.Stat(abs)
				if os.IsNotExist(err) {
					fmt.Println("Process: ", abs)
				} else {
					fmt.Println("Rename: ", abs)
				}

				// err = os.Rename(abs, newFilePath)
				// if err != nil {
				// 	log.Fatal(err)
				// }
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
