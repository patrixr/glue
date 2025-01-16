package runtime

type RTValue interface {
	String() string
	Type() Type
}

type RTString interface {
	RTValue
}

type RTDict interface {
	RTValue
	Map() map[interface{}]interface{}
}

type RTArray interface {
	RTValue
	Map() []interface{}
}

type RTNumber interface {
	RTValue
}

type RTFunction interface {
	RTValue
}

type RTBool interface {
	RTValue
	Value() bool
}

type RTNil interface {
	RTValue
}
