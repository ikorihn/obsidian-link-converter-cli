package olconv

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertUnderDir_Integration(t *testing.T) {
	// Create temporary test directory
	tempDir := t.TempDir()

	// Copy test files to temp directory
	copyTestVault(t, "testdata/sample_vault", tempDir)

	// Run conversion
	err := LinkToWikilink(tempDir)
	require.NoError(t, err)

	// Verify index.md conversion
	indexContent, err := os.ReadFile(filepath.Join(tempDir, "index.md"))
	require.NoError(t, err)
	indexStr := string(indexContent)

	// Check that markdown links are converted to wikilinks
	assert.Contains(t, indexStr, "[[basic|Basic Note]]")
	assert.Contains(t, indexStr, "[[note with spaces|Note with Spaces]]")
	assert.Contains(t, indexStr, "[[sub1/samename|Sub Note 1]]")
	assert.Contains(t, indexStr, "[[sub2/samename|Sub Note 2]]")
	assert.Contains(t, indexStr, "[[important|Notes Directory]]")
	assert.Contains(t, indexStr, "[[basic|Click here]]")
	assert.Contains(t, indexStr, "[[special|See this note]]")

	// Check that external links are NOT converted
	assert.Contains(t, indexStr, "[GitHub](https://github.com)")
	assert.Contains(t, indexStr, "[Example](https://example.com)")

	// Check that code block links are NOT converted
	assert.Contains(t, indexStr, "[Code Link](should-not-convert.md)")

	// Check that code span links are NOT converted
	assert.Contains(t, indexStr, "`[inline code](not-converted.md)`")

	// Verify basic.md conversion
	basicContent, err := os.ReadFile(filepath.Join(tempDir, "basic.md"))
	require.NoError(t, err)
	basicStr := string(basicContent)

	assert.Contains(t, basicStr, "[[index]]")
	assert.Contains(t, basicStr, "[[note with spaces]]")
}

func TestReverseConvertUnderDir_Integration(t *testing.T) {
	// Create temporary test directory
	tempDir := t.TempDir()

	// Copy test files to temp directory
	copyTestVault(t, "testdata/sample_vault", tempDir)

	// Run reverse conversion on wikilinks file
	err := WikilinkToLink(tempDir)
	require.NoError(t, err)

	// Verify wikilinks.md conversion
	wikilinksContent, err := os.ReadFile(filepath.Join(tempDir, "wikilinks.md"))
	require.NoError(t, err)
	wikilinksStr := string(wikilinksContent)

	// Check that wikilinks are converted to markdown links
	assert.Contains(t, wikilinksStr, "[basic](basic.md)")
	assert.Contains(t, wikilinksStr, "[note with spaces](note with spaces.md)")
	assert.Contains(t, wikilinksStr, "[Basic Note Link](basic.md)")
	assert.Contains(t, wikilinksStr, "[First Same Name](sub1/samename.md)")
	assert.Contains(t, wikilinksStr, "[Second Same Name](sub2/samename.md)")
	assert.Contains(t, wikilinksStr, "[index](index.md)")

	// Check that external links are NOT converted
	assert.Contains(t, wikilinksStr, "[link](https://example.com)")

	// Check that code block wikilinks are NOT converted
	assert.Contains(t, wikilinksStr, "[[should-not-convert]]")

	// Check that code span wikilinks are NOT converted
	assert.Contains(t, wikilinksStr, "`[[inline-code]]`")
}

func TestFileMapping_Integration(t *testing.T) {
	// Test file mapping functionality
	files, err := ListMdFiles("testdata/sample_vault")
	require.NoError(t, err)

	filemap := FileListToMap(files)

	// Check that samename files are properly mapped
	samenames := filemap["samename"]
	assert.Len(t, samenames, 2)
	assert.Contains(t, samenames, "testdata/sample_vault/sub1/samename.md")
	assert.Contains(t, samenames, "testdata/sample_vault/sub2/samename.md")

	// Check unique files
	basics := filemap["basic"]
	assert.Len(t, basics, 1)
	assert.Contains(t, basics, "testdata/sample_vault/basic.md")

	// Check file with spaces
	noteWithSpaces := filemap["note with spaces"]
	assert.Len(t, noteWithSpaces, 1)
	assert.Contains(t, noteWithSpaces, "testdata/sample_vault/note with spaces.md")
}

func TestRoundTripConversion_Integration(t *testing.T) {
	// Create temporary test directory
	tempDir := t.TempDir()

	// Copy test files to temp directory
	copyTestVault(t, "testdata/sample_vault", tempDir)

	// Read original content
	originalContent, err := os.ReadFile(filepath.Join(tempDir, "index.md"))
	require.NoError(t, err)

	// Convert markdown links to wikilinks
	err = LinkToWikilink(tempDir)
	require.NoError(t, err)

	// Convert wikilinks back to markdown links
	err = WikilinkToLink(tempDir)
	require.NoError(t, err)

	// Read final content
	finalContent, err := os.ReadFile(filepath.Join(tempDir, "index.md"))
	require.NoError(t, err)

	// The content should be similar (not exactly the same due to formatting differences)
	// but should contain the same essential links
	originalStr := string(originalContent)
	finalStr := string(finalContent)

	// Count markdown links in both
	originalMdLinks := strings.Count(originalStr, "](")
	finalMdLinks := strings.Count(finalStr, "](")

	// Should have same number of markdown links
	assert.Equal(t, originalMdLinks, finalMdLinks)
}

func TestEdgeCases_Integration(t *testing.T) {
	// Create temporary test directory
	tempDir := t.TempDir()

	// Copy test files to temp directory
	copyTestVault(t, "testdata/sample_vault", tempDir)

	// Run conversion
	err := LinkToWikilink(tempDir)
	require.NoError(t, err)

	// Verify edge_cases.md conversion
	edgeCasesContent, err := os.ReadFile(filepath.Join(tempDir, "edge_cases.md"))
	require.NoError(t, err)
	edgeCasesStr := string(edgeCasesContent)

	// Check Japanese characters
	assert.Contains(t, edgeCasesStr, "[[日本語|日本語ファイル]]")

	// Check multiple links in one line
	assert.Contains(t, edgeCasesStr, "[[basic|this]]")
	assert.Contains(t, edgeCasesStr, "[[important|that]]")

	// Check that malformed links are not affected
	assert.Contains(t, edgeCasesStr, "[incomplete link")
	assert.Contains(t, edgeCasesStr, "[[incomplete wikilink")

	// Check empty links - they get converted to empty wikilinks
	assert.Contains(t, edgeCasesStr, "[[|]]")
	assert.Contains(t, edgeCasesStr, "[[]]")
}

func TestJapaneseFiles_Integration(t *testing.T) {
	// Create temporary test directory
	tempDir := t.TempDir()

	// Copy test files to temp directory
	copyTestVault(t, "testdata/sample_vault", tempDir)

	// Run conversion
	err := LinkToWikilink(tempDir)
	require.NoError(t, err)

	// Verify Japanese file conversion
	japaneseContent, err := os.ReadFile(filepath.Join(tempDir, "日本語.md"))
	require.NoError(t, err)
	japaneseStr := string(japaneseContent)

	// Check that links are converted properly
	assert.Contains(t, japaneseStr, "[[index|インデックス]]")
	assert.Contains(t, japaneseStr, "[[basic|基本ノート]]")
}

func TestEmptyFiles_Integration(t *testing.T) {
	// Create temporary test directory
	tempDir := t.TempDir()

	// Create an empty file
	emptyFile := filepath.Join(tempDir, "empty.md")
	err := os.WriteFile(emptyFile, []byte(""), 0644)
	require.NoError(t, err)

	// Run conversion (should not crash on empty files)
	err = LinkToWikilink(tempDir)
	require.NoError(t, err)

	// File should still be empty
	content, err := os.ReadFile(emptyFile)
	require.NoError(t, err)
	assert.Empty(t, content)
}

// Helper function to copy test vault to temporary directory
func copyTestVault(t *testing.T, src, dst string) {
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		err = os.MkdirAll(filepath.Dir(dstPath), 0755)
		if err != nil {
			return err
		}

		return os.WriteFile(dstPath, data, info.Mode())
	})
	require.NoError(t, err)
}
