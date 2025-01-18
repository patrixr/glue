package lua

import (
	"os"
	"path/filepath"

	"github.com/patrixr/glue/pkg/core"
	. "github.com/patrixr/glue/pkg/scaffold"
	"github.com/patrixr/q"
)

type LuaScaffold struct {
	glue *core.Glue
}

func NewLuaScaffold(glue *core.Glue) Scaffold {
	return &LuaScaffold{glue}
}

func (s *LuaScaffold) Setup(folder string) error {
	q.AssertNoError(os.MkdirAll(folder, 0770))

	s.glue.Log.Info("Initializing Glue folder", "path", folder)

	libFolderName := ".glue"
	libFolder := filepath.Join(folder, libFolderName)
	libFile := filepath.Join(libFolder, "typedefs.lua")

	// Write all the typedefs under .glue/typedefs.lua
	q.AssertNoError(os.MkdirAll(libFolder, 0770))
	q.AssertNoError(os.WriteFile(libFile, []byte(s.Typegen()), 0644))

	s.glue.Log.Info("Installing glue type definitions", "path", libFile)

	// Load the luarc configuration
	luarc, err := InitLuaRC(folder)

	q.AssertNoError(err)

	// Inject the glue typedefs into the luarc configuration
	luarc.AddLibrary("./" + libFolderName)
	luarc.SaveToFolder(folder)

	glueScript := filepath.Join(folder, "glue.lua")

	if _, err := core.TryFindGlueFile(folder); err != nil {
		q.AssertNoError(os.WriteFile(glueScript, []byte("-- Input code below"), 0644))
		s.glue.Log.Info("Glue script initialized", "script", glueScript)
	}

	s.glue.Log.Info("Glue scaffold setup complete", "workdir", folder)

	return nil
}

func (s *LuaScaffold) Typegen() string {
	return GenerateTypeDefinitions(s.glue)
}
