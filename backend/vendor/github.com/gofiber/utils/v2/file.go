package utils

import (
	"io"
	"net/http"
	"os"
	pathpkg "path"
	"path/filepath"
	"slices"
)

// Walk walks the filesystem rooted at root in lexical order, calling walkFn for each file or
// directory, including root. It follows filepath.Walk step for step, so filepath.SkipDir and
// filepath.SkipAll must be returned as they are, not wrapped.
//
// Without Lstat, Walk stats an entry by opening it and lists a directory with Readdir. An
// entry that cannot be opened for stat is reported with a nil FileInfo and the open error. A
// directory whose listing fails is reported with its FileInfo and the listing error.
//
// Unlike filepath.Walk, Walk follows symbolic links and, over http.Dir, never reports one: a link to
// an ancestor repeats the tree until ELOOP and one out of the tree is walked, so check paths against
// root before serving them. A named pipe can block the walk on open. fs and walkFn must not be nil,
// and a Stat returning neither a FileInfo nor an error panics.
func Walk(fs http.FileSystem, root string, walkFn filepath.WalkFunc) error {
	info, err := stat(fs, root)
	if err != nil {
		err = walkFn(root, nil, err)
	} else {
		err = walkInternal(fs, root, info, walkFn)
	}
	if isSentinel(err, filepath.SkipDir) || isSentinel(err, filepath.SkipAll) {
		return nil
	}
	return err
}

// isSentinel compares by identity, as filepath.Walk does; errors.Is would accept a wrapped sentinel.
func isSentinel(err, sentinel error) bool {
	return err == sentinel //nolint:errorlint,err113 // filepath.Walk compares the sentinels exactly
}

// ReadFile returns the raw content of a file
func ReadFile(path string, fs http.FileSystem) ([]byte, error) {
	if fs != nil {
		file, err := fs.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close() //nolint: errcheck // No need to check error
		return io.ReadAll(file)
	}
	return os.ReadFile(path) // #nosec G304
}

func readDirNames(fs http.FileSystem, dirname string) ([]string, error) {
	fis, err := readDir(fs, dirname)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(fis))
	for i := range fis {
		names[i] = fis[i].Name()
	}
	slices.Sort(names)
	return names, nil
}

// walkInternal is path/filepath's walk: it hands SkipDir and SkipAll back to the caller,
// except that a directory's own SkipDir ends there.
func walkInternal(fs http.FileSystem, path string, info os.FileInfo, walkFn filepath.WalkFunc) error {
	if !info.IsDir() {
		return walkFn(path, info, nil)
	}

	// Read before walkFn, so a read failure is reported by that one call.
	names, err := readDirNames(fs, path)
	err1 := walkFn(path, info, err)
	if err != nil || err1 != nil {
		if isSentinel(err1, filepath.SkipDir) {
			return nil
		}
		return err1
	}

	for _, name := range names {
		filename := pathpkg.Join(path, name)
		fileInfo, err := stat(fs, filename)
		if err != nil {
			if err := walkFn(filename, fileInfo, err); err != nil && !isSentinel(err, filepath.SkipDir) {
				return err
			}
		} else {
			err = walkInternal(fs, filename, fileInfo, walkFn)
			if err != nil {
				if !fileInfo.IsDir() || !isSentinel(err, filepath.SkipDir) {
					return err
				}
			}
		}
	}
	return nil
}

func readDir(fs http.FileSystem, name string) ([]os.FileInfo, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close() //nolint: errcheck // No need to check error
	return f.Readdir(0)
}

func stat(fs http.FileSystem, name string) (os.FileInfo, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close() //nolint: errcheck // No need to check error
	return f.Stat()
}
