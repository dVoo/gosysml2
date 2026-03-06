package sysml

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

const (
	parseCacheVersion   = "v1"
	grammarVersion      = "sysml-parser-2026-03-08"
	cacheManifestName   = "manifest.json"
	cacheLockName       = "lock"
	cacheFilesDirectory = "files"
)

type ParseCacheOption func(*parseCacheConfig)

type parseCacheConfig struct {
	dir        string
	temporary  bool
	persistent bool
}

type ParseCache struct {
	dir        string
	temporary  bool
	persistent bool

	mu       sync.RWMutex
	manifest *cacheManifest
	entries  map[string]*cachedFileUnit
}

type parseModeFingerprint struct {
	BuildModel           bool `json:"build_model"`
	ResolveModel         bool `json:"resolve_model"`
	DiscardTree          bool `json:"discard_tree"`
	DisableRewrites      bool `json:"disable_rewrites"`
	HasLibraryResolution bool `json:"has_library_resolution"`
}

type cachedImport struct {
	Namespace string   `json:"namespace"`
	IsAll     bool     `json:"is_all"`
	IsRec     bool     `json:"is_recursive"`
	Filters   []string `json:"filters,omitempty"`
}

type cachedFileUnit struct {
	SourcePath         string               `json:"source_path"`
	ContentHash        string               `json:"content_hash"`
	ParseMode          parseModeFingerprint `json:"parse_mode"`
	GrammarVersion     string               `json:"grammar_version"`
	LibraryFingerprint string               `json:"library_fingerprint"`
	ExportSignature    string               `json:"export_signature"`
	TopLevelPackages   []string             `json:"top_level_packages,omitempty"`
	Imports            []cachedImport       `json:"imports,omitempty"`
	Diagnostics        []*Error             `json:"diagnostics,omitempty"`
	Rewrites           []string             `json:"rewrites,omitempty"`
	Result             *ParseResult         `json:"-"`
	ResolvedForRepo    bool                 `json:"-"`
}

type cacheManifest struct {
	Version        string                    `json:"version"`
	GrammarVersion string                    `json:"grammar_version"`
	Files          map[string]cachedFileMeta `json:"files"`
	ReverseDeps    map[string][]string       `json:"reverse_deps"`
}

type cachedFileMeta struct {
	SourcePath         string               `json:"source_path"`
	ContentHash        string               `json:"content_hash"`
	ExportSignature    string               `json:"export_signature"`
	EntryFile          string               `json:"entry_file"`
	ParseMode          parseModeFingerprint `json:"parse_mode"`
	LibraryFingerprint string               `json:"library_fingerprint"`
	TopLevelPackages   []string             `json:"top_level_packages,omitempty"`
	Imports            []cachedImport       `json:"imports,omitempty"`
}

type RepoIndex struct {
	elementIndex     map[string]Element
	shortNameIndex   map[string][]Element
	libraryRegistry  *LibraryRegistry
	elementOwnerPath map[string]string
}

type fileUnit struct {
	Source          string
	Result          *ParseResult
	ExportSignature string
	TopLevelPkgs    []string
	Imports         []cachedImport
}

func WithCacheDir(dir string) ParseCacheOption {
	return func(c *parseCacheConfig) {
		c.dir = dir
	}
}

func WithTemporaryCacheDir() ParseCacheOption {
	return func(c *parseCacheConfig) {
		c.temporary = true
	}
}

func WithCachePersistence(enabled bool) ParseCacheOption {
	return func(c *parseCacheConfig) {
		c.persistent = enabled
	}
}

func NewParseCache(opts ...ParseCacheOption) (*ParseCache, error) {
	cfg := parseCacheConfig{persistent: true}
	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.dir == "" {
		cfg.temporary = true
	}

	var err error
	if cfg.dir == "" {
		cfg.dir, err = os.MkdirTemp("", "gosysml2-cache-")
		if err != nil {
			return nil, fmt.Errorf("create parse cache dir: %w", err)
		}
	}
	if err := os.MkdirAll(filepath.Join(cfg.dir, cacheFilesDirectory), 0o755); err != nil {
		return nil, fmt.Errorf("create parse cache layout: %w", err)
	}

	cache := &ParseCache{
		dir:        cfg.dir,
		temporary:  cfg.temporary,
		persistent: cfg.persistent,
		manifest: &cacheManifest{
			Version:        parseCacheVersion,
			GrammarVersion: grammarVersion,
			Files:          make(map[string]cachedFileMeta),
			ReverseDeps:    make(map[string][]string),
		},
		entries: make(map[string]*cachedFileUnit),
	}
	cache.loadManifest()
	cache.ensureLockFile()
	return cache, nil
}

func (c *ParseCache) Dir() string {
	if c == nil {
		return ""
	}
	return c.dir
}

func (c *ParseCache) Close() error {
	if c == nil {
		return nil
	}
	if c.temporary && !c.persistent {
		return os.RemoveAll(c.dir)
	}
	return nil
}

func (c *ParseCache) Clear() error {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]*cachedFileUnit)
	c.manifest = &cacheManifest{
		Version:        parseCacheVersion,
		GrammarVersion: grammarVersion,
		Files:          make(map[string]cachedFileMeta),
		ReverseDeps:    make(map[string][]string),
	}

	if err := os.RemoveAll(filepath.Join(c.dir, cacheFilesDirectory)); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(c.dir, cacheFilesDirectory), 0o755); err != nil {
		return err
	}
	return c.persistManifestLocked()
}

func (c *ParseCache) loadManifest() {
	path := filepath.Join(c.dir, cacheManifestName)
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var manifest cacheManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return
	}
	if manifest.Version != parseCacheVersion || manifest.GrammarVersion != grammarVersion {
		return
	}
	if manifest.Files == nil {
		manifest.Files = make(map[string]cachedFileMeta)
	}
	if manifest.ReverseDeps == nil {
		manifest.ReverseDeps = make(map[string][]string)
	}
	c.manifest = &manifest
}

func (c *ParseCache) ensureLockFile() {
	lockPath := filepath.Join(c.dir, cacheLockName)
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o644)
	if err == nil {
		_ = f.Close()
	}
}

func (c *ParseCache) persistManifestLocked() error {
	if c == nil || !c.persistent {
		return nil
	}
	data, err := json.MarshalIndent(c.manifest, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(c.dir, cacheManifestName), data, 0o644)
}

func WithParseCache(cache *ParseCache) ParseOption {
	return func(c *parseConfig) {
		c.cache = cache
	}
}

func WithDefaultParseCache() ParseOption {
	var (
		once  sync.Once
		cache *ParseCache
		err   error
	)
	return func(c *parseConfig) {
		once.Do(func() {
			cache, err = NewParseCache(WithTemporaryCacheDir(), WithCachePersistence(false))
		})
		if err != nil {
			c.cacheErr = err
			return
		}
		c.cache = cache
	}
}

func WithParseCacheDir(dir string) ParseOption {
	var (
		once  sync.Once
		cache *ParseCache
		err   error
	)
	return func(c *parseConfig) {
		once.Do(func() {
			cache, err = NewParseCache(WithCacheDir(dir), WithCachePersistence(true))
		})
		if err != nil {
			c.cacheErr = err
			return
		}
		c.cache = cache
	}
}

func parseModeForConfig(cfg *parseConfig) parseModeFingerprint {
	return parseModeFingerprint{
		BuildModel:           cfg.buildModel,
		ResolveModel:         cfg.resolveModel,
		DiscardTree:          cfg.discardTree,
		DisableRewrites:      cfg.disableRewrites,
		HasLibraryResolution: cfg.libraryRegistry != nil || cfg.autoLoadLibraries,
	}
}

func libraryFingerprint(reg *LibraryRegistry, cfg *parseConfig) string {
	if reg == nil {
		return "none"
	}
	reg.mu.RLock()
	defer reg.mu.RUnlock()

	paths := append([]string(nil), reg.libraryPaths...)
	sort.Strings(paths)
	parts := []string{
		fmt.Sprintf("loaded=%t", reg.loaded),
		fmt.Sprintf("elements=%d", len(reg.elementIndex)),
	}
	parts = append(parts, paths...)
	if cfg != nil && cfg.libraryPath != "" {
		parts = append(parts, "cfg="+cfg.libraryPath)
	}
	sum := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return fmt.Sprintf("%x", sum[:])
}

func parseCacheEligibleSource(source string) (string, bool) {
	if source == "" || strings.Contains(source, "://") {
		return "", false
	}
	ext := strings.ToLower(filepath.Ext(source))
	if ext != ".sysml" && ext != ".kerml" {
		return "", false
	}
	abs, err := filepath.Abs(source)
	if err != nil {
		return "", false
	}
	return abs, true
}

func cacheEntryKey(sourcePath string, mode parseModeFingerprint, libraryFP string) string {
	payload := fmt.Sprintf("%s|%t|%t|%t|%t|%t|%s",
		sourcePath,
		mode.BuildModel,
		mode.ResolveModel,
		mode.DiscardTree,
		mode.DisableRewrites,
		mode.HasLibraryResolution,
		libraryFP,
	)
	sum := sha256.Sum256([]byte(payload))
	return fmt.Sprintf("%x", sum[:])
}

func cacheEntryFile(entryKey string) string {
	return entryKey + ".json"
}

func (c *ParseCache) get(sourcePath string, mode parseModeFingerprint, libraryFP, contentHash string) (*cachedFileUnit, bool) {
	if c == nil {
		return nil, false
	}
	key := cacheEntryKey(sourcePath, mode, libraryFP)

	c.mu.RLock()
	defer c.mu.RUnlock()

	meta, ok := c.manifest.Files[key]
	if !ok {
		return nil, false
	}
	if meta.ContentHash != contentHash || meta.LibraryFingerprint != libraryFP || meta.ParseMode != mode {
		return nil, false
	}
	entry, ok := c.entries[key]
	if !ok || entry == nil || entry.Result == nil {
		return nil, false
	}
	return entry, true
}

func (c *ParseCache) meta(sourcePath string, mode parseModeFingerprint, libraryFP string) (cachedFileMeta, bool) {
	if c == nil {
		return cachedFileMeta{}, false
	}
	key := cacheEntryKey(sourcePath, mode, libraryFP)
	c.mu.RLock()
	defer c.mu.RUnlock()
	meta, ok := c.manifest.Files[key]
	return meta, ok
}

func (c *ParseCache) missingPaths(existing map[string]struct{}) []string {
	if c == nil {
		return nil
	}
	c.mu.RLock()
	defer c.mu.RUnlock()

	missing := make([]string, 0)
	seen := make(map[string]struct{})
	for _, meta := range c.manifest.Files {
		if _, ok := existing[meta.SourcePath]; ok {
			continue
		}
		if _, ok := seen[meta.SourcePath]; ok {
			continue
		}
		seen[meta.SourcePath] = struct{}{}
		missing = append(missing, meta.SourcePath)
	}
	sort.Strings(missing)
	return missing
}

func (c *ParseCache) put(unit *cachedFileUnit) error {
	if c == nil || unit == nil {
		return nil
	}
	key := cacheEntryKey(unit.SourcePath, unit.ParseMode, unit.LibraryFingerprint)
	meta := cachedFileMeta{
		SourcePath:         unit.SourcePath,
		ContentHash:        unit.ContentHash,
		ExportSignature:    unit.ExportSignature,
		EntryFile:          cacheEntryFile(key),
		ParseMode:          unit.ParseMode,
		LibraryFingerprint: unit.LibraryFingerprint,
		TopLevelPackages:   append([]string(nil), unit.TopLevelPackages...),
		Imports:            append([]cachedImport(nil), unit.Imports...),
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = unit
	c.manifest.Files[key] = meta

	if !c.persistent {
		return nil
	}
	entryData, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(c.dir, cacheFilesDirectory, meta.EntryFile), entryData, 0o644); err != nil {
		return err
	}
	return c.persistManifestLocked()
}

func (c *ParseCache) removeMissingPaths(existing map[string]struct{}) error {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	changed := false
	for key, meta := range c.manifest.Files {
		if _, ok := existing[meta.SourcePath]; ok {
			continue
		}
		delete(c.manifest.Files, key)
		delete(c.entries, key)
		_ = os.Remove(filepath.Join(c.dir, cacheFilesDirectory, meta.EntryFile))
		changed = true
	}
	for path := range c.manifest.ReverseDeps {
		if _, ok := existing[path]; !ok {
			delete(c.manifest.ReverseDeps, path)
			changed = true
		}
	}
	if changed {
		return c.persistManifestLocked()
	}
	return nil
}

func (c *ParseCache) reverseDependencyClosure(seeds map[string]struct{}) map[string]struct{} {
	out := make(map[string]struct{}, len(seeds))
	if c == nil {
		return out
	}
	c.mu.RLock()
	defer c.mu.RUnlock()

	queue := make([]string, 0, len(seeds))
	for seed := range seeds {
		out[seed] = struct{}{}
		queue = append(queue, seed)
	}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, dep := range c.manifest.ReverseDeps[current] {
			if _, seen := out[dep]; seen {
				continue
			}
			out[dep] = struct{}{}
			queue = append(queue, dep)
		}
	}
	return out
}

func (c *ParseCache) replaceManifest(units []*fileUnit, mode parseModeFingerprint, libraryFP string) error {
	if c == nil {
		return nil
	}
	reverseDeps := buildReverseDeps(units)

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.manifest == nil {
		c.manifest = &cacheManifest{}
	}
	c.manifest.Version = parseCacheVersion
	c.manifest.GrammarVersion = grammarVersion
	c.manifest.ReverseDeps = reverseDeps
	if c.manifest.Files == nil {
		c.manifest.Files = make(map[string]cachedFileMeta)
	}
	for _, unit := range units {
		entryKey := cacheEntryKey(unit.Source, mode, libraryFP)
		meta := c.manifest.Files[entryKey]
		meta.SourcePath = unit.Source
		meta.ExportSignature = unit.ExportSignature
		meta.ParseMode = mode
		meta.LibraryFingerprint = libraryFP
		meta.TopLevelPackages = append([]string(nil), unit.TopLevelPkgs...)
		meta.Imports = append([]cachedImport(nil), unit.Imports...)
		if entry, ok := c.entries[entryKey]; ok && entry != nil {
			meta.ContentHash = entry.ContentHash
			meta.EntryFile = cacheEntryFile(entryKey)
		}
		c.manifest.Files[entryKey] = meta
	}
	return c.persistManifestLocked()
}

func buildFileUnit(result *ParseResult, source string) *fileUnit {
	if result == nil {
		return &fileUnit{Source: source}
	}
	return &fileUnit{
		Source:          source,
		Result:          result,
		ExportSignature: computeExportSignature(result.Model),
		TopLevelPkgs:    topLevelPackageNames(result.Model),
		Imports:         collectImports(result.Model),
	}
}

func buildRepoIndex(units []*fileUnit, reg *LibraryRegistry) *RepoIndex {
	idx := &RepoIndex{
		elementIndex:     make(map[string]Element),
		shortNameIndex:   make(map[string][]Element),
		libraryRegistry:  reg,
		elementOwnerPath: make(map[string]string),
	}
	for _, unit := range units {
		if unit == nil || unit.Result == nil || unit.Result.Model == nil {
			continue
		}
		unit.Result.Model.Walk(func(elem Element) bool {
			qn := elem.QualifiedName()
			if qn != "" {
				idx.elementIndex[qn] = elem
				idx.elementOwnerPath[qn] = unit.Source
			}
			if snElem, ok := elem.(interface{ DeclaredShortName() string }); ok {
				sn := snElem.DeclaredShortName()
				if sn != "" {
					idx.shortNameIndex[sn] = append(idx.shortNameIndex[sn], elem)
				}
			}
			return true
		})
	}
	return idx
}

func topLevelPackageNames(model *Model) []string {
	if model == nil {
		return nil
	}
	names := make([]string, 0, len(model.Packages()))
	for _, pkg := range model.Packages() {
		if pkg == nil || pkg.Name() == "" {
			continue
		}
		names = append(names, pkg.Name())
	}
	sort.Strings(names)
	return names
}

func collectImports(model *Model) []cachedImport {
	if model == nil {
		return nil
	}
	imports := make([]cachedImport, 0)
	model.Walk(func(elem Element) bool {
		imp, ok := elem.(*Import)
		if !ok || imp == nil {
			return true
		}
		imports = append(imports, cachedImport{
			Namespace: imp.ImportedNamespace,
			IsAll:     imp.IsAll,
			IsRec:     imp.IsRecursive,
			Filters:   append([]string(nil), imp.FilterExpressions...),
		})
		return true
	})
	sort.Slice(imports, func(i, j int) bool {
		return imports[i].Namespace < imports[j].Namespace
	})
	return imports
}

func computeExportSignature(model *Model) string {
	if model == nil {
		return ""
	}
	lines := make([]string, 0, 64)
	for _, pkg := range model.Packages() {
		if pkg != nil && pkg.Name() != "" {
			lines = append(lines, "pkg:"+pkg.Name())
		}
	}
	model.Walk(func(elem Element) bool {
		if elem == nil {
			return true
		}
		if qn := elem.QualifiedName(); qn != "" {
			lines = append(lines, "qn:"+qn)
		}
		if snElem, ok := elem.(interface{ DeclaredShortName() string }); ok {
			if sn := snElem.DeclaredShortName(); sn != "" {
				lines = append(lines, "short:"+elem.QualifiedName()+"="+sn)
			}
		}
		if imp, ok := elem.(*Import); ok && imp != nil && imp.ImportedNamespace != "" {
			lines = append(lines, "import:"+imp.ImportedNamespace)
		}
		return true
	})
	sort.Strings(lines)
	sum := sha256.Sum256([]byte(strings.Join(lines, "\n")))
	return fmt.Sprintf("%x", sum[:])
}

func sharesTopLevelPackage(a, b []string) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}
	set := make(map[string]struct{}, len(a))
	for _, pkg := range a {
		set[pkg] = struct{}{}
	}
	for _, pkg := range b {
		if _, ok := set[pkg]; ok {
			return true
		}
	}
	return false
}

func importsTouchPackage(imports []cachedImport, pkg string) bool {
	for _, imp := range imports {
		ns := strings.TrimSpace(imp.Namespace)
		if ns == "" {
			continue
		}
		ns = strings.TrimSuffix(ns, "::**")
		ns = strings.TrimSuffix(ns, "::*")
		if ns == pkg || strings.HasPrefix(ns, pkg+"::") {
			return true
		}
	}
	return false
}

func buildReverseDeps(units []*fileUnit) map[string][]string {
	reverse := make(map[string][]string)
	for _, unit := range units {
		if unit == nil {
			continue
		}
		if _, ok := reverse[unit.Source]; !ok {
			reverse[unit.Source] = nil
		}
	}

	for _, a := range units {
		if a == nil {
			continue
		}
		for _, b := range units {
			if b == nil || a.Source == b.Source {
				continue
			}

			depends := sharesTopLevelPackage(a.TopLevelPkgs, b.TopLevelPkgs)
			if !depends {
				for _, pkg := range b.TopLevelPkgs {
					if importsTouchPackage(a.Imports, pkg) {
						depends = true
						break
					}
				}
			}
			if !depends {
				continue
			}
			reverse[b.Source] = append(reverse[b.Source], a.Source)
		}
	}

	for path := range reverse {
		if len(reverse[path]) == 0 {
			continue
		}
		sort.Strings(reverse[path])
		reverse[path] = compactStrings(reverse[path])
	}
	return reverse
}

func compactStrings(in []string) []string {
	if len(in) < 2 {
		return in
	}
	out := in[:1]
	for _, s := range in[1:] {
		if s == out[len(out)-1] {
			continue
		}
		out = append(out, s)
	}
	return out
}

func parseFileUnitBytes(ctx context.Context, input []byte, source string, opts ...ParseOption) (*fileUnit, *ParseResult) {
	result := parseBytesWithSourceContext(ctx, input, source, opts...)
	return buildFileUnit(result, source), result
}

func parseBytesWithCache(ctx context.Context, cfg *parseConfig, input []byte, source, cacheSource string, opts ...ParseOption) *ParseResult {
	if ctx == nil {
		ctx = context.Background()
	}
	registry := getOrCreateRegistry(cfg)
	mode := parseModeForConfig(cfg)
	libraryFP := libraryFingerprint(registry, cfg)
	contentHash := hashBytes(input)

	if entry, ok := cfg.cache.get(cacheSource, mode, libraryFP, contentHash); ok {
		if cfg.captureHash && entry.Result != nil {
			entry.Result.Hash = contentHash
		}
		if entry.Result != nil {
			entry.Result.Source = source
		}
		return entry.Result
	}

	result := parseBytesWithSourceContext(ctx, input, source, opts...)
	if result != nil && cfg.captureHash && result.Hash == "" {
		result.Hash = contentHash
	}
	unit := buildFileUnit(result, cacheSource)
	cacheEntry := &cachedFileUnit{
		SourcePath:         cacheSource,
		ContentHash:        contentHash,
		ParseMode:          mode,
		GrammarVersion:     grammarVersion,
		LibraryFingerprint: libraryFP,
		ExportSignature:    unit.ExportSignature,
		TopLevelPackages:   unit.TopLevelPkgs,
		Imports:            unit.Imports,
		Result:             result,
	}
	if result != nil && result.ParseError != nil {
		cacheEntry.Diagnostics = append(cacheEntry.Diagnostics, result.ParseError.Items...)
	}
	if result != nil {
		cacheEntry.Rewrites = append([]string(nil), result.Rewrites...)
	}
	_ = cfg.cache.put(cacheEntry)
	return result
}

func parseDirWithCache(ctx context.Context, dir string, cfg *parseConfig, opts []ParseOption) ([]*ParseResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	registry := getOrCreateRegistry(cfg)
	mode := parseModeForConfig(cfg)
	libraryFP := libraryFingerprint(registry, cfg)

	type fileContent struct {
		path string
		data []byte
		hash string
	}

	files := make([]fileContent, 0)
	existing := make(map[string]struct{})
	if err := walkModelFiles(ctx, dir, func(path string) bool {
		absPath, ok := parseCacheEligibleSource(path)
		if !ok {
			return true
		}
		data, err := os.ReadFile(absPath)
		if err != nil {
			files = append(files, fileContent{
				path: absPath,
				data: nil,
				hash: "",
			})
			return true
		}
		files = append(files, fileContent{
			path: absPath,
			data: data,
			hash: hashBytes(data),
		})
		existing[absPath] = struct{}{}
		return true
	}); err != nil {
		return nil, err
	}

	deletedPaths := cfg.cache.missingPaths(existing)
	if err := cfg.cache.removeMissingPaths(existing); err != nil {
		return nil, err
	}

	results := make([]*ParseResult, 0, len(files))
	units := make([]*fileUnit, 0, len(files))
	entriesByPath := make(map[string]*cachedFileUnit, len(files))
	changedPaths := make(map[string]struct{})
	exportChangedSeeds := make(map[string]struct{})

	parseOpts := append([]ParseOption(nil), opts...)
	if cfg.resolveModel {
		parseOpts = append(parseOpts, WithoutResolution())
	}

	for _, file := range files {
		if file.data == nil {
			results = append(results, walkErrorResult(file.path, fmt.Errorf("reading %s failed", file.path)))
			continue
		}
		oldMeta, hadOldMeta := cfg.cache.meta(file.path, mode, libraryFP)
		entry, hit := cfg.cache.get(file.path, mode, libraryFP, file.hash)
		if hit {
			entriesByPath[file.path] = entry
			unit := buildFileUnit(entry.Result, file.path)
			unit.ExportSignature = entry.ExportSignature
			unit.TopLevelPkgs = append([]string(nil), entry.TopLevelPackages...)
			unit.Imports = append([]cachedImport(nil), entry.Imports...)
			units = append(units, unit)
			results = append(results, entry.Result)
			continue
		}

		changedPaths[file.path] = struct{}{}
		unit, result := parseFileUnitBytes(ctx, file.data, file.path, parseOpts...)
		if result != nil && cfg.captureHash && result.Hash == "" {
			result.Hash = file.hash
		}
		cacheEntry := &cachedFileUnit{
			SourcePath:         file.path,
			ContentHash:        file.hash,
			ParseMode:          mode,
			GrammarVersion:     grammarVersion,
			LibraryFingerprint: libraryFP,
			ExportSignature:    unit.ExportSignature,
			TopLevelPackages:   unit.TopLevelPkgs,
			Imports:            unit.Imports,
			Result:             result,
		}
		if result != nil && result.ParseError != nil {
			cacheEntry.Diagnostics = append(cacheEntry.Diagnostics, result.ParseError.Items...)
		}
		if result != nil {
			cacheEntry.Rewrites = append([]string(nil), result.Rewrites...)
		}
		entriesByPath[file.path] = cacheEntry
		units = append(units, unit)
		results = append(results, result)
		if !hadOldMeta || oldMeta.ExportSignature != unit.ExportSignature {
			exportChangedSeeds[file.path] = struct{}{}
		}
	}

	for _, deleted := range deletedPaths {
		exportChangedSeeds[deleted] = struct{}{}
	}

	if cfg.resolveModel && cfg.buildModel {
		dirtySeeds := make(map[string]struct{}, len(exportChangedSeeds)+len(changedPaths))
		for path := range exportChangedSeeds {
			dirtySeeds[path] = struct{}{}
		}
		for path := range changedPaths {
			dirtySeeds[path] = struct{}{}
		}
		dirty := cfg.cache.reverseDependencyClosure(dirtySeeds)
		idx := buildRepoIndex(units, registry)
		for _, unit := range units {
			if unit == nil || unit.Result == nil || unit.Result.Model == nil {
				continue
			}
			entry := entriesByPath[unit.Source]
			if entry == nil {
				continue
			}
			unit.Result.Model.applyRepoIndex(idx)
			_, needsResolve := dirty[unit.Source]
			if !needsResolve && entry.ResolvedForRepo {
				continue
			}
			unit.Result.Model.ResolveReferencesWithIndex(idx)
			entry.ResolvedForRepo = true
		}
	}

	for _, unit := range units {
		if unit == nil {
			continue
		}
		if entry := entriesByPath[unit.Source]; entry != nil {
			entry.ExportSignature = unit.ExportSignature
			entry.TopLevelPackages = append([]string(nil), unit.TopLevelPkgs...)
			entry.Imports = append([]cachedImport(nil), unit.Imports...)
			if err := cfg.cache.put(entry); err != nil {
				return nil, err
			}
		}
	}
	if err := cfg.cache.replaceManifest(units, mode, libraryFP); err != nil {
		return nil, err
	}
	return results, nil
}
