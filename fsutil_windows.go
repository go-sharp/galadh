package galadh

import (
	"os"
	"syscall"
)

func isHidden(fi os.FileInfo) bool {
	data, ok := fi.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		return false
	}

	return data.FileAttributes&syscall.FILE_ATTRIBUTE_HIDDEN == FILE_ATTRIBUTE_HIDDEN
}
