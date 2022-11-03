package plan

import (
	"cuelang.org/go/cue"
	cueformat "cuelang.org/go/cue/format"
	"cuelang.org/go/cue/token"
)

func WrapValue(cueValue cue.Value) *Value {
	return &Value{cueValue: cueValue}
}

type Value struct {
	cueValue cue.Value
}

func (val *Value) Path() cue.Path {
	return val.cueValue.Path()
}

func (val *Value) Value() cue.Value {
	return val.cueValue
}

func (val *Value) Pos() token.Pos {
	return val.cueValue.Pos()
}

func (val *Value) Decode(target any) error {
	return val.cueValue.Decode(target)
}

func (val *Value) Source() string {
	syn := val.cueValue.Syntax(
		cue.Final(),         // close structs and lists
		cue.Concrete(false), // allow incomplete values
		cue.DisallowCycles(true),
		cue.All(),
		cue.Docs(true),
	)
	data, _ := cueformat.Node(syn, cueformat.Simplify())
	return string(data)
}
