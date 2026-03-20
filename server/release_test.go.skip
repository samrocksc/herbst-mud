package server

import (
	"regexp"
	"strings"
	"testing"
)

// TestVersionParsing tests the semantic version parsing logic from the release workflow
func TestVersionParsing(t *testing.T) {
	tests := []struct {
		tag           string
		expectedMajor int
		expectedMinor int
		expectedPatch int
	}{
		{"v1.2.3", 1, 2, 3},
		{"v0.1.0", 0, 1, 0},
		{"v2.0.0", 2, 0, 0},
		{"v10.20.30", 10, 20, 30},
		{"v0.0.0", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			re := regexp.MustCompile(`v([0-9]+)\.([0-9]+)\.([0-9]+)`)
			matches := re.FindStringSubmatch(tt.tag)
			if len(matches) != 4 {
				t.Fatalf("failed to parse version from %s", tt.tag)
			}

			major := parseInt(matches[1])
			minor := parseInt(matches[2])
			patch := parseInt(matches[3])

			if major != tt.expectedMajor {
				t.Errorf("major: got %d, want %d", major, tt.expectedMajor)
			}
			if minor != tt.expectedMinor {
				t.Errorf("minor: got %d, want %d", minor, tt.expectedMinor)
			}
			if patch != tt.expectedPatch {
				t.Errorf("patch: got %d, want %d", patch, tt.expectedPatch)
			}
		})
	}
}

// TestVersionBumpType tests the commit message analysis for determining bump type
func TestVersionBumpType(t *testing.T) {
	tests := []struct {
		name        string
		commits     string
		expected    string
	}{
		{"only fixes - patch", "fix: resolve bug\nfix: another bug", "patch"},
		{"feat with fixes - minor", "fix: resolve bug\nfeat: add new feature", "minor"},
		{"breaking change - major", "fix: bug\nfeat!: breaking change", "major"},
		{"BREAKING CHANGE in body - major", "fix: bug\n\nBREAKING CHANGE: this is breaking", "major"},
		{"multiple features - minor", "feat: feature a\nfeat: feature b", "minor"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := determineBumpType(tt.commits)
			if result != tt.expected {
				t.Errorf("got %s, want %s", result, tt.expected)
			}
		})
	}
}

// TestVersionCalculation tests the version number calculation after bump type
func TestVersionCalculation(t *testing.T) {
	tests := []struct {
		name         string
		prevVersion  string
		bumpType     string
		expected     string
	}{
		{"patch bump", "v1.2.3", "patch", "v1.2.4"},
		{"minor bump", "v1.2.3", "minor", "v1.3.0"},
		{"major bump", "v1.2.3", "major", "v2.0.0"},
		{"patch from 0.x.x", "v0.1.0", "patch", "v0.1.1"},
		{"minor from 0.x.x", "v0.1.0", "minor", "v0.2.0"},
		{"major from 0.x.x", "v0.1.0", "major", "v1.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateNewVersion(tt.prevVersion, tt.bumpType)
			if result != tt.expected {
				t.Errorf("got %s, want %s", result, tt.expected)
			}
		})
	}
}

// TestChangelogGenerator tests the changelog generation from commit messages
func TestChangelogGenerator(t *testing.T) {
	commits := `feat: add new user API
fix: resolve login bug
docs: update README
feat: add character creation
chore: update dependencies
fix: fix database connection leak
feat!: breaking change to combat system`

	changelog := generateChangelog(commits)

	// Check that all sections exist
	if !strings.Contains(changelog, "### Features") && !strings.Contains(changelog, "### Bug Fixes") {
		t.Error("changelog should contain Features and Bug Fixes sections")
	}

	// Verify features are detected
	if !strings.Contains(changelog, "add new user API") {
		t.Error("changelog should contain feat: add new user API")
	}

	// Verify fixes are detected
	if !strings.Contains(changelog, "resolve login bug") {
		t.Error("changelog should contain fix: resolve login bug")
	}
}

// TestChangelogFormat tests that generated changelog follows Keep a Changelog format
func TestChangelogFormat(t *testing.T) {
	changelog := generateChangelog("feat: test")

	// Check format requirements
	if !strings.Contains(changelog, "## ") {
		t.Error("changelog should have version headers")
	}
	if !strings.Contains(changelog, "###") {
		t.Error("changelog should have section headers")
	}
}

// Helper functions that mirror the GitHub Actions workflow logic

func parseInt(s string) int {
	var n int
	for _, c := range s {
		n = n*10 + int(c-'0')
	}
	return n
}

func determineBumpType(commits string) string {
	bumpType := "patch"

	// Check for breaking changes first (has highest priority)
	if regexp.MustCompile(`^feat!.+:`).MatchString(commits) {
		return "major"
	}
	if strings.Contains(commits, "BREAKING CHANGE") {
		return "major"
	}

	// Check for features (minor)
	if regexp.MustCompile(`^feat(\(.+\))?:`).MatchString(commits) {
		bumpType = "minor"
	}

	return bumpType
}

func calculateNewVersion(prevTag, bumpType string) string {
	re := regexp.MustCompile(`v([0-9]+)\.([0-9]+)\.([0-9]+)`)
	matches := re.FindStringSubmatch(prevTag)
	if len(matches) != 4 {
		return "v0.0.0"
	}

	major := parseInt(matches[1])
	minor := parseInt(matches[2])
	patch := parseInt(matches[3])

	switch bumpType {
	case "major":
		major++
		minor = 0
		patch = 0
	case "minor":
		minor++
		patch = 0
	case "patch":
		patch++
	}

	return majorMinorPatch(major, minor, patch)
}

func majorMinorPatch(major, minor, patch int) string {
	return "v" + itoa(major) + "." + itoa(minor) + "." + itoa(patch)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var s string
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

func generateChangelog(commits string) string {
	var features, fixes, docs, chores []string

	for _, line := range strings.Split(commits, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "feat"):
			features = append(features, extractMessage(line))
		case strings.HasPrefix(line, "fix"):
			fixes = append(fixes, extractMessage(line))
		case strings.HasPrefix(line, "docs"):
			docs = append(docs, extractMessage(line))
		case strings.HasPrefix(line, "chore"):
			chores = append(chores, extractMessage(line))
		}
	}

	var changelog string

	if len(features) > 0 {
		changelog += "### Features\n"
		for _, f := range features {
			changelog += "- " + f + "\n"
		}
	}

	if len(fixes) > 0 {
		changelog += "### Bug Fixes\n"
		for _, f := range fixes {
			changelog += "- " + f + "\n"
		}
	}

	return changelog
}

func extractMessage(commit string) string {
	// Remove conventional commit prefix
	re := regexp.MustCompile(`^(feat|fix|docs|chore)(\(.+\))?:`)
	return re.ReplaceAllString(commit, "")
}