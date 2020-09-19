// +build !windows

package galadh

import (
	"os"
	"strings"
)

func isHidden(fi os.FileInfo) bool {
	return strings.HasPrefix(fi.Name(), ".")
}

func isExecutable(fi os.FileInfo) bool {
	mode := fi.Mode()
	return ((mode&0100 == 0100) || (mode&0010 == 0010) || (mode&0001 == 0001)) && !fi.IsDir()
}
