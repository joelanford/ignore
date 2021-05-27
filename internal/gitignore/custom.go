package gitignore

import (
	"github.com/go-git/go-billy/v5"
	"os"
)

// ReadCustomPatterns reads patterns from each ignoreFile found recursively
// traversing through the directory structure. The result is in the ascending
// order of priority (last higher).
func ReadCustomPatterns(fs billy.Filesystem, ignoreFile string, path []string) (ps []Pattern, err error) {
	ps, _ = readIgnoreFile(fs, path, ignoreFile)

	var fis []os.FileInfo
	fis, err = fs.ReadDir(fs.Join(path...))
	if err != nil {
		return
	}

	for _, fi := range fis {
		if fi.IsDir() {
			var subps []Pattern
			subps, err = ReadCustomPatterns(fs, ignoreFile, append(path, fi.Name()))
			if err != nil {
				return
			}

			if len(subps) > 0 {
				ps = append(ps, subps...)
			}
		}
	}

	return
}
