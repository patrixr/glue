package runtime

import "fmt"

type Arguments struct {
	R    Runtime
	data []RTValue
}

func NewArguments(R Runtime, data []RTValue) *Arguments {
	return &Arguments{R, data}
}

func (args *Arguments) Len() int {
	return len(args.data)
}

func (args *Arguments) CheckString(i int) (RTString, error) {
	if i > len(args.data) {
		return nil, fmt.Errorf("Argument index out of bounds: %d", i)
	}

	return args.R.CheckString(args.data[i])
}

func (args *Arguments) CheckNumber(i int) (RTNumber, error) {
	if i > len(args.data) {
		return nil, fmt.Errorf("Argument index out of bounds: %d", i)
	}

	return args.R.CheckNumber(args.data[i])
}

func (args *Arguments) CheckBool(i int) (RTBool, error) {
	if i > len(args.data) {
		return nil, fmt.Errorf("Argument index out of bounds: %d", i)
	}

	return args.R.CheckBool(args.data[i])
}

func (args *Arguments) CheckDict(i int) (RTDict, error) {
	if i > len(args.data) {
		return nil, fmt.Errorf("Argument index out of bounds: %d", i)
	}

	return args.R.CheckDict(args.data[i])
}

func (args *Arguments) CheckFunction(i int) (RTFunction, error) {
	if i > len(args.data) {
		return nil, fmt.Errorf("Argument index out of bounds: %d", i)
	}

	return args.R.CheckFunction(args.data[i])
}

func (args *Arguments) EnsureString(i int) RTString {
	if i >= len(args.data) {
		args.R.RaiseError("Argument index out of bounds: %d", i)
	}

	return args.R.EnsureString(args.data[i])
}

func (args *Arguments) EnsureNumber(i int) RTNumber {
	if i > len(args.data) {
		args.R.RaiseError("Argument index out of bounds: %d", i)
	}

	return args.R.EnsureNumber(args.data[i])
}

func (args *Arguments) EnsureBool(i int) RTBool {
	if i > len(args.data) {
		args.R.RaiseError("Argument index out of bounds: %d", i)
	}

	return args.R.EnsureBool(args.data[i])
}

func (args *Arguments) EnsureDict(i int) RTDict {
	if i > len(args.data) {
		args.R.RaiseError("Argument index out of bounds: %d", i)
	}

	return args.R.EnsureDict(args.data[i])
}

func (args *Arguments) EnsureFunction(i int) RTFunction {
	if i > len(args.data) {
		args.R.RaiseError("Argument index out of bounds: %d", i)
	}

	return args.R.EnsureFunction(args.data[i])
}

func (args *Arguments) Get(i int) RTValue {
	if i > len(args.data) {
		args.R.RaiseError("Argument index out of bounds: %d", i)
	}

	return args.data[i]
}
