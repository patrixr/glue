package runtime

type Type interface {
	Is(other Type) bool
	Id() int
	Name() string
}

const (
	StringId = iota
	NumberId
	BoolId
	DictId
	FuncId
	ArrayId
	NilId
	AnyId
)

var (
	STRING = PrimitiveType{StringId, "string"}
	NUMBER = PrimitiveType{NumberId, "number"}
	BOOL   = PrimitiveType{BoolId, "bool"}
	DICT   = PrimitiveType{DictId, "dict"}
	FUNC   = PrimitiveType{FuncId, "func"}
	ARRAY  = PrimitiveType{ArrayId, "array"}
	NIL    = PrimitiveType{NilId, "nil"}
	ANY    = PrimitiveType{AnyId, "any"}
)

type PrimitiveType struct {
	id   int
	name string
}

func (typ PrimitiveType) Is(other Type) bool {
	return typ.id == other.Id()
}

func (typ PrimitiveType) Id() int {
	return typ.id
}

func (typ PrimitiveType) Name() string {
	return typ.name
}

type TypedArrayType struct {
	Type
	ItemType Type
}

func TypedArray(itemType Type) TypedArrayType {
	return TypedArrayType{ARRAY, itemType}
}

type CustomStructType struct {
	Type
	Fields []Field
}

type Field struct {
	Name     string
	Type     Type
	Desc     string
	Optional bool
}

func CustomStruct(name string, fields []Field) CustomStructType {
	return CustomStructType{DICT, fields}
}

func NewField(name string, typ Type, desc string) Field {
	optional := false

	if name[len(name)-1] == '?' {
		optional = true
		name = name[:len(name)-1]
	}

	AssertValidSymbolName(name)

	return Field{name, typ, desc, optional}
}

type CustomFunctionType struct {
	Type
	Args []ArgDef
}

func Custom(name string, fields []Field) CustomStructType {
	return CustomStructType{DICT, fields}
}

func TypeName(t Type) string {
	return t.Name()
}

type DescribableType struct {
	Type
	desc string
}

func (t DescribableType) Describe() string {
	return t.desc
}

func TypeWithDesc(t Type, desc string) DescribableType {
	return DescribableType{t, desc}
}
