package main

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func get_letter(n int) string {
	result := []rune{}
	if n > len(letterBytes) {
		index := n / len(letterBytes)
		ex := n - len(letterBytes)*index
		for i := 0; i < index; i++ {
			result = append(result, []rune(letterBytes)[i])
		}

		if ex > 0 {
			result = append(result, []rune(letterBytes)[ex])
		}
		return string(result)
	}
	return string([]rune(letterBytes)[n])
}

func get_file_ext(file_name string) string {
	return filepath.Ext(file_name)
}

// Clean filename of all special characters and numbers
// from the file name begining until it find first letter.
// And return cleaner file_name
func name_cleaner(file fs.FileInfo) string {
	file_ext := get_file_ext(file.Name())
	file_name := strings.TrimSpace(strings.Replace(file.Name(), file_ext, "", 1))
	var index int

	for i, l := range file_name {
		index = i
		if unicode.IsLetter(l) {
			break
		}
	}

	return file_name[index:] + file_ext
}

func recreateFileName(path string, index int) string {
	// Get new name
	file_ext := get_file_ext(filepath.Base(path))
	file_name := strings.TrimSpace(strings.Replace(filepath.Base(path), file_ext, "", 1))

	if index > 0 {
		file_name = file_name[:len(file_name)-2]
	}
	file_name = file_name + " " + get_letter(index) + file_ext
	abs, _ := filepath.Abs(filepath.Join(strings.Replace(path, filepath.Base(path), "", 1), file_name))
	return abs
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
				i := 0
				newFilePath := name_cleaner(file)
				if newFilePath == file.Name() {
					continue
				}

				abs, err := filepath.Abs(filepath.Join(path, newFilePath))
				original_abs, _ := filepath.Abs(filepath.Join(path, file.Name()))
				for {
					if err != nil {
						log.Fatal(err)
					}
					_, err = os.Stat(abs)
					if os.IsNotExist(err) {
						break
					} else {
						abs = recreateFileName(abs, i)
					}
					i++
				}
				err = os.Rename(original_abs, abs)
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
