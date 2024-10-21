package luatools

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTempFile(t *testing.T, content string) string {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "luarc_test_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tmpfile.Name()
}

func TestLoadLuaRC(t *testing.T) {
	validContent := `{"key1": "value1", "key2": "value2"}`
	tmpfile := createTempFile(t, validContent)
	defer os.Remove(tmpfile)

	luarc, err := LoadLuaRC(tmpfile)
	assert.NoError(t, err)
	expected := LuaRC{"key1": "value1", "key2": "value2"}
	assert.Equal(t, expected, *luarc)

	invalidContent := `{"key1": "value1", "key2": }`
	tmpfile = createTempFile(t, invalidContent)
	defer os.Remove(tmpfile)

	_, err = LoadLuaRC(tmpfile)
	assert.Error(t, err)
}

func TestToJSON(t *testing.T) {
	luarc := LuaRC{"key1": "value1", "key2": "value2"}
	expected := `{"key1":"value1","key2":"value2"}`

	jsonStr, err := luarc.ToJSON()
	assert.NoError(t, err)
	assert.JSONEq(t, expected, jsonStr)
}

func TestAddLibrary(t *testing.T) {
	testCases := []struct {
		name         string
		initialJSON  string
		libraryToAdd string
		expectedJSON string
	}{
		{
			name:         "Add to empty WorkspaceLibrary",
			initialJSON:  `{}`,
			libraryToAdd: "new_lib",
			expectedJSON: `{"workspace.library":["new_lib"]}`,
		},
		{
			name:         "Add to existing WorkspaceLibrary",
			initialJSON:  `{"workspace.library":["existing_lib"]}`,
			libraryToAdd: "new_lib",
			expectedJSON: `{"workspace.library":["existing_lib","new_lib"]}`,
		},
		{
			name:         "Add to empty WorkspaceObject.Library",
			initialJSON:  `{"workspace":{}}`,
			libraryToAdd: "new_lib",
			expectedJSON: `{"workspace":{"library":["new_lib"]}}`,
		},
		{
			name:         "Add to existing WorkspaceObject.Library",
			initialJSON:  `{"workspace":{"library":["existing_lib"], "extra": "value"}}`, //with extra field
			libraryToAdd: "new_lib",
			expectedJSON: `{"workspace":{"library":["existing_lib","new_lib"], "extra": "value"}}`,
		},
		{
			name:         "Don't add duplicate to WorkspaceLibrary",
			initialJSON:  `{"workspace.library":["existing_lib"]}`,
			libraryToAdd: "existing_lib",
			expectedJSON: `{"workspace.library":["existing_lib"]}`,
		},
		{
			name:         "Don't add duplicate to WorkspaceObject.Library",
			initialJSON:  `{"workspace":{"library":["existing_lib"]}}`,
			libraryToAdd: "existing_lib",
			expectedJSON: `{"workspace":{"library":["existing_lib"]}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			luarc := LuaRC{}
			err := json.Unmarshal([]byte(tc.initialJSON), &luarc)
			assert.NoError(t, err)

			luarc.AddLibrary(tc.libraryToAdd)

			actualJSON, err := luarc.ToJSON()
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expectedJSON, actualJSON)
		})
	}
}
