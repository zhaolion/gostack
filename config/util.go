package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/afero"
)

func absPath(inPath string) string {
	if strings.HasPrefix(inPath, "$") {
		end := strings.Index(inPath, string(os.PathSeparator))
		inPath = os.Getenv(inPath[1:end]) + inPath[end:]
	}

	if filepath.IsAbs(inPath) {
		return filepath.Clean(inPath)
	}

	p, err := filepath.Abs(inPath)
	if err == nil {
		return filepath.Clean(p)
	}

	return ""
}

// Check if File / Directory Exists
func exists(fs afero.Fs, path string) (bool, error) {
	_, err := fs.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func checkObject(obj interface{}) error {
	objVal := reflect.ValueOf(obj)
	if objVal.Kind() == reflect.Ptr && objVal.Elem().Kind() == reflect.Struct {
		return nil
	}

	return InvalidConfigTypeError("should be a pointer to a struct")
}

// copyObject return inner type copy
func copyObject(ptr interface{}) interface{} {
	return reflect.New(reflect.Indirect(reflect.ValueOf(ptr)).Type()).Interface()
}

func fileInfo(path string) (name, ext string) {
	if fileExt := filepath.Ext(path); len(fileExt) > 1 {
		ext = fileExt[1:]
	}

	name = strings.TrimSuffix(path, filepath.Ext(path))
	return
}
