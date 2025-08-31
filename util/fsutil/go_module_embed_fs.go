package fsutil

import (
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

type goModuleEmbedFS struct {
	fs        fs.FS
	goModName string
}

func (g *goModuleEmbedFS) Open(name string) (fs.File, error) {
	if name == "go.mod" {
		return g.fs.Open(g.goModName)
	}
	return g.fs.Open(name)
}

func GoModuleEmbedFS(embeddedFS fs.FS, embeddedGoModName string) fs.FS {
	return &goModuleEmbedFS{
		fs:        embeddedFS,
		goModName: embeddedGoModName,
	}
}

func CopyFS(dst string, src fs.FS) error {
	embeddedFS, ok := src.(*goModuleEmbedFS)
	if !ok {
		return os.CopyFS(dst, src)
	}

	// Copied & modified from os.CopyFS
	return fs.WalkDir(src, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fpath, err := filepath.Localize(p)
		if err != nil {
			return err
		}

		if fpath == embeddedFS.goModName {
			fpath = "go.mod"
		}

		newPath := path.Join(dst, fpath)

		switch d.Type() {
		case os.ModeDir:
			return os.MkdirAll(newPath, 0777)
		case os.ModeSymlink:
			target, err := fs.ReadLink(src, p)
			if err != nil {
				return err
			}
			return os.Symlink(target, newPath)
		case 0:
			r, err := src.Open(p)
			if err != nil {
				return err
			}
			defer r.Close()
			info, err := r.Stat()
			if err != nil {
				return err
			}
			w, err := os.OpenFile(newPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666|info.Mode()&0777)
			if err != nil {
				return err
			}

			if _, err := io.Copy(w, r); err != nil {
				w.Close()
				return &os.PathError{Op: "Copy", Path: newPath, Err: err}
			}
			return w.Close()
		default:
			return &os.PathError{Op: "CopyFS", Path: p, Err: os.ErrInvalid}
		}
	})
}
