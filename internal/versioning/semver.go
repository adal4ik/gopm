package versioning

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
)

func extractVersion(fileName, packageName string) string {
	prefix := packageName + "-"
	if !strings.HasPrefix(fileName, prefix) {
		return ""
	}
	suffix := ".tar.gz"
	if !strings.HasSuffix(fileName, suffix) {
		return ""
	}

	return fileName[len(prefix) : len(fileName)-len(suffix)]
}

func FindBestMatch(fileNames []string, packageName, constraintStr string) (string, error) {
	if constraintStr == "" || constraintStr == "latest" {
		constraintStr = ">=0.0.0"
	}
	constraint, err := semver.NewConstraint(constraintStr)
	if err != nil {
		return "", fmt.Errorf("invalid version constraint '%s': %w", constraintStr, err)
	}

	var bestVersion *semver.Version
	var bestFileName string

	for _, fileName := range fileNames {
		versionStr := extractVersion(fileName, packageName)
		if versionStr == "" {
			continue
		}

		v, err := semver.NewVersion(versionStr)
		if err != nil {
			continue
		}

		if constraint.Check(v) {
			if bestVersion == nil || v.GreaterThan(bestVersion) {
				bestVersion = v
				bestFileName = fileName
			}
		}
	}

	if bestVersion == nil {
		return "", fmt.Errorf("no suitable version found for package '%s' with constraint '%s'", packageName, constraintStr)
	}

	return bestFileName, nil
}
