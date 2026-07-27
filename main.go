package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ---- 7x5 bitmap font (1 = paint a contribution, 0 = skip) ----
// Each letter is 5 columns wide. We add one empty column as a spacer.
var font = map[rune][]string{
	'A': {
		"01110",
		"10001",
		"10001",
		"11111",
		"10001",
		"10001",
		"10001",
	},
	'B': {
		"11110",
		"10001",
		"10001",
		"11110",
		"10001",
		"10001",
		"11110",
	},
	'C': {
		"01111",
		"10000",
		"10000",
		"10000",
		"10000",
		"10000",
		"01111",
	},
	'D': {
		"11110",
		"10001",
		"10001",
		"10001",
		"10001",
		"10001",
		"11110",
	},
	'E': {
		"11111",
		"10000",
		"10000",
		"11110",
		"10000",
		"10000",
		"11111",
	},
	'F': {
		"11111",
		"10000",
		"10000",
		"11110",
		"10000",
		"10000",
		"10000",
	},
	'G': {
		"01111",
		"10000",
		"10000",
		"10011",
		"10001",
		"10001",
		"01111",
	},
	'H': {
		"10001",
		"10001",
		"10001",
		"11111",
		"10001",
		"10001",
		"10001",
	},
	'I': {
		"11111",
		"00100",
		"00100",
		"00100",
		"00100",
		"00100",
		"11111",
	},
	'J': {
		"00111",
		"00010",
		"00010",
		"00010",
		"00010",
		"10010",
		"01100",
	},
	'K': {
		"10001",
		"10010",
		"10100",
		"11000",
		"10100",
		"10010",
		"10001",
	},
	'L': {
		"10000",
		"10000",
		"10000",
		"10000",
		"10000",
		"10000",
		"11111",
	},
	'M': {
		"10001",
		"11011",
		"10101",
		"10101",
		"10001",
		"10001",
		"10001",
	},
	'N': {
		"10001",
		"11001",
		"10101",
		"10011",
		"10001",
		"10001",
		"10001",
	},
	'O': {
		"01110",
		"10001",
		"10001",
		"10001",
		"10001",
		"10001",
		"01110",
	},
	'P': {
		"11110",
		"10001",
		"10001",
		"11110",
		"10000",
		"10000",
		"10000",
	},
	'Q': {
		"01110",
		"10001",
		"10001",
		"10001",
		"10101",
		"10010",
		"01101",
	},
	'R': {
		"11110",
		"10001",
		"10001",
		"11110",
		"10100",
		"10010",
		"10001",
	},
	'S': {
		"01110",
		"10001",
		"10000",
		"01110",
		"00001",
		"10001",
		"01110",
	},
	'T': {
		"11111",
		"00100",
		"00100",
		"00100",
		"00100",
		"00100",
		"00100",
	},
	'U': {
		"10001",
		"10001",
		"10001",
		"10001",
		"10001",
		"10001",
		"01110",
	},
	'V': {
		"10001",
		"10001",
		"10001",
		"10001",
		"10001",
		"01010",
		"00100",
	},
	'W': {
		"10001",
		"10001",
		"10001",
		"10101",
		"10101",
		"11011",
		"10001",
	},
	'X': {
		"10001",
		"10001",
		"01010",
		"00100",
		"01010",
		"10001",
		"10001",
	},
	'Y': {
		"10001",
		"10001",
		"01010",
		"00100",
		"00100",
		"00100",
		"00100",
	},
	'Z': {
		"11111",
		"00001",
		"00010",
		"00100",
		"01000",
		"10000",
		"11111",
	},
	' ': {
		"00000",
		"00000",
		"00000",
		"00000",
		"00000",
		"00000",
		"00000",
	},
}

// Letter width = 5, plus 1 spacer column = 6 per letter.
const (
	letterWidth  = 5
	letterSpacer = 1
	cellWidth    = letterWidth + letterSpacer
)

// ---- CLI flags ----
var (
	pattern       = flag.String("pattern", "SARIF", "Text to draw on the heatmap")
	startDate     = flag.String("start", "", "Start date (YYYY-MM-DD). Default = next Monday")
	commitsPerDay = flag.Int("commits", 3, "Number of commits per filled cell (1-10)")
	dryRun        = flag.Bool("dry-run", false, "Preview the heatmap without creating commits")
)

func main() {
	flag.Parse()

	// Validate CLI flags
	trimmedPattern := strings.TrimSpace(*pattern)
	if trimmedPattern == "" {
		fmt.Fprintln(os.Stderr, "ERROR: pattern string cannot be empty")
		os.Exit(1)
	}

	if *commitsPerDay < 1 || *commitsPerDay > 10 {
		fmt.Fprintln(os.Stderr, "ERROR: commits flag must be between 1 and 10")
		os.Exit(1)
	}

	// 1. Determine start date (next Monday by default)
	start, err := getStartDate(*startDate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: invalid start date: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Starting from: %s\n", start.Format("2006-01-02"))

	// 2. Build the 7‑row bitmap grid for the whole timeline
	grid := buildGrid(trimmedPattern)

	// 3. Preview the heatmap in the terminal
	previewGrid(grid, start)

	if *dryRun {
		fmt.Println("\n[DRY RUN] No commits were created.")
		return
	}

	// 4. Generate the actual commits
	if err := generateCommits(grid, start, *commitsPerDay); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ All commits created! Push to GitHub to see the heatmap.")
}

// ---- helpers ----

// getStartDate parses a YYYY-MM-DD string or returns next Monday in UTC.
func getStartDate(s string) (time.Time, error) {
	if s != "" {
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			return time.Time{}, fmt.Errorf("must be in YYYY-MM-DD format: %w", err)
		}
		return t.UTC(), nil
	}

	// Default: next Monday in UTC
	now := time.Now().UTC()
	weekday := now.Weekday()
	daysUntilMonday := (8 - int(weekday)) % 7
	if daysUntilMonday == 0 {
		daysUntilMonday = 7
	}
	y, m, d := now.AddDate(0, 0, daysUntilMonday).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC), nil
}

// buildGrid turns the pattern string into a 7 x (weeks) grid of booleans.
func buildGrid(text string) [][]bool {
	upperText := strings.ToUpper(text)
	totalCols := len(upperText) * cellWidth
	rows := 7

	grid := make([][]bool, rows)
	for r := 0; r < rows; r++ {
		grid[r] = make([]bool, totalCols)
	}

	for charIdx, ch := range upperText {
		letterRows, ok := font[ch]
		if !ok {
			// Skip unsupported characters (creates blank spacer columns)
			continue
		}
		offset := charIdx * cellWidth

		for r := 0; r < rows && r < len(letterRows); r++ {
			rowStr := letterRows[r]
			for c := 0; c < letterWidth && c < len(rowStr); c++ {
				if rowStr[c] == '1' {
					grid[r][offset+c] = true
				}
			}
		}
	}
	return grid
}

// previewGrid prints the contribution heatmap using ANSI colours.
func previewGrid(grid [][]bool, start time.Time) {
	fmt.Println("\n📅 Heatmap preview (7 rows = Sun–Sat):")
	fmt.Println()

	dayLabels := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	for r := 0; r < 7; r++ {
		fmt.Printf("%s ", dayLabels[r])

		for c := 0; c < len(grid[r]); c++ {
			if grid[r][c] {
				fmt.Printf("\x1b[42m  \x1b[49m") // bright green block
			} else {
				fmt.Printf("\x1b[40m  \x1b[49m") // dark grey block
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// generateCommits walks through the grid and creates backdated commits.
func generateCommits(grid [][]bool, start time.Time, commitsPerCell int) error {
	// Ensure we are in a git repo (or init one)
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		fmt.Println("Initialising new Git repository...")
		if err := exec.Command("git", "init").Run(); err != nil {
			return fmt.Errorf("git init failed: %w", err)
		}
	}

	// Align start date to Sunday (start of week 0 in GitHub heatmap)
	startSunday := start.AddDate(0, 0, -int(start.Weekday()))

	cols := len(grid[0]) // total number of week columns

	// Walk week by week, day by day (Sunday to Saturday)
	for c := 0; c < cols; c++ {
		for r := 0; r < 7; r++ {
			if !grid[r][c] {
				continue
			}
			currentDate := startSunday.AddDate(0, 0, c*7+r)
			for i := 0; i < commitsPerCell; i++ {
				if err := createCommit(currentDate); err != nil {
					return fmt.Errorf("commit failed for %s: %w", currentDate.Format("2006-01-02"), err)
				}
			}
			fmt.Print(".")
		}
	}
	fmt.Println() // newline after progress dots
	return nil
}

// createCommit makes a single empty commit with a backdated timestamp.
func createCommit(date time.Time) error {
	dateStr := fmt.Sprintf("%04d-%02d-%02dT12:00:00", date.Year(), date.Month(), date.Day())

	cmd := exec.Command("git", "commit", "--allow-empty", "-m", "contribution-art")
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_DATE="+dateStr,
		"GIT_COMMITTER_DATE="+dateStr,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(output))
	}
	return nil
}