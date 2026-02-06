package sysml

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// LibraryResolver defines the interface for resolving library imports.
type LibraryResolver interface {
	// ResolveImport resolves an import namespace to a library package.
	// Supports patterns like "ScalarValues", "ScalarValues::*", "ISQ::mass"
	ResolveImport(namespace string) (*Package, error)

	// FindElement finds an element by qualified name across all loaded libraries.
	// Returns nil if not found.
	FindElement(qualifiedName string) Element

	// RegisterLibrary parses and registers a library file at the given path.
	RegisterLibrary(path string) error
}

// LibraryRegistry manages loaded SysML standard libraries.
// It provides thread-safe access to library elements and supports
// resolving imports from user models to library definitions.
type LibraryRegistry struct {
	// libraries maps package names to their parsed models
	libraries map[string]*Model

	// elementIndex provides fast lookup of elements by qualified name
	elementIndex map[string]Element

	// libraryPaths contains search paths for library files
	libraryPaths []string

	// mu provides thread-safe access to the registry
	mu sync.RWMutex

	// loaded tracks whether standard libraries have been loaded
	loaded bool
}

// LibraryOption configures library registry behavior.
type LibraryOption func(*LibraryRegistry)

// WithLibraryPaths sets custom library search paths.
func WithLibraryPaths(paths ...string) LibraryOption {
	return func(r *LibraryRegistry) {
		r.libraryPaths = paths
	}
}

// NewLibraryRegistry creates a new library registry with optional configuration.
// Default library path is "./libraries/sysml.library" if not specified.
func NewLibraryRegistry(opts ...LibraryOption) *LibraryRegistry {
	registry := &LibraryRegistry{
		libraries:    make(map[string]*Model, 10),
		elementIndex: make(map[string]Element, 100),
		libraryPaths: []string{"./libraries/sysml.library"},
		loaded:       false,
	}

	for _, opt := range opts {
		opt(registry)
	}

	return registry
}

// DiscoverLibraries walks the library paths and finds all .sysml files.
// Returns absolute paths to all discovered library files.
func (r *LibraryRegistry) DiscoverLibraries(rootPath string) ([]string, error) {
	var files []string

	absRoot, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("resolving library path %q: %w", rootPath, err)
	}

	// Check if directory exists
	info, err := os.Stat(absRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return files, nil // Return empty slice if path doesn't exist
		}
		return nil, fmt.Errorf("checking library path %q: %w", absRoot, err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("library path %q is not a directory", absRoot)
	}

	// Walk directory tree
	err = filepath.WalkDir(absRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			// Log but continue walking
			return nil
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Check for .sysml or .kerml extension (case-insensitive)
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".sysml" || ext == ".kerml" {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walking library directory %q: %w", absRoot, err)
	}

	return files, nil
}

// LoadLibrary parses a library file and returns the model.
// All packages in the model are marked as IsLibrary = true.
func (r *LibraryRegistry) LoadLibrary(path string) (*Model, error) {
	result := ParseFile(path)

	if result.Errors != nil && result.Errors.HasErrors() {
		return nil, fmt.Errorf("parsing library file %q: %w", path, result.Errors)
	}

	if result.Model == nil {
		return nil, fmt.Errorf("library file %q produced no model", path)
	}

	// Mark all packages as library packages
	for _, pkg := range result.Model.Packages {
		markPackageAsLibrary(pkg)
	}

	// Build index for the library model
	result.Model.BuildIndex()

	return result.Model, nil
}

// markPackageAsLibrary recursively marks a package and all nested packages as library packages.
func markPackageAsLibrary(pkg *Package) {
	pkg.IsLibrary = true
	for _, child := range pkg.Packages() {
		markPackageAsLibrary(child)
	}
}

// RegisterLibrary parses and registers a library file.
// The library is added to the registry and its elements are indexed.
func (r *LibraryRegistry) RegisterLibrary(path string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	model, err := r.LoadLibrary(path)
	if err != nil {
		return err
	}

	// Register each package in the model
	for _, pkg := range model.Packages {
		pkgName := pkg.Name()
		if pkgName == "" {
			continue // Skip unnamed packages
		}

		// Store the model for this package
		r.libraries[pkgName] = model

		// Index all elements from this package
		model.Walk(func(elem Element) bool {
			qn := elem.QualifiedName()
			if qn != "" {
				r.elementIndex[qn] = elem
			}
			return true
		})
	}

	return nil
}

// RegisterStandardLibrary discovers and loads all library files from configured paths.
// Parse errors are logged but don't stop the loading process.
func (r *LibraryRegistry) RegisterStandardLibrary() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.loaded {
		return nil // Already loaded
	}

	var allErrors []error

	for _, libPath := range r.libraryPaths {
		files, err := r.discoverLibrariesUnlocked(libPath)
		if err != nil {
			allErrors = append(allErrors, err)
			continue
		}

		for _, file := range files {
			model, err := r.loadLibraryUnlocked(file)
			if err != nil {
				// Log but continue with other files
				allErrors = append(allErrors, fmt.Errorf("loading %q: %w", file, err))
				continue
			}

			// Register each package
			for _, pkg := range model.Packages {
				pkgName := pkg.Name()
				if pkgName == "" {
					continue
				}

				r.libraries[pkgName] = model

				// Index all elements
				model.Walk(func(elem Element) bool {
					qn := elem.QualifiedName()
					if qn != "" {
						r.elementIndex[qn] = elem
					}
					return true
				})
			}
		}
	}

	r.loaded = true

	if len(allErrors) > 0 {
		// Return first error but indicate partial success
		return fmt.Errorf("loaded with %d errors: %w", len(allErrors), allErrors[0])
	}

	return nil
}

// discoverLibrariesUnlocked is the internal version without locking.
func (r *LibraryRegistry) discoverLibrariesUnlocked(rootPath string) ([]string, error) {
	var files []string

	absRoot, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(absRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return files, nil
		}
		return nil, err
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("not a directory: %s", absRoot)
	}

	filepath.WalkDir(absRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if strings.ToLower(filepath.Ext(path)) == ".sysml" {
			files = append(files, path)
		}
		return nil
	})

	return files, nil
}

// loadLibraryUnlocked is the internal version without locking.
func (r *LibraryRegistry) loadLibraryUnlocked(path string) (*Model, error) {
	result := ParseFile(path)

	if result.Errors != nil && result.Errors.HasErrors() {
		return nil, result.Errors
	}

	if result.Model == nil {
		return nil, fmt.Errorf("no model produced")
	}

	for _, pkg := range result.Model.Packages {
		markPackageAsLibrary(pkg)
	}

	result.Model.BuildIndex()
	return result.Model, nil
}

// FindElement finds an element by qualified name across all loaded libraries.
// Thread-safe for concurrent access.
func (r *LibraryRegistry) FindElement(qualifiedName string) Element {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Fast path: direct lookup
	if elem := r.elementIndex[qualifiedName]; elem != nil {
		return elem
	}

	// Try partial qualified name resolution
	// e.g., if looking for "ScalarValues::Real", try searching through packages
	parts := strings.Split(qualifiedName, "::")
	if len(parts) < 2 {
		return nil
	}

	pkgName := parts[0]
	model := r.libraries[pkgName]
	if model == nil {
		return nil
	}

	// Search in the model's index
	return model.FindByQualifiedName(qualifiedName)
}

// ResolveImport resolves an import namespace to a library package.
// Supports patterns:
//   - Simple package: "ScalarValues"
//   - Wildcard: "ScalarValues::*"
//   - Specific element: "ISQ::mass"
//
// Thread-safe for concurrent access.
func (r *LibraryRegistry) ResolveImport(namespace string) (*Package, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Parse the namespace
	namespace = strings.TrimSpace(namespace)

	// Handle wildcard patterns
	isWildcard := strings.HasSuffix(namespace, "::")
	if isWildcard {
		namespace = strings.TrimSuffix(namespace, "::")
	}

	// Handle recursive wildcard (::**)
	isRecursive := strings.HasSuffix(namespace, "::**")
	if isRecursive {
		namespace = strings.TrimSuffix(namespace, "::**")
	}

	// Split by :: to get package path
	parts := strings.Split(namespace, "::")
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty namespace")
	}

	// Get the root package name
	rootName := parts[0]

	// Find the model containing this package
	model := r.libraries[rootName]
	if model == nil {
		return nil, fmt.Errorf("library package %q not found", rootName)
	}

	// Find the specific package
	var pkg *Package
	for _, p := range model.Packages {
		if p.Name() == rootName {
			pkg = p
			break
		}
	}

	if pkg == nil {
		return nil, fmt.Errorf("package %q not found in library", rootName)
	}

	// Navigate to nested package if needed
	for i := 1; i < len(parts); i++ {
		found := false
		for _, child := range pkg.Packages() {
			if child.Name() == parts[i] {
				pkg = child
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("nested package %q not found in %q", parts[i], pkg.QualifiedName())
		}
	}

	return pkg, nil
}

// GetLibraryCount returns the number of loaded library packages.
func (r *LibraryRegistry) GetLibraryCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.libraries)
}

// GetElementCount returns the number of indexed elements.
func (r *LibraryRegistry) GetElementCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.elementIndex)
}

// IsLoaded returns true if standard libraries have been loaded.
func (r *LibraryRegistry) IsLoaded() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.loaded
}

// GetLibraryNames returns a sorted list of loaded library package names.
func (r *LibraryRegistry) GetLibraryNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.libraries))
	for name := range r.libraries {
		names = append(names, name)
	}

	// Sort for consistent ordering
	sort.Strings(names)
	return names
}
