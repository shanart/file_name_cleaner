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

func checkAlphaChar(charVariable rune) bool {

	if (unicode.Is(unicode.Cyrillic, charVariable)) ||
		(unicode.IsLetter(charVariable)) ||
		(strings.ContainsRune("!@#$%^&*()_+~", charVariable)) {
		return true
	}

	return false
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

func name_cleaner(file fs.FileInfo, alter string) string {

	file_ext := get_file_ext(file.Name())
	file_name := sanitize_name(file.Name(), &file_ext)
	var index int

	for i, L := range file_name {
		if checkAlphaChar(rune(L)) {
			break
		}
		index = i
	}

	a_file_name := file_name[index:]
	if len(alter) > 0 {
		a_file_name = a_file_name[index:] + alter
	}

	ex, _ := is_path_exists(a_file_name + file_ext)
	if ex {
		return name_cleaner(file, "_"+randomString(3))
	} else {
		return file_name[index:] + file_ext
	}
}

func main() {

	fmt.Println()

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

				// abs, err := filepath.Abs(filepath.Join(path, file.Name()))
				// if err != nil {
				// 	log.Fatal(err)
				// }

				newFileName := name_cleaner(file, "")
				newFilePath, _ := filepath.Abs(filepath.Join(path, newFileName))

				fmt.Println(newFilePath)
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
		name_cleaner(file, "")
	}
}
