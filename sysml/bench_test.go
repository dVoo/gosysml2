package sysml

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/dVoo/gosysml2_oc/low"
)

// testdataDir returns the path to testdata relative to the module root.
func testdataDir(t testing.TB) string {
	t.Helper()
	dir := filepath.Join("..", "..", "testdata")
	if _, err := os.Stat(dir); err != nil {
		t.Skipf("testdata directory not found at %s: %v", dir, err)
	}
	return dir
}

// loadTestFile loads a single test file for benchmarking.
func loadTestFile(b testing.TB, name string) []byte {
	b.Helper()
	content, err := os.ReadFile(filepath.Join(testdataDir(b), name))
	if err != nil {
		b.Fatalf("loading test file %s: %v", name, err)
	}
	return content
}

// loadAllTestFilesRecursive loads all .sysml files recursively from testdata.
func loadAllTestFilesRecursive(b testing.TB) map[string][]byte {
	b.Helper()
	dir := testdataDir(b)
	files := make(map[string][]byte)

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".sysml" {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			relPath, _ := filepath.Rel(dir, path)
			files[relPath] = content
		}
		return nil
	})

	if err != nil {
		b.Fatalf("walking testdata: %v", err)
	}
	return files
}

// BenchmarkParseStringSingle benchmarks full parse pipeline for single file.
func BenchmarkParseStringSingle(b *testing.B) {
	content := loadTestFile(b, filepath.Join("Vehicle Example", "SysML v2 Spec Annex A SimpleVehicleModel.sysml"))
	input := string(content)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ParseString(input)
	}
}

// BenchmarkParseStringAll benchmarks full parse pipeline for all files.
func BenchmarkParseStringAll(b *testing.B) {
	files := loadAllTestFilesRecursive(b)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, content := range files {
			_ = ParseString(string(content))
		}
	}
}

// BenchmarkParseFile benchmarks ParseFile including file I/O.
func BenchmarkParseFile(b *testing.B) {
	testFile := filepath.Join(testdataDir(b), "Vehicle Example", "SysML v2 Spec Annex A SimpleVehicleModel.sysml")
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ParseFile(testFile)
	}
}

// BenchmarkParseBytes benchmarks ParseBytes (more efficient than ParseString).
func BenchmarkParseBytes(b *testing.B) {
	content := loadTestFile(b, filepath.Join("Vehicle Example", "SysML v2 Spec Annex A SimpleVehicleModel.sysml"))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ParseBytes(content, "test")
	}
}

// BenchmarkModelBuild benchmarks model building by reusing parse tree.
func BenchmarkModelBuild(b *testing.B) {
	content := loadTestFile(b, filepath.Join("Vehicle Example", "SysML v2 Spec Annex A SimpleVehicleModel.sysml"))
	// Parse once to get the tree
	tree, _ := low.ParseBytes(content)
	if tree == nil {
		b.Fatal("failed to parse test file")
	}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = buildModel(tree, nil)
	}
}

// BenchmarkModelWalk benchmarks walking a parsed model.
func BenchmarkModelWalk(b *testing.B) {
	content := loadTestFile(b, filepath.Join("Vehicle Example", "SysML v2 Spec Annex A SimpleVehicleModel.sysml"))
	result := ParseBytes(content, "test")
	if result.Model == nil {
		b.Fatal("failed to parse test file")
	}
	model := result.Model
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		count := 0
		Walk(model, func(elem Element, depth int) bool {
			count++
			return true
		})
	}
}

// BenchmarkFindByKind benchmarks FindByKind performance.
func BenchmarkFindByKind(b *testing.B) {
	content := loadTestFile(b, filepath.Join("Vehicle Example", "SysML v2 Spec Annex A SimpleVehicleModel.sysml"))
	result := ParseBytes(content, "test")
	if result.Model == nil {
		b.Fatal("failed to parse test file")
	}
	model := result.Model
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = FindPackages(model)
	}
}

// BenchmarkFindAllGeneric benchmarks generic FindAll[T] performance.
func BenchmarkFindAllGeneric(b *testing.B) {
	content := loadTestFile(b, filepath.Join("Vehicle Example", "SysML v2 Spec Annex A SimpleVehicleModel.sysml"))
	result := ParseBytes(content, "test")
	if result.Model == nil {
		b.Fatal("failed to parse test file")
	}
	model := result.Model
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = FindAll[*Package](model)
	}
}

// BenchmarkIteratorAll benchmarks iter.Seq All() iterator.
func BenchmarkIteratorAll(b *testing.B) {
	content := loadTestFile(b, filepath.Join("Vehicle Example", "SysML v2 Spec Annex A SimpleVehicleModel.sysml"))
	result := ParseBytes(content, "test")
	if result.Model == nil {
		b.Fatal("failed to parse test file")
	}
	model := result.Model
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		count := 0
		for range All(model) {
			count++
		}
	}
}

// BenchmarkIteratorOfType benchmarks OfType[*Package] iterator.
func BenchmarkIteratorOfType(b *testing.B) {
	content := loadTestFile(b, filepath.Join("Vehicle Example", "SysML v2 Spec Annex A SimpleVehicleModel.sysml"))
	result := ParseBytes(content, "test")
	if result.Model == nil {
		b.Fatal("failed to parse test file")
	}
	model := result.Model
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		count := 0
		for range OfType[*Package](model) {
			count++
		}
	}
}

// BenchmarkIteratorOfKind benchmarks OfKind iterator.
func BenchmarkIteratorOfKind(b *testing.B) {
	content := loadTestFile(b, filepath.Join("Vehicle Example", "SysML v2 Spec Annex A SimpleVehicleModel.sysml"))
	result := ParseBytes(content, "test")
	if result.Model == nil {
		b.Fatal("failed to parse test file")
	}
	model := result.Model
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		count := 0
		for range OfKind(model, KindPackage) {
			count++
		}
	}
}

// BenchmarkParseDirectoryParallel benchmarks parallel parsing with different worker counts.
func BenchmarkParseDirectoryParallel(b *testing.B) {
	testDir := testdataDir(b)
	files := loadAllTestFilesRecursive(b)
	numFiles := len(files)

	b.Logf("Benchmarking parallel parsing with %d files", numFiles)

	// Benchmark with different worker counts
	workerCounts := []int{1, 2, 4, runtime.NumCPU()}

	for _, workers := range workerCounts {
		b.Run(fmt.Sprintf("workers=%d", workers), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := ParseDirectoryParallel(testDir, workers, WithDiscardTree())
				if err != nil {
					b.Fatalf("ParseDirectoryParallel failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkParseDirectorySequential benchmarks sequential parsing for comparison.
func BenchmarkParseDirectorySequential(b *testing.B) {
	testDir := testdataDir(b)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ParseDirectory(testDir)
		if err != nil {
			b.Fatalf("ParseDirectory failed: %v", err)
		}
	}
}

// BenchmarkParseDirectoryStream benchmarks streaming parse mode.
func BenchmarkParseDirectoryStream(b *testing.B) {
	testDir := testdataDir(b)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		count := 0
		err := ParseDirectoryStream(testDir, func(r *ParseResult) error {
			if r.Success() {
				count++
			}
			return nil
		})
		if err != nil {
			b.Fatalf("ParseDirectoryStream failed: %v", err)
		}
	}
}

// BenchmarkModelBuildIndex benchmarks index building.
func BenchmarkModelBuildIndex(b *testing.B) {
	content := loadTestFile(b, filepath.Join("Vehicle Example", "SysML v2 Spec Annex A SimpleVehicleModel.sysml"))
	result := ParseBytes(content, "test")
	if result.Model == nil {
		b.Fatal("failed to parse test file")
	}
	model := result.Model
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		model.BuildIndex()
	}
}

// BenchmarkModelResolveReferences benchmarks reference resolution.
func BenchmarkModelResolveReferences(b *testing.B) {
	content := loadTestFile(b, filepath.Join("Vehicle Example", "SysML v2 Spec Annex A SimpleVehicleModel.sysml"))
	result := ParseBytes(content, "test")
	if result.Model == nil {
		b.Fatal("failed to parse test file")
	}
	model := result.Model
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		model.ResolveReferences()
	}
}

// BenchmarkCountAll benchmarks CountAll function using iterators.
func BenchmarkCountAll(b *testing.B) {
	content := loadTestFile(b, filepath.Join("Vehicle Example", "SysML v2 Spec Annex A SimpleVehicleModel.sysml"))
	result := ParseBytes(content, "test")
	if result.Model == nil {
		b.Fatal("failed to parse test file")
	}
	model := result.Model
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = CountAll(model)
	}
}

// BenchmarkVisit benchmarks the visitor pattern.
func BenchmarkVisit(b *testing.B) {
	content := loadTestFile(b, filepath.Join("Vehicle Example", "SysML v2 Spec Annex A SimpleVehicleModel.sysml"))
	result := ParseBytes(content, "test")
	if result.Model == nil {
		b.Fatal("failed to parse test file")
	}
	model := result.Model
	visitor := &BaseVisitor{}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Visit(model, visitor)
	}
}

// BenchmarkParseWithDiscardTree benchmarks parsing with tree discarded.
func BenchmarkParseWithDiscardTree(b *testing.B) {
	content := loadTestFile(b, filepath.Join("Vehicle Example", "SysML v2 Spec Annex A SimpleVehicleModel.sysml"))
	input := string(content)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ParseString(input, WithDiscardTree())
	}
}
