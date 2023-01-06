package photosorter

import (
	"os"
	"path"
	"testing"
)

var imgPath = path.Join("..", "..", "test", "img")

type SortTestCase struct {
	format string
	paths  []string
}

func (stc *SortTestCase) Run(t *testing.T) {
	src, dst := imgPath, imgPath+"/output"
	SortDir(src, dst, stc.format)

	for _, p := range stc.paths {
		_, err := os.Stat(dst + p)
		if os.IsNotExist(err) {
			t.Errorf("Missing image %s", p)
		}
	}

	os.RemoveAll(dst)
}

func TestSortDir(t *testing.T) {
	t.Run("Month directory structure", func(t *testing.T) {
		stc := &SortTestCase{
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
}
