package lua

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/patrixr/q"
)

type LuaRC map[string]interface{}

// LoadLuaRC tries to loads a LuaRC configuration from a specified file.
// It reads the file content and unmarshals the JSON data into a LuaRC map.
func LoadLuaRC(file string) (LuaRC, error) {
	luarc := make(LuaRC)

	data, err := os.ReadFile(file)

	if err != nil {
		return luarc, err
	}

	err = json.Unmarshal(data, &luarc)
	if err != nil {
		return luarc, err
	}

	return luarc, nil
}

// InitLuaRC initializes a LuaRC configuration from a specified folder.
// It looks for a file named ".luarc.json" within the folder.
// If the file does not exist, it initializes an empty one
func InitLuaRC(folder string) (LuaRC, error) {
	luarcPath := filepath.Join(folder, ".luarc.json")
	_, err := os.Stat(luarcPath)

	if os.IsNotExist(err) {
		return LuaRC{}, nil
	}

	if err != nil {
		return nil, err
	}

	return LoadLuaRC(luarcPath)
}

// SaveToFolder saves the LuaRC configuration to a specified folder.
// It writes the configuration to a file named ".luarc.json" within the folder.
// The file permissions are set to 0644 by default, but if the file already exists,
// it retains the existing file permissions
func (luarc LuaRC) SaveToFolder(folder string) error {
	var mode fs.FileMode = 0644

	luarcPath := filepath.Join(folder, ".luarc.json")
	stat, err := os.Stat(luarcPath)

	if err != nil && stat != nil {
		mode = stat.Mode()
	}

	json, err := luarc.ToJSON()

	if err != nil {
		return err
	}

	if err = os.WriteFile(luarcPath, []byte(json), mode); err != nil {
		return err
	}

	return nil
}

// ToJSON converts the LuaRC configuration to its JSON string representation.
// It marshals the LuaRC map into a JSON formatted string.
// If the marshalling process encounters an error, it returns an empty string and the error.
// Otherwise, it returns the JSON string and a nil error.
func (luarc LuaRC) ToJSON() (string, error) {
	bytes, err := json.Marshal(luarc)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// AddLibrary adds a library name to the LuaRC configuration.
// It ensures that the library name is added to the "workspace.library" key.
// If the "workspace.library" key does not exist, it creates it.
// If the library name already exists in the list, it does not add it again.
func (luarc LuaRC) AddLibrary(name string) {
	equals := func(obj interface{}, _ int) bool {
		str, _ := obj.(string)
		return str == name
	}

	apply := func(libs []interface{}) []interface{} {
		if found, _, _ := q.Find(libs, equals); found {
			return libs
		}
		return append(libs, name)
	}

	if libs, ok := luarc["workspace.library"].([]interface{}); ok {
		luarc["workspace.library"] = apply(libs)
		return
	}

	if workspace, ok := luarc["workspace"].(map[string]interface{}); ok {
		if libs, ok := (workspace)["library"].([]interface{}); ok {
			(workspace)["library"] = apply(libs)
			return
		}

		luarc["workspace"] = map[string]interface{}{
			"library": []string{name},
		}
		return
	}

	luarc["workspace.library"] = []string{name}
}
