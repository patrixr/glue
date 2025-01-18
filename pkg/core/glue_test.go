package core

import (
	"testing"

	"github.com/patrixr/glue/pkg/runtime"
	"github.com/stretchr/testify/assert"
)

func Test_FunctionCalling(t *testing.T) {
	called := false
	glue := NewGlue()

	defer glue.Close()

	glue.Plug("foo", FUNCTION).
		Do(func(R runtime.Runtime, args *runtime.Arguments) (runtime.RTValue, error) {
			called = true
			return nil, nil
		})

	t.Run("should fail to call a non-existing function", func(t *testing.T) {
		glue := NewGlue()
		err := glue.execString("bar()")
		assert.NotNil(t, err)
	})

	t.Run("should allow calling an existing function", func(t *testing.T) {
		err := glue.execString("foo()")
		assert.Nil(t, err)
		assert.True(t, called)
	})

	t.Run("Native lua functions", func(t *testing.T) {
		t.Run("should be forbidden by default", func(t *testing.T) {
			err := glue.execString("os.getenv(\"ENV\")")
			assert.NotNil(t, err)
		})
	})
}
