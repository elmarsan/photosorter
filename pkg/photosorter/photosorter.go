package photosorter

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

type Image struct {
	f      *os.File
	tm     time.Time
	dstDir string
}

func (img *Image) save() {
	dst := img.dst()

	f, err := os.Create(dst)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, img.f)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Image saved %s", dst)
}

func (img *Image) dst() string {
	year := strconv.Itoa(img.tm.Year())
	month := img.tm.Month().String()

	// Check if destinarion directory exists, otherwise it's created
	dir := strings.Join([]string{img.dstDir, year, month}, "/")
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Printf("Error creating dst directory %s \n", dir)
		}
	}

	name := strings.Split(img.f.Name(), "/")[1]
	dst := strings.Join([]string{img.dstDir, year, month, name}, "/")

	return dst
}

func SortDir(srcPath string, dstPath string) {
	log.Printf("source directory: %s, destination directory: %s", srcPath, dstPath)

	files, err := ioutil.ReadDir(srcPath)
	if err != nil {
		log.Fatalf("Reading %s src directory", srcPath)
	}

	for _, file := range files {
		srcPath := strings.Join([]string{srcPath, file.Name()}, "/")

		f, err := os.Open(srcPath)
		if err != nil {
			log.Panic(err)
		}

		defer f.Close()

		meta, err := exif.Decode(f)
		tm, err := meta.DateTime()
		if err != nil {
			log.Panic("While getting metadata time")
		}

		img := &Image{
			f:      f,
			tm:     tm,
			dstDir: dstPath,
		}

		img.save()
	}
}
