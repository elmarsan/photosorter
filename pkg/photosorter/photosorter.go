package photosorter

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
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
		return nil, fmt.Errorf("Reading file %v\n", err)
	}

	// Extract exif
	rawExif, err := exif.SearchFileAndExtractExif(src)
	if err != nil {
		return nil, fmt.Errorf("exif not found\n")
	}

	// Exif ifd mapping
	im, err := exifcommon.NewIfdMappingWithStandard()
	if err != nil {
		return nil, fmt.Errorf("exif format\n")
	}

	// Exif index for explore different tags
	ti := exif.NewTagIndex()
	_, index, err := exif.Collect(im, ti, rawExif)
	if err != nil {
		return nil, fmt.Errorf("exif format\n")
	}

	// Exif original creation time
	tagName := "DateTime"
	rootIfd := index.RootIfd
	results, err := rootIfd.FindTagWithName(tagName)
	if err != nil {
		return nil, fmt.Errorf("exif Tag: DateTime\n")
	}

	// Parse exif date to time.Time
	timeString, err := results[0].Format()
	if err != nil {
		return nil, fmt.Errorf("exif parsing Tag: DateTime\n")
	}

	tm, err := time.Parse("2006:01:02 15:04:05", timeString)
	if err != nil {
		return nil, fmt.Errorf("exif parsing Tag: DateTime\n")
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

	// fmt.Printf("save on %s\n", dst)

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
			fmt.Printf("Error creating dst directory %s \n", dir)
		}
	}

	fields := strings.Split(img.src, "/")
	fname := fields[len(fields)-1]
	dst := strings.Join([]string{dir, fname}, "/")

	return dst
}

// DirSortReport holds information about the sort directory proccess
type DirSortReport struct {
	// number of images processed.
	Imgs int
	// map of unprocessed files and the error message.
	Unprocessed map[string]string
}

// NewDirSortReport creates a new DirSortReport with default values.
func NewDirSortReport() *DirSortReport {
	return &DirSortReport{
		Imgs:        0,
		Unprocessed: make(map[string]string),
	}
}

// SortDir sorts the files in the specified source directory and saves them to the destination directory.
//
// Only files with the .jpg extension will be processed.
// The function returns a report with the elapsed time, number of images processed,
// and a map of unprocessed files and the error message.
func SortDir(src string, dst string, format string) (*DirSortReport, error) {
	report := NewDirSortReport()

	createCh := make(chan string)
	saveCh := make(chan (*Image))
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Goroutine for creating images
	// It receives by channel img src and tries to create it
	// If error occurs it adds the cause to report otherwise send img over saveCh
	go func() {
		wg.Add(1)
		defer wg.Done()

		for {
			path := <-createCh
			if path == "" {
				break
			}

			img, err := NewImage(path)
			if err != nil {
				mu.Lock()
				report.Unprocessed[path] = err.Error()
				mu.Unlock()
			} else {
				saveCh <- img
			}
		}

		close(createCh)
		saveCh <- nil
	}()

	// Goroutine for saving images
	// It receives by channel Image struct and tries to save it
	// If error occurs it adds the cause to report otherwise img is saved in dst directory
	go func() {
		wg.Add(1)
		defer wg.Done()

		for {
			img := <-saveCh
			if img == nil {
				break
			}

			err := img.Save(dst, format)
			mu.Lock()
			if err != nil {
				report.Unprocessed[img.src] = err.Error()
			} else {
				report.Imgs++
			}
			mu.Unlock()
		}

		close(saveCh)
	}()

	// Get file paths of the src directory
	paths, err := ScanDir(src)
	if err != nil {
		return nil, err
	}

	nfiles := len(*paths)

	// progress bar used for display sorting process information in real time
	progress := int64(nfiles)
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

	for i, p := range *paths {
		description := fmt.Sprintf("[cyan][%d/%d][reset]: %s", i, nfiles, p)
		bar.Describe(description)

		// Check if the file has the .jpg extension
		jpg, err := regexp.MatchString(".jpg", p)
		if err != nil {
			fmt.Printf("regex err: %v", err)
		}

		if jpg {
			createCh <- p
		} else {
			// File with no jpg extension are placed in unprocessed directly
			report.Unprocessed[p] = "Not .jpg extension"
		}

		bar.Add(1)
	}

	// Send empty as message for clossing channel
	createCh <- ""
	wg.Wait()

	return report, nil
}

// ScanDir scans the specified directory and returns a slice of strings
// containing the paths of all files in the directory and its subdirectories.
// It also checks if the directory exist.
func ScanDir(src string) (*[]string, error) {
	paths := []string{}

	// Check if the source directory exists
	_, err := ioutil.ReadDir(src)
	if os.IsNotExist(err) {
		return nil, err
	}

	// Walk through the files in the source directory
	filepath.Walk(src, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			paths = append(paths, path)
		}

		return nil
	})

	return &paths, nil
}
