package modules_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/patrixr/glue/pkg/modules"
	"github.com/patrixr/q"
	"github.com/stretchr/testify/assert"
)

func TestBlockInString(t *testing.T) {

	type testCase struct {
		name   string
		input  string
		output string
		props  modules.BlockOpts
	}

	testCases := []testCase{
		testCase{
			name: "Block should be removed from the text",
			input: q.Paragraph(`
			    line 1
			    line 2
			    # BEGIN MANAGED BLOCK
			    Test Block
			    # END MANAGED BLOCK
			    line 3
		    `),
			output: q.Paragraph(`
		    	line 1
			    line 2
			    line 3
		    `),
			props: modules.BlockOpts{
				Block: "Test block",
				State: false,
			},
		},
		testCase{
			name: "Block with custom markers should be removed from the text",
			input: q.Paragraph(`
			    line 1
			    line 2
			    # AAA
			    Test Block
			    # BBB
			    line 3
		    `),
			output: q.Paragraph(`
		    	line 1
			    line 2
			    line 3
		    `),
			props: modules.BlockOpts{
				Block:       "Test block",
				State:       false,
				Markerbegin: "AAA",
				Markerend:   "BBB",
			},
		},
		testCase{
			name: "Existing block should be replaced",
			input: q.Paragraph(`
				line 1
				line 2
				# BEGIN MANAGED BLOCK
				Something else
				# END MANAGED BLOCK
				line 3
		    `),
			output: q.Paragraph(`
				line 1
				line 2
				# BEGIN MANAGED BLOCK
				Test Block
				# END MANAGED BLOCK
				line 3
		    `),
			props: modules.BlockOpts{
				Block: "Test Block",
				State: true,
			},
		},
		testCase{
			name: "Existing block should be replaced even if empty",
			input: q.Paragraph(`
				line 1
				line 2
				# BEGIN MANAGED BLOCK
				# END MANAGED BLOCK
				line 3
		    `),
			output: q.Paragraph(`
				line 1
				line 2
				# BEGIN MANAGED BLOCK
				Test Block
				# END MANAGED BLOCK
				line 3
		    `),
			props: modules.BlockOpts{
				Block: "Test Block",
				State: true,
			},
		},
		testCase{
			name: "Block should be added at the end by default",
			input: q.Paragraph(`
				line 1
				line 2
				line 3`),
			output: q.Paragraph(`
				line 1
				line 2
				line 3
				# BEGIN MANAGED BLOCK
				Test Block
				# END MANAGED BLOCK
		    `),
			props: modules.BlockOpts{
				Block: "Test Block",
				State: true,
			},
		},
		testCase{
			name: "Block should be after the Insertafter",
			input: q.Paragraph(`
				line 1
				line 2
				line 3
		    `),
			output: q.Paragraph(`
				line 1
				line 2
				# BEGIN MANAGED BLOCK
				Test Block
				# END MANAGED BLOCK
				line 3
		    `),
			props: modules.BlockOpts{
				Block:       "Test Block",
				State:       true,
				Insertafter: "line 2",
			},
		},
		testCase{
			name: "Block should be after the Insertafter on the last line",
			input: q.Paragraph(`
				line 1
				line 2
				line 3
		    `),
			output: q.Paragraph(`
				line 1
				line 2
				line 3
				# BEGIN MANAGED BLOCK
				Test Block
				# END MANAGED BLOCK
		    `),
			props: modules.BlockOpts{
				Block:       "Test Block",
				State:       true,
				Insertafter: "line 3",
			},
		},
		testCase{
			name: "block should be added at the Insertafter marker using a Regex",
			input: q.Paragraph(`
				line 1
				line 2
				line 3
		    `),
			output: q.Paragraph(`
				line 1
				line 2
				# BEGIN MANAGED BLOCK
				Test Block
				# END MANAGED BLOCK
				line 3
		    `),
			props: modules.BlockOpts{
				Block:       "Test Block",
				State:       true,
				Insertafter: ".*2",
			},
		},
		testCase{
			name: "block should be added at the Insertbefore marker if specified",
			input: q.Paragraph(`
				line 1
				line 2
				line 3
		    `),
			output: q.Paragraph(`
				line 1
				# BEGIN MANAGED BLOCK
				Test Block
				# END MANAGED BLOCK
				line 2
				line 3
		    `),
			props: modules.BlockOpts{
				Block:        "Test Block",
				State:        true,
				Insertbefore: "line 2",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			out := modules.BlockInString(testCase.input, testCase.props)
			assert.Equal(t, testCase.output, out)
		})

		t.Run(testCase.name+" in file", func(t *testing.T) {
			tmp := t.TempDir()
			tmpfile := filepath.Join(tmp, "sample")

			file, err := os.Create(tmpfile)

			assert.Nil(t, err)

			_, err = file.Write([]byte(testCase.input))

			file.Close()

			assert.Nil(t, err)

			testCase.props.Path = tmpfile

			err = modules.BlockInFile(testCase.props)

			assert.Nil(t, err)

			out, err := os.ReadFile(tmpfile)

			assert.Nil(t, err)

			assert.Equal(t, testCase.output, string(out))
		})
	}

	t.Run("a backup is created with the older version when backup field is set to true", func(t *testing.T) {
		tmp := t.TempDir()
		tmpfile := filepath.Join(tmp, "sample")

		file, err := os.Create(tmpfile)

		assert.Nil(t, err)

		_, err = file.Write([]byte(q.Paragraph(`
			line 1
			line 2
			line 3
		`)))

		file.Close()

		assert.Nil(t, err)

		err = modules.BlockInFile(modules.BlockOpts{
			Block:        "Test Block",
			State:        true,
			Insertbefore: "line 3",
			Path:         tmpfile,
			Backup:       true,
		})

		assert.Nil(t, err)

		files, err := os.ReadDir(tmp)

		assert.Nil(t, err)

		assert.Equal(t, 2, len(files))

		var backupFile string
		for _, file := range files {
			if strings.Contains(file.Name(), "backup") {
				backupFile = filepath.Join(tmp, file.Name())
				break
			}
		}

		assert.NotEmpty(t, backupFile, "Backup file was not created")
	})

	t.Run("a file is created if it is missing and the create flag is set to true", func(t *testing.T) {
		tmp := t.TempDir()
		tmpfile := filepath.Join(tmp, "sample")

		err := modules.BlockInFile(modules.BlockOpts{
			Block:        "Test Block",
			State:        true,
			Insertbefore: "line 3",
			Path:         tmpfile,
			Create:       true,
		})

		assert.NoError(t, err)
		assert.FileExists(t, tmpfile)
		content, err := os.ReadFile(tmpfile)
		assert.NoError(t, err)
		assert.Equal(t,
			q.Paragraph(`
				# BEGIN MANAGED BLOCK
				Test Block
				# END MANAGED BLOCK
			`),
			string(content))
	})

	t.Run("an error is returned if the file doesn't exist and create flag is set to false", func(t *testing.T) {
		tmp := t.TempDir()
		tmpfile := filepath.Join(tmp, "sample")

		err := modules.BlockInFile(modules.BlockOpts{
			Block:        "Test Block",
			State:        true,
			Insertbefore: "line 3",
			Path:         tmpfile,
			Create:       false,
		})

		assert.Error(t, err)
	})
}
