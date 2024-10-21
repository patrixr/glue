package luatools

import (
	"encoding/json"
	"github.com/patrixr/q"
	"os"
)

type LuaRC map[string]interface{}

func LoadLuaRC(file string) (*LuaRC, error) {
	luarc := make(LuaRC)

	data, err := os.ReadFile(file)

	if err != nil {
		return &luarc, err
	}

	err = json.Unmarshal(data, &luarc)
	if err != nil {
		return &luarc, err
	}

	return &luarc, nil
}

func (luarc *LuaRC) ToJSON() (string, error) {
	bytes, err := json.Marshal(luarc)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (luarc *LuaRC) AddLibrary(name string) {
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

	if libs, ok := (*luarc)["workspace.library"].([]interface{}); ok {
		(*luarc)["workspace.library"] = apply(libs)
		return
	}

	if workspace, ok := (*luarc)["workspace"].(map[string]interface{}); ok {
		if libs, ok := (workspace)["library"].([]interface{}); ok {
			(workspace)["library"] = apply(libs)
			return
		}

		(*luarc)["workspace"] = map[string]interface{}{
			"library": []string{name},
		}
		return
	}

	(*luarc)["workspace.library"] = []string{name}
}
