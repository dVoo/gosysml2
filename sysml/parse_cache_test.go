package sysml

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func collectDirResults(t *testing.T, dir string, opts ...ParseOption) map[string]*ParseResult {
	t.Helper()

	results := make(map[string]*ParseResult)
	for result := range ParseDir(context.Background(), dir, DirOptions{
		Workers:      1,
		ParseOptions: opts,
	}) {
		if result == nil {
			continue
		}
		results[result.Source] = result
	}
	return results
}

func mustWriteModel(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestParseFileCacheReusesEntry(t *testing.T) {
	root := t.TempDir()
	cacheDir := filepath.Join(root, "cache")
	modelPath := filepath.Join(root, "model.sysml")
	mustWriteModel(t, modelPath, "package P { part def A; }")

	cache, err := NewParseCache(WithCacheDir(cacheDir), WithCachePersistence(true))
	if err != nil {
		t.Fatalf("NewParseCache: %v", err)
	}

	first := ParseFile(modelPath, WithParseCache(cache), WithDiscardTree())
	if err := first.Err(); err != nil {
		t.Fatalf("first parse: %v", err)
	}
	second := ParseFile(modelPath, WithParseCache(cache), WithDiscardTree())
	if err := second.Err(); err != nil {
		t.Fatalf("second parse: %v", err)
	}

	if first != second {
		t.Fatalf("expected cached ParseFile result to be reused")
	}
	if got := len(cache.entries); got != 1 {
		t.Fatalf("expected 1 cache entry, got %d", got)
	}
	if _, err := os.Stat(filepath.Join(cacheDir, cacheManifestName)); err != nil {
		t.Fatalf("expected manifest file: %v", err)
	}
}

func TestParseDirCacheReresolvesDependencyClosureOnChangedFile(t *testing.T) {
	root := t.TempDir()
	cache, err := NewParseCache(WithCacheDir(filepath.Join(root, "cache")), WithCachePersistence(true))
	if err != nil {
		t.Fatalf("NewParseCache: %v", err)
	}

	aPath := filepath.Join(root, "a.sysml")
	bPath := filepath.Join(root, "b.sysml")
	mustWriteModel(t, aPath, "package P { part def A; }")
	mustWriteModel(t, bPath, "package P { part b : A; }")

	first := collectDirResults(t, root, WithParseCache(cache), WithDiscardTree())
	b1 := first[bPath]
	if b1 == nil || b1.Model == nil {
		t.Fatalf("missing cached result for %s", bPath)
	}
	pkg1 := b1.Model.Packages()[0]
	usage1 := pkg1.Parts()[0]
	firstResolved := usage1.TypeRef.Resolved()
	if firstResolved == nil {
		t.Fatalf("expected part usage to resolve on first ParseDir")
	}

	mustWriteModel(t, aPath, "package P {   part def A; }\n")

	second := collectDirResults(t, root, WithParseCache(cache), WithDiscardTree())
	b2 := second[bPath]
	if b2 == nil || b2.Model == nil {
		t.Fatalf("missing second cached result for %s", bPath)
	}
	usage2 := b2.Model.Packages()[0].Parts()[0]
	secondResolved := usage2.TypeRef.Resolved()
	if secondResolved == nil {
		t.Fatalf("expected part usage to resolve on second ParseDir")
	}
	if firstResolved == secondResolved {
		t.Fatalf("expected dependent file to be re-resolved after changed dependency")
	}
}

func TestParseCacheReloadsManifestFromDirectory(t *testing.T) {
	root := t.TempDir()
	cacheDir := filepath.Join(root, "cache")
	modelPath := filepath.Join(root, "model.sysml")
	mustWriteModel(t, modelPath, "package P { part def A; }")

	cache1, err := NewParseCache(WithCacheDir(cacheDir), WithCachePersistence(true))
	if err != nil {
		t.Fatalf("NewParseCache cache1: %v", err)
	}
	result := ParseFile(modelPath, WithParseCache(cache1), WithDiscardTree())
	if err := result.Err(); err != nil {
		t.Fatalf("ParseFile with cache1: %v", err)
	}

	cache2, err := NewParseCache(WithCacheDir(cacheDir), WithCachePersistence(true))
	if err != nil {
		t.Fatalf("NewParseCache cache2: %v", err)
	}
	if got := len(cache2.manifest.Files); got == 0 {
		t.Fatalf("expected manifest entries to be reloaded")
	}

	reparsed := ParseFile(modelPath, WithParseCache(cache2), WithDiscardTree())
	if err := reparsed.Err(); err != nil {
		t.Fatalf("ParseFile with cache2: %v", err)
	}
	if got := len(cache2.entries); got != 1 {
		t.Fatalf("expected one in-memory entry after reparsing, got %d", got)
	}
}
