package static

import (
	"net/http"
	"path"
	"strings"

	"github.com/hverr/asset"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/httpfs"
)

type AssetFileSystem struct {
	http.FileSystem
	root    string
	indexes bool
	fs      vfs.FileSystem
}

func LoadAssetFileSystem(root string, indexes bool) (*AssetFileSystem, error) {
	fs := asset.Fs()
	if asset.Error() != nil {
		return nil, asset.Error()
	}

	if _, err := fs.Stat(root); err != nil {
		return nil, err
	}

	return &AssetFileSystem{
		FileSystem: httpfs.New(fs),
		root:       root,
		indexes:    indexes,
		fs:         fs,
	}, nil
}

func (fs *AssetFileSystem) Exists(prefix, filepath string) bool {
	var p string
	if p = strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		name := path.Join(fs.root, p)
		stats, err := fs.fs.Stat(name)
		if err != nil {
			return false
		}
		if !fs.indexes && stats.IsDir() {
			return false
		}
		return true
	}
	return false
}

func (fs *AssetFileSystem) Open(filename string) (http.File, error) {
	filename = path.Join(fs.root, filename)
	return fs.FileSystem.Open(filename)
}
