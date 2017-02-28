package main

// a single recursive copy with an intelligent decision tree, prepares
// dst parent folder if not exists, and recreates all symlinks using
// the first os.Readlink path retaining relative link structure

// loosely based on the two-function solution found here:
// @link: https://gist.github.com/m4ng0squ4sh/92462b38df26839a3ca324697c8cba04

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func copy(src, dst string) error {
	if src == "" {
		return nil
	}

	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	srcStat, err := os.Lstat(src)
	if err != nil {
		return err
	}

	if _, err := os.Lstat(dst); err != nil && !os.IsNotExist(err) {
		return err
	} else if err == nil {
		return fmt.Errorf("%s already exists...", dst)
	}

	if dstParent, err := os.Lstat(filepath.Dir(dst)); err == nil && !dstParent.IsDir() {
		return fmt.Errorf("cannot copy into a file: %s", filepath.Dir(dst))
	} else if err != nil && !os.IsNotExist(err) {
		return err
	} else if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	if srcStat.IsDir() {
		if err = os.Mkdir(dst, srcStat.Mode()); err != nil {
			return err
		}

		entries, err := ioutil.ReadDir(src)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if err = copy(filepath.Join(src, entry.Name()), filepath.Join(dst, entry.Name())); err != nil {
				return err
			}
		}
	} else {
		srcLstat, err := os.Lstat(src)
		if err != nil {
			return err
		} else if srcLstat.Mode()&os.ModeSymlink != 0 {
			if in, err := os.Readlink(src); err != nil {
				return err
			} else if err := os.Symlink(in, dst); err != nil {
				return err
			}
			return nil
		}

		in, err := os.Open(src)
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer func() {
			if e := out.Close(); err == nil && e != nil {
				err = e
			}
		}()

		if _, err = io.Copy(out, in); err != nil {
			return err
		} else if err = out.Sync(); err != nil {
			return err
		} else if err = out.Chmod(srcStat.Mode()); err != nil {
			return err
		}
	}

	return nil
}
