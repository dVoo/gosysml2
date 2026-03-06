package sysml

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func sysmlReleaseRoot(t testing.TB) string {
	t.Helper()

	candidates := make([]string, 0, 2)
	if env := os.Getenv("SYSML_V2_RELEASE_DIR"); env != "" {
		candidates = append(candidates, env)
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		candidates = append(candidates, filepath.Join(home, "projects", "SysML-v2-Release"))
	}

	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
	}

	t.Skip("SysML-v2-Release checkout not found; set SYSML_V2_RELEASE_DIR if needed")
	return ""
}

func validationTestdataDir(t *testing.T) string {
	t.Helper()
	dir := filepath.Join(sysmlReleaseRoot(t), "validationdata")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Skip("validationdata directory not found in SysML-v2-Release checkout")
	}
	return dir
}

func standardLibraryDir(t testing.TB) string {
	t.Helper()
	dir := filepath.Join(sysmlReleaseRoot(t), "libraries", "sysml.library")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Skip("standard library directory not found in SysML-v2-Release checkout")
	}
	return dir
}

func standardLibraryExists(t testing.TB) bool {
	t.Helper()
	dir := filepath.Join(sysmlReleaseRoot(t), "libraries", "sysml.library")
	_, err := os.Stat(dir)
	return err == nil
}

func externalTestPath(t testing.TB, legacy string) string {
	t.Helper()

	root := sysmlReleaseRoot(t)
	normalized := filepath.Clean(legacy)
	normalized = filepath.ToSlash(normalized)

	const validationPrefix = "../validationdata/"
	const libraryPrefix = "../libraries/sysml.library/"

	switch {
	case strings.HasPrefix(normalized, validationPrefix):
		rest := strings.TrimPrefix(normalized, validationPrefix)
		return filepath.Join(root, "validationdata", filepath.FromSlash(rest))
	case strings.HasPrefix(normalized, libraryPrefix):
		rest := strings.TrimPrefix(normalized, libraryPrefix)
		return filepath.Join(root, "libraries", "sysml.library", filepath.FromSlash(rest))
	default:
		return legacy
	}
}
