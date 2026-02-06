package low

import (
	"os"
	"path/filepath"
	"testing"
)

// testdataDir returns the path to testdata relative to the module root.
// The gosysml2 module is in a subdirectory, so testdata is at ../../testdata/.
func testdataDir(t testing.TB) string {
	t.Helper()
	dir := filepath.Join("..", "..", "testdata")
	if _, err := os.Stat(dir); err != nil {
		t.Skipf("testdata directory not found at %s: %v", dir, err)
	}
	return dir
}

// loadTestFile loads a single test file for benchmarking.
func loadTestFile(t testing.TB, name string) []byte {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(testdataDir(t), name))
	if err != nil {
		t.Fatalf("loading test file %s: %v", name, err)
	}
	return content
}

// loadAllTestFiles loads all .sysml files from testdata.
func loadAllTestFiles(t testing.TB) map[string][]byte {
	t.Helper()
	dir := testdataDir(t)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading testdata dir: %v", err)
	}
	files := make(map[string][]byte, len(entries))
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".sysml" {
			content, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				t.Fatalf("reading %s: %v", e.Name(), err)
			}
			files[e.Name()] = content
		}
	}
	return files
}

// loadAllTestFilesRecursive loads all .sysml files recursively from testdata.
func loadAllTestFilesRecursive(t testing.TB) map[string][]byte {
	t.Helper()
	dir := testdataDir(t)
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
			// Store relative path from testdata dir
			relPath, _ := filepath.Rel(dir, path)
			files[relPath] = content
		}
		return nil
	})

	if err != nil {
		t.Fatalf("walking testdata: %v", err)
	}
	return files
}

// BenchmarkLexSingle benchmarks lexing a single large file.
func BenchmarkLexSingle(b *testing.B) {
	// Use the largest test file
	content := loadTestFile(b, filepath.Join("Vehicle Example", "SysML v2 Spec Annex A SimpleVehicleModel.sysml"))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		lexer := NewLexerFromBytes(content)
		_ = lexer.AllTokens()
	}
}

// BenchmarkLexAll benchmarks lexing all test files sequentially.
func BenchmarkLexAll(b *testing.B) {
	files := loadAllTestFilesRecursive(b)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, content := range files {
			lexer := NewLexerFromBytes(content)
			_ = lexer.AllTokens()
		}
	}
}

// BenchmarkLexAllPerFile benchmarks lexing with sub-benchmarks per file.
func BenchmarkLexAllPerFile(b *testing.B) {
	files := loadAllTestFilesRecursive(b)

	for name, content := range files {
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				lexer := NewLexerFromBytes(content)
				_ = lexer.AllTokens()
			}
		})
	}
}

// BenchmarkParseSingle benchmarks parsing a single file to parse tree.
func BenchmarkParseSingle(b *testing.B) {
	content := loadTestFile(b, filepath.Join("Vehicle Example", "SysML v2 Spec Annex A SimpleVehicleModel.sysml"))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ParseBytes(content)
	}
}

// BenchmarkParseAll benchmarks parsing all test files sequentially.
func BenchmarkParseAll(b *testing.B) {
	files := loadAllTestFilesRecursive(b)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, content := range files {
			_, _ = ParseBytes(content)
		}
	}
}

// BenchmarkParseAllPerFile benchmarks parsing with sub-benchmarks per file.
func BenchmarkParseAllPerFile(b *testing.B) {
	files := loadAllTestFilesRecursive(b)

	for name, content := range files {
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = ParseBytes(content)
			}
		})
	}
}

// BenchmarkValidateSingle benchmarks validation without tree building.
func BenchmarkValidateSingle(b *testing.B) {
	content := loadTestFile(b, filepath.Join("Vehicle Example", "SysML v2 Spec Annex A SimpleVehicleModel.sysml"))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ValidateBytes(content)
	}
}

// BenchmarkValidateAll benchmarks validating all test files.
func BenchmarkValidateAll(b *testing.B) {
	files := loadAllTestFilesRecursive(b)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, content := range files {
			_ = ValidateBytes(content)
		}
	}
}

// BenchmarkValidateAllPerFile benchmarks validation with sub-benchmarks per file.
func BenchmarkValidateAllPerFile(b *testing.B) {
	files := loadAllTestFilesRecursive(b)

	for name, content := range files {
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = ValidateBytes(content)
			}
		})
	}
}

// BenchmarkLexerCreation benchmarks lexer creation overhead.
func BenchmarkLexerCreation(b *testing.B) {
	content := loadTestFile(b, filepath.Join("Vehicle Example", "SysML v2 Spec Annex A SimpleVehicleModel.sysml"))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = NewLexerFromBytes(content)
	}
}

// BenchmarkParserCreation benchmarks parser creation overhead.
func BenchmarkParserCreation(b *testing.B) {
	content := loadTestFile(b, filepath.Join("Vehicle Example", "SysML v2 Spec Annex A SimpleVehicleModel.sysml"))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = NewParserFromBytes(content)
	}
}
