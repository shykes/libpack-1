package libpack

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"gopkg.in/libgit2/git2go.v23"
)

func getAnnotation(db *DB, name string) (string, error) {
	return db.Get(mkAnnotation(name))
}

func setAnnotation(db *DB, name, value string) error {
	return db.Set(mkAnnotation(name), value)
}

func walkAnnotations(db *DB, h func(name, value string)) error {
	return db.Walk("/", func(k string, obj git.Object) error {
		blob, isBlob := obj.(*git.Blob)
		if !isBlob {
			return nil
		}
		targetPath, err := parseAnnotation(k)
		if err != nil {
			return err
		}
		h(targetPath, string(blob.Contents()))
		return nil
	})
}

func mkAnnotation(target string) string {
	target = treePath(target)
	if target == "/" {
		return "0"
	}
	return fmt.Sprintf("%d/%s", strings.Count(target, "/")+1, target)
}

func parseAnnotation(annot string) (target string, err error) {
	annot = treePath(annot)
	parts := strings.Split(annot, "/")
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid annotation path")
	}
	lvl, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return "", err
	}
	if len(parts)-1 != int(lvl) {
		return "", fmt.Errorf("invalid annotation path")
	}
	return path.Join(parts[1:]...), nil
}
