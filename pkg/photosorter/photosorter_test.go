package photosorter

import (
	"os"
	"path"
	"testing"
)

var imgPath = path.Join("..", "..", "test", "img")
var invalidImgPath = path.Join("..", "..", "test", "invalid_img")

type SortTestCase struct {
	src    string
	dst    string
	format string
	paths  []string
}

func (stc *SortTestCase) Run(t *testing.T) {
	SortDir(stc.src, stc.dst, stc.format)

	for _, p := range stc.paths {
		_, err := os.Stat(stc.dst + p)
		if os.IsNotExist(err) {
			t.Errorf("Missing image %s", p)
		}
	}

	os.RemoveAll(stc.dst)
}

func TestSortDir(t *testing.T) {
	t.Run("Month directory structure", func(t *testing.T) {
		stc := &SortTestCase{
			src:    imgPath,
			dst:    imgPath + "/output",
			format: "month",
			paths: []string{
				"/2007/September/mountain.jpg",
				"/2008/July/bug.jpg",
				"/2008/July/butterfly.jpg",
				"/2008/July/lizard.jpg",
				"/2008/November/tree.jpg",
				"/2014/September/leaf.jpg",
				"/2015/February/light.jpg",
			},
		}

		stc.Run(t)
	})

	t.Run("Year directory structure", func(t *testing.T) {
		stc := &SortTestCase{
			src:    imgPath,
			dst:    imgPath + "/output",
			format: "year",
			paths: []string{
				"/2007/mountain.jpg",
				"/2008/bug.jpg",
				"/2008/butterfly.jpg",
				"/2008/lizard.jpg",
				"/2008/tree.jpg",
				"/2014/leaf.jpg",
				"/2015/light.jpg",
			},
		}

		stc.Run(t)
	})

	t.Run("No exif", func(t *testing.T) {
		dst := invalidImgPath + "/output"

		stc := &SortTestCase{
			src:    invalidImgPath,
			dst:    dst,
			format: "year",
		}

		_, err := os.Stat(dst)
		if err == nil {
			t.Errorf("Output directory %s should not exist", dst)
		}

		stc.Run(t)
	})
}
