package files

import (
	"fmt"
	"gopm/internal/config" // Мы зависим от типов, определенных в config
	"path/filepath"
)

// FindFilesByTarget находит все файлы, соответствующие одному таргету,
// применяя маску исключения.
func FindFilesByTarget(target config.Target) ([]string, error) {
	matches, err := filepath.Glob(target.Path)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern in path '%s': %w", target.Path, err)
	}

	if target.Exclude == "" {
		return matches, nil
	}

	var finalFiles []string
	for _, match := range matches {
		isExcluded, err := filepath.Match(target.Exclude, filepath.Base(match))
		if err != nil {
			return nil, fmt.Errorf("invalid pattern in exclude '%s': %w", target.Exclude, err)
		}

		if !isExcluded {
			finalFiles = append(finalFiles, match)
		}
	}
	return finalFiles, nil
}
