package fs

import (
	"path/filepath"
	"strings"
)

const Separator = string(filepath.Separator)

func DirHierarchy(dir string) []string {
	dir, _ = filepath.Abs(dir)

	if dir == "" || dir == Separator {
		return []string{Separator}
	}

	dirs := strings.Split(dir, Separator)
	dirs[0] = Separator

	result := make([]string, 0, len(dirs))
	for i := len(dirs) - 1; i >= 0; i-- {
		result = append(result, filepath.Join(dirs[:i+1]...))
	}

	return result
}
