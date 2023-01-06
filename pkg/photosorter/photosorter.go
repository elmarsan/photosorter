package photosorter

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	"github.com/schollz/progressbar/v3"
)

// Image represents an image file and holds original creation time from metadata.
type Image struct {
	// File path of the image.
	src string
	// Raw image data.
	data []byte
	// original creation time of the image, extracted from the EXIF metadata.
	tm time.Time
}

// NewImage reads the content of an image file from the file system,
// extracts the EXIF metadata from the file,
// parses the EXIF metadata to extract the original creation time of the image,
// and then returns an *Image struct containing the image data, file path, and original creation time.
//
// Returns an error if any of the steps in this process fail,
// such as if the file cannot be read, if the EXIF metadata cannot be extracted or parsed,
// or if the original creation time cannot be extracted from the EXIF metadata.
func NewImage(src string) (*Image, error) {
	// Get src img content
	d, err := os.ReadFile(src)
	if err != nil {
		return nil, fmt.Errorf("Reading file %v", err)
	}

	// Extract exif
	rawExif, err := exif.SearchFileAndExtractExif(src)
	if err != nil && err.Error() != "no exif data" {
		return nil, fmt.Errorf("missing exif")
	}

	// Exif ifd mapping
	im, err := exifcommon.NewIfdMappingWithStandard()
	if err != nil {
		return nil, fmt.Errorf("exif format")
	}

	// Exif index for explore different tags
	ti := exif.NewTagIndex()
	_, index, err := exif.Collect(im, ti, rawExif)
	if err != nil {
		return nil, fmt.Errorf("exif format")
	}

	// Exif original creation time
	tagName := "DateTime"
	rootIfd := index.RootIfd
	results, err := rootIfd.FindTagWithName(tagName)
	if err != nil {
		return nil, fmt.Errorf("exif Tag: DateTime")
	}

	// Parse exif date to time.Time
	timeString, err := results[0].Format()
	if err != nil {
		return nil, fmt.Errorf("exif parsing Tag: DateTime")
	}

	tm, err := time.Parse("2006:01:02 15:04:05", timeString)
	if err != nil {
		return nil, fmt.Errorf("exif parsing Tag: DateTime")
	}

	return &Image{
		data: d,
		src:  src,
		tm:   tm,
	}, nil
}

// Save writes the image data to the specified destination directory.
//
// It uses the dst function to generate the file path for the image in the
// destination directory, and then writes the image data to that file. If the
// file cannot be written, it returns an error.
func (img *Image) Save(dir string, format string) error {
	dst := img.dst(dir, format)

	err := os.WriteFile(dst, img.data, 0755)
	if err != nil {
		return fmt.Errorf("Creating file %s. Cause: %v", dst, err)
	}

	return nil
}

// dst generates the destination file path for an image.
//
// It takes the destination directory and the format for organizing the
// images in the destination directory. If the format is "year", the image will
// be placed in a subdirectory named after the year of the image's original
// creation time. If the format is anything else, the image will be placed in a
// subdirectory named after the year, and a subdirectory named after the month
// of the image's original creation time.
//
// It creates the necessary subdirectories if they do not exist, and
// returns the file path for the image in the destination directory.
func (img *Image) dst(dstDir string, format string) string {
	year := strconv.Itoa(img.tm.Year())
	month := img.tm.Month().String()

	tree := []string{dstDir}

	if format == "year" {
		tree = append(tree, year)
	} else {
		tree = append(tree, year, month)
	}

	dir := strings.Join(tree, "/")
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Printf("Error creating dst directory %s \n", dir)
		}
	}

	fields := strings.Split(img.src, "/")
	fname := fields[len(fields)-1]
	dst := strings.Join([]string{dir, fname}, "/")

	return dst
}

func SortDir(src string, dst string, format string) {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		log.Fatalf("Reading %s src directory", src)
	}

	progress := int64(len(files))
	bar := progressbar.NewOptions(int(progress),
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[cyan]=[reset]",
			SaucerHead:    "[cyan]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	for i, file := range files {
		src := strings.Join([]string{src, file.Name()}, "/")
		description := fmt.Sprintf("[cyan][%d/%d][reset]: %s", i, len(files), src)
		bar.Describe(description)

		img, err := NewImage(src)
		if err != nil {
			fmt.Printf("Image %s could not be processed. Cause: %v\n", src, err)
		}

		err = img.Save(dst, format)
		if err != nil {
			fmt.Printf("Image %s could not be saved. Cause: %v\n", src, err)
		}

		bar.Add(1)
	}
}
