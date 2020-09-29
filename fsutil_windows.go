package galadh

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func isHidden(fi os.FileInfo) bool {
	data, ok := fi.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		return false
	}

	return data.FileAttributes&syscall.FILE_ATTRIBUTE_HIDDEN == syscall.FILE_ATTRIBUTE_HIDDEN
}

func isExecutable(fi os.FileInfo) bool {
	return strings.EqualFold(filepath.Ext(fi.Name()), ".exe")
}
