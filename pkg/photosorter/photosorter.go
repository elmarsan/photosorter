package photosorter

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func SortDir(srcPath string, dstPath string) {
	_, err := os.Stat(dstPath)
	if os.IsNotExist(err) {
		err := os.Mkdir(dstPath, 0755)
		if err != nil {
			log.Panicf("While creating destination directory %s", dstPath)
			return
		}

		return
	}

	log.Printf("source directory: %s, destination directory: %s", srcPath, dstPath)

	files, err := ioutil.ReadDir(srcPath)
	if err != nil {
		log.Fatalf("Reading %s src directory", srcPath)
	}

	for _, file := range files {
		srcPath := strings.Join([]string{srcPath, file.Name()}, "/")
		dstPath := strings.Join([]string{dstPath, file.Name()}, "/")

		copyFile(srcPath, dstPath)
	}
}

func copyFile(srcPath string, dstPath string) {
	src, err := os.Open(srcPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Printf("File %s copied to: %s successfully", srcPath, dstPath)
}
