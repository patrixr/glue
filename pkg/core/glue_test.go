package core_test

import (
	"testing"

	"github.com/patrixr/glue/pkg/core"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
)

func Test_FunctionCalling(t *testing.T) {
	called := false
	glue := core.NewGlue()

	defer glue.Close()

	glue.Plug().
		Name("foo").
		Do(func(L *lua.LState) (int, error) {
			called = true
			return 0, nil
		})

	t.Run("should fail to call a non-existing function", func(t *testing.T) {
		glue := core.NewGlue()
		err := glue.ExecuteString("bar()")
		assert.NotNil(t, err)
		assert.True(t, len(glue.ExecutionTrace) == 0)
	})

	t.Run("should allow calling an existing function", func(t *testing.T) {
		err := glue.ExecuteString("foo()")
		assert.Nil(t, err)
		assert.True(t, len(glue.ExecutionTrace) == 1)
		assert.Equal(t, glue.ExecutionTrace[0].Name, "foo")
		assert.True(t, called)
	})

	t.Run("Native lua functions", func(t *testing.T) {
		t.Run("should be forbidden by default", func(t *testing.T) {
			err := glue.ExecuteString("os.getenv(\"ENV\")")
			assert.NotNil(t, err)
		})
	})
}
