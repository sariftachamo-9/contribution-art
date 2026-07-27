package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGetStartDate(t *testing.T) {
	t.Run("default next Monday", func(t *testing.T) {
		start, err := getStartDate("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if start.Weekday() != time.Monday {
			t.Errorf("expected weekday Monday, got %v", start.Weekday())
		}
	})

	t.Run("valid custom date", func(t *testing.T) {
		start, err := getStartDate("2025-06-01")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "2025-06-01"
		if start.Format("2006-01-02") != expected {
			t.Errorf("expected date %s, got %s", expected, start.Format("2006-01-02"))
		}
	})

	t.Run("invalid date format", func(t *testing.T) {
		_, err := getStartDate("06-01-2025")
		if err == nil {
			t.Errorf("expected error for invalid date format, got nil")
		}
	})
}

func TestBuildGrid(t *testing.T) {
	t.Run("grid dimensions and lowercase conversion", func(t *testing.T) {
		grid := buildGrid("sarif")
		if len(grid) != 7 {
			t.Fatalf("expected 7 rows, got %d", len(grid))
		}
		expectedCols := 5 * cellWidth // 5 letters * 6 = 30
		if len(grid[0]) != expectedCols {
			t.Fatalf("expected %d columns, got %d", expectedCols, len(grid[0]))
		}
	})

	t.Run("spacer columns are false", func(t *testing.T) {
		grid := buildGrid("A")
		// Letter 'A' occupies columns 0..4, column 5 is spacer
		for r := 0; r < 7; r++ {
			if grid[r][5] {
				t.Errorf("expected spacer column at index 5 to be false, got true at row %d", r)
			}
		}
	})
}

func TestGenerateCommits(t *testing.T) {
	// Create a temp directory for git commit testing
	tempDir := t.TempDir()
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current wd: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to chdir to tempDir: %v", err)
	}
	defer func() {
		_ = os.Chdir(origDir)
	}()

	// Configure git author/committer so git commit doesn't fail if git user config is missing
	_ = exec.Command("git", "config", "--global", "user.name", "Test User").Run()
	_ = exec.Command("git", "config", "--global", "user.email", "test@example.com").Run()

	// 1x1 simple test grid: only 1 day is true
	grid := make([][]bool, 7)
	for r := 0; r < 7; r++ {
		grid[r] = make([]bool, 1)
	}
	grid[1][0] = true // Monday of week 0

	startDate := time.Date(2025, 6, 2, 0, 0, 0, 0, time.UTC) // Mon Jun 2, 2025
	commitsPerCell := 2

	err = generateCommits(grid, startDate, commitsPerCell)
	if err != nil {
		t.Fatalf("generateCommits failed: %v", err)
	}

	// Verify git log has 2 commits
	out, err := exec.Command("git", "rev-list", "--count", "HEAD").Output()
	if err != nil {
		t.Fatalf("git rev-list failed: %v", err)
	}
	countStr := strings.TrimSpace(string(out))
	if countStr != "2" {
		t.Errorf("expected 2 commits, got %s", countStr)
	}
}
