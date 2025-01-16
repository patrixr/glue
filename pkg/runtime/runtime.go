package runtime

type CustomFunc func(runtime Runtime, args *Arguments) RTValue

type Runtime interface {
	Lang() string
	ExecFile(path string) error
	ExecString(source string) error
	String(str string) RTString
	Close()
	RaiseError(format string, args ...interface{})
	CheckString(v RTValue) (RTString, error)
	CheckNumber(v RTValue) (RTNumber, error)
	CheckDict(v RTValue) (RTDict, error)
	CheckBool(v RTValue) (RTBool, error)
	CheckFunction(v RTValue) (RTFunction, error)
	CheckArray(v RTValue) (RTArray, error)
	EnsureString(v RTValue) RTString
	EnsureNumber(v RTValue) RTNumber
	EnsureDict(v RTValue) RTDict
	EnsureBool(v RTValue) RTBool
	EnsureFunction(v RTValue) RTFunction
	EnsureArray(v RTValue) RTArray
	InvokeFunction(fn RTFunction, params ...RTValue) error
	InvokeFunctionSafe(fn RTFunction, params ...RTValue) error
	SetGlobal(name string, val RTValue) ([]string, error)
	SetFunction(
		name string,
		desc string,
		args []ArgDef,
		fn CustomFunc,
	) error
}

type ArgDef struct {
	Type Type
	Name string
	Desc string
}
