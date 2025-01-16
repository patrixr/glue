package runtime

import "strconv"

type Type int

const (
	STRING Type = iota
	NUMBER
	BOOL
	DICT
	FUNC
	ARRAY
	NIL
	ANY
)

func (typ Type) Is(other Type) bool {
	return typ == other
}

func TypeName(id Type) string {
	switch id {
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case BOOL:
		return "BOOL"
	case DICT:
		return "DICT"
	case ARRAY:
		return "ARRAY"
	case FUNC:
		return "FUNC"
	case ANY:
		return "ANY"
	case NIL:
		return "NIL"
	default:
		panic("Invalid type id " + strconv.Itoa(int(id)))
	}
}
