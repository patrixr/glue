package modules_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/patrixr/glue/pkg/modules"
	"github.com/stretchr/testify/assert"
)

func createTestDir(t *testing.T, files map[string]string) string {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "copy-test")
	assert.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	for path, content := range files {
		filePath := filepath.Join(tempDir, path)
		dir := filepath.Dir(filePath)
		err := os.MkdirAll(dir, 0755)
		assert.NoError(t, err)

		if strings.HasPrefix(content, "link:") { // Check for link
			linkTarget := strings.TrimPrefix(content, "link:")
			err = os.Symlink(filepath.Join(tempDir, linkTarget), filePath) // Create symlink
			assert.NoError(t, err)

		} else {
			err = os.WriteFile(filePath, []byte(content), 0644) // Create regular file
			assert.NoError(t, err)
		}
	}
	return tempDir
}

func createSymlink(t *testing.T, src, dst string) {
	t.Helper()
	err := os.Symlink(src, dst)
	assert.NoError(t, err)
}

func TestCopy(t *testing.T) {
	testCases := []struct {
		name          string
		opts          modules.CopyOpts
		srcFiles      map[string]string
		wantFiles     map[string]string
		dontWantFiles []string
		wantErr       bool
	}{
		{
			name: "Copy Dir - Create",
			opts: modules.CopyOpts{Source: "src", Dest: "dest", Strategy: modules.StrategyMerge},
			srcFiles: map[string]string{
				"src/file1.txt":        "content1",
				"src/file2.txt":        "content2",
				"src/subdir/file3.txt": "content3",
			},
			wantFiles: map[string]string{
				"dest/file1.txt":        "content1",
				"dest/file2.txt":        "content2",
				"dest/subdir/file3.txt": "content3",
			},
		},
		{
			name: "Copy Dir - Replace",
			opts: modules.CopyOpts{Source: "src", Dest: "dest", Strategy: modules.StrategyReplace},
			srcFiles: map[string]string{
				"src/file1.txt":       "content1",
				"src/file2.txt":       "content2",
				"src/subdir/yay.txt":  "yay",
				"dest/subdir/nay.txt": "nay",
			},
			wantFiles: map[string]string{
				"dest/file1.txt":     "content1",
				"dest/file2.txt":     "content2",
				"src/subdir/yay.txt": "yay",
			},
			dontWantFiles: []string{"dest/subdir/nay.txt"},
		},
		{
			name: "Copy Dir - Replace",
			opts: modules.CopyOpts{Source: "src", Dest: "dest", Strategy: modules.StrategyMerge},
			srcFiles: map[string]string{
				"src/file1.txt":       "content1",
				"src/file2.txt":       "content2",
				"src/subdir/yay.txt":  "yay",
				"dest/subdir/nay.txt": "nay",
			},
			wantFiles: map[string]string{
				"dest/file1.txt":      "content1",
				"dest/file2.txt":      "content2",
				"src/subdir/yay.txt":  "yay",
				"dest/subdir/nay.txt": "nay",
			},
		},
		{
			name:     "Copy a file",
			opts:     modules.CopyOpts{Source: "file1.txt", Dest: "dest/file2.txt"},
			srcFiles: map[string]string{"file1.txt": "source content"},
			wantFiles: map[string]string{
				"dest/file2.txt": "source content",
			},
		},
		{
			name:    "Invalid source",
			opts:    modules.CopyOpts{Source: "i/dont/exist.txt", Dest: "."},
			wantErr: true,
		},
		{
			name: "Symlink Deep",
			opts: modules.CopyOpts{Source: "src", Dest: "dest", Symlink: modules.SymlinkDeep},
			srcFiles: map[string]string{
				"src/file.txt": "content",
				"src/link":     "link:src/file.txt",
			},
			wantFiles: map[string]string{
				"dest/file.txt": "content",
				"dest/link":     "content", // Should be the content of the linked file
			},
		},
		{
			name: "Symlink Shallow",
			opts: modules.CopyOpts{Source: "src", Dest: "dest", Symlink: modules.SymlinkShallow},
			srcFiles: map[string]string{
				"src/file.txt": "content",
				"src/link":     "link:src/file.txt",
			},
			wantFiles: map[string]string{
				"dest/file.txt": "content",
				"dest/link":     "link:src/file.txt", // Remains a symlink
			},
		},
		{
			name: "Symlink Skip",
			opts: modules.CopyOpts{Source: "src", Dest: "dest", Symlink: modules.SymlinkSkip},
			srcFiles: map[string]string{
				"src/file.txt": "content",
				"src/link":     "link:src/file.txt",
			},
			wantFiles: map[string]string{
				"dest/file.txt": "content",
			},
			dontWantFiles: []string{"dest/link"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := createTestDir(t, tc.srcFiles)

			tc.opts.Source = filepath.Join(tempDir, tc.opts.Source)
			tc.opts.Dest = filepath.Join(tempDir, tc.opts.Dest)

			err := modules.Copy(tc.opts)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			for path, wantContent := range tc.wantFiles {
				filePath := filepath.Join(tempDir, path)

				if strings.HasPrefix(wantContent, "link:") {
					gotLink, err := os.Readlink(filePath)
					assert.NoError(t, err)
					assert.Equal(t, filepath.Join(tempDir, strings.TrimPrefix(wantContent, "link:")), gotLink)
				} else {
					gotContent, err := os.ReadFile(filePath)
					assert.NoError(t, err)
					assert.Equal(t, wantContent, string(gotContent))
				}
			}

			for _, path := range tc.dontWantFiles {
				_, err := os.Stat(filepath.Join(tempDir, path))
				assert.True(t, err != nil && os.IsNotExist(err), "Path should not exist: "+path)
			}
		})
	}
}
