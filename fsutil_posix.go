// +build !windows

package galadh

import (
	"os"
	"strings"
)

func isHidden(fi os.FileInfo) bool {
	return strings.HasPrefix(fi.Name(), ".")
}
