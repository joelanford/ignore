package ignore

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/joelanford/ignore/internal/gitignore"
)

func NewMatcher(root, ignoreFile string) (*Matcher, error) {
	indexfs := osfs.New(root)
	rootSegments := strings.Split(root, string(filepath.Separator))
	patterns, err := gitignore.ReadCustomPatterns(indexfs, ignoreFile, nil)
	if err != nil {
		return nil, fmt.Errorf("read patterns: %v", err)
	}
	base := gitignore.NewMatcher(patterns)
	return &Matcher{rootSegments, base}, nil
}

type Matcher struct {
	root    []string
	matcher gitignore.Matcher
}

func (m Matcher) Match(path string, isDir bool) bool {
	if !filepath.IsAbs(path) {
		rootPath := strings.Join(m.root, string(filepath.Separator))
		path = filepath.Clean(filepath.Join(rootPath, path))
	}
	pathSegments := strings.Split(path, string(filepath.Separator))
	return m.matcher.Match(pathSegments, isDir)
}
