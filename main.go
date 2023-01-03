package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	src := flag.String("src", "", "source folder")
	dst := flag.String("dst", "", "destination folder")

	flag.Parse()

	if *src == "" {
		log.Fatal("Missing --src arg")
	}

	if *dst == "" {
		*dst = "./output"
	}

	_, err := os.Stat(*dst)
	if os.IsNotExist(err) {
		err := os.Mkdir(*dst, 0755)
		if err != nil {
			log.Panicf("While creating destination directory %s", *dst)
			return
		}

		return
	}

	log.Printf("source directory: %s, destination directory: %s", *src, *dst)

	files, err := ioutil.ReadDir(*src)
	if err != nil {
		log.Fatalf("Reading %s src directory", *src)
	}

	for _, file := range files {
		path := fmt.Sprintf("%s/%s", *src, file.Name())

		srcFile, err := os.Open(path)
		if err != nil {
			log.Panicf("Reading %s\n", path)
		}

		defer srcFile.Close()

		meta, err := exif.Decode(srcFile)

		tm, err := meta.DateTime()
		if err != nil {
			log.Panicf("Getting datetime of %s\n", path)
		}

		log.Printf("%s date: %v\n", file.Name(), tm)

		dstPath := fmt.Sprintf("%s/%s-%s", *dst, tm.Format("2006-01-02"), file.Name())

		dstFile, err := os.Create(dstPath)
		if err != nil {
			log.Panicf("Error creating file %s\n", dstPath)
		}

		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			log.Panicf("Error writting %s\n", dstPath)
		}
	}
}
