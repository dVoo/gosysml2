package sysml

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"

	"github.com/dVoo/gosysml2/low"
)

const batmobileBenchmarkPath = "/tmp/Batmobile.sysml"

func loadBatmobileBenchmarkFile(tb testing.TB) []byte {
	tb.Helper()
	content, err := os.ReadFile(batmobileBenchmarkPath)
	if err != nil {
		tb.Skipf("benchmark sample not available at %s: %v", batmobileBenchmarkPath, err)
	}
	return content
}

func makeBatmobileBenchmarkDir(tb testing.TB, copies int) string {
	tb.Helper()
	content := loadBatmobileBenchmarkFile(tb)
	dir := tb.TempDir()
	for i := 0; i < copies; i++ {
		name := filepath.Join(dir, "batmobile_"+strconv.Itoa(i)+".sysml")
		if err := os.WriteFile(name, content, 0o644); err != nil {
			tb.Fatalf("writing benchmark file %s: %v", name, err)
		}
	}
	return dir
}

func BenchmarkBatmobileNormalize(b *testing.B) {
	input := string(loadBatmobileBenchmarkFile(b))
	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _ = normalizeUnsupportedRequirementSyntax(input)
	}
}

func BenchmarkBatmobileLowParseBytes(b *testing.B) {
	content := loadBatmobileBenchmarkFile(b)
	b.ReportAllocs()
	b.SetBytes(int64(len(content)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree, errs := low.ParseBytes(content)
		if errs.HasErrors() || tree == nil {
			b.Fatal("failed to parse benchmark file")
		}
	}
}

func BenchmarkBatmobileHighParseBytes(b *testing.B) {
	content := loadBatmobileBenchmarkFile(b)
	b.ReportAllocs()
	b.SetBytes(int64(len(content)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result := ParseBytes(content, batmobileBenchmarkPath)
		if err := result.Err(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBatmobileHighParseBytesDiscardTree(b *testing.B) {
	content := loadBatmobileBenchmarkFile(b)
	b.ReportAllocs()
	b.SetBytes(int64(len(content)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result := ParseBytes(content, batmobileBenchmarkPath, WithDiscardTree())
		if err := result.Err(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBatmobileHighParseBytesNoResolve(b *testing.B) {
	content := loadBatmobileBenchmarkFile(b)
	b.ReportAllocs()
	b.SetBytes(int64(len(content)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result := ParseBytes(content, batmobileBenchmarkPath, WithoutResolution())
		if err := result.Err(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBatmobileHighParseBytesNoModelBuild(b *testing.B) {
	content := loadBatmobileBenchmarkFile(b)
	b.ReportAllocs()
	b.SetBytes(int64(len(content)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result := ParseBytes(content, batmobileBenchmarkPath, WithoutModelBuild(), WithDiscardTree())
		if err := result.Err(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBatmobileModelBuildOnly(b *testing.B) {
	content := loadBatmobileBenchmarkFile(b)
	tree, errs := low.ParseBytes(content)
	if errs.HasErrors() || tree == nil {
		b.Fatal("failed to parse benchmark file")
	}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		model := buildModelWithOptions(tree, nil, newParseRewriteHints(), false)
		if model == nil {
			b.Fatal("failed to build model")
		}
	}
}

func BenchmarkBatmobileBuildIndexAndResolve(b *testing.B) {
	content := loadBatmobileBenchmarkFile(b)
	tree, errs := low.ParseBytes(content)
	if errs.HasErrors() || tree == nil {
		b.Fatal("failed to parse benchmark file")
	}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		model := buildModelWithOptions(tree, nil, newParseRewriteHints(), false)
		model.BuildIndexAndResolve()
	}
}

func BenchmarkBatmobileParseDir(b *testing.B) {
	dir := makeBatmobileBenchmarkDir(b, 64)
	workerCounts := []int{1, 2, runtime.NumCPU()}

	for _, workers := range workerCounts {
		b.Run("workers="+strconv.Itoa(workers), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				count := 0
				for result := range ParseDir(context.Background(), dir, DirOptions{
					Workers: workers,
					ParseOptions: []ParseOption{
						WithDiscardTree(),
					},
				}) {
					if err := result.Err(); err != nil {
						b.Fatal(err)
					}
					count++
				}
				if count != 64 {
					b.Fatalf("expected 64 parse results, got %d", count)
				}
			}
		})
	}
}
