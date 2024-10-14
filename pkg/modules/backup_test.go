package modules

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBackup(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("Backup a valid file", func(t *testing.T) {
		originalFile := filepath.Join(tempDir, "testfile.txt")
		expectedContent := "test content"

		// Create a test file
		err := os.WriteFile(originalFile, []byte(expectedContent), 0644)
		assert.NoError(t, err)

		// Call the Backup function
		err = Backup(originalFile)
		assert.NoError(t, err)

		// Read the directory contents and find the backup file
		files, err := os.ReadDir(tempDir)
		assert.NoError(t, err)

		var backupFile string
		for _, file := range files {
			if strings.HasPrefix(file.Name(), "testfile.txt.backup.") {
				backupFile = filepath.Join(tempDir, file.Name())
				break
			}
		}

		assert.NotEmpty(t, backupFile, "Backup file was not created")

		// Verify the content of the backup file
		content, err := os.ReadFile(backupFile)
		assert.NoError(t, err)
		assert.Equal(t, expectedContent, string(content))
	})

	t.Run("Backup a non-existent file", func(t *testing.T) {
		err := Backup(filepath.Join(tempDir, "nonexistent.txt"))
		assert.Error(t, err)
	})

	t.Run("Attempt to backup a directory", func(t *testing.T) {
		dirPath := filepath.Join(tempDir, "testdir")
		err := os.Mkdir(dirPath, 0755)
		assert.NoError(t, err)

		err = Backup(dirPath)
		expectedErr := "Cannot create backup of a folder: " + dirPath
		assert.EqualError(t, err, expectedErr)
	})
}
