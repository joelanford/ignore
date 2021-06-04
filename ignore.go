package ignore

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	)

func NewMatcher(root, ignoreFile string) (gitignore.Matcher, error) {
	patterns := []gitignore.Pattern{}
	if err := fs.WalkDir(os.DirFS(root), ".", func(path string, d fs.DirEntry, err error) error {
		if d.Name() != ignoreFile || d.IsDir() {
			return nil
		}
		ps, err := loadPatterns(root, path)
		if err != nil {
			return err
		}
		patterns = append(patterns, ps...)
		return nil
	}); err != nil {
		return nil, err
	}
	return gitignore.NewMatcher(patterns), nil
}

func loadPatterns(root, path string) ([]gitignore.Pattern, error) {
	file := filepath.Join(root, path)
	domain := strings.Split(filepath.Dir(path), string(filepath.Separator))
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	patterns := []gitignore.Pattern{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		patterns = append(patterns, gitignore.ParsePattern(line, domain))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return patterns, nil
}
