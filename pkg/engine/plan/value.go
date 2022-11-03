package plan

import (
	"cuelang.org/go/cue"
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

func (val *Value) Lookup(s string) *Value {
	return WrapValue(val.cueValue.LookupPath(cue.ParsePath(s)))
}

func (val *Value) Pos() token.Pos {
	return val.cueValue.Pos()
}

func (val *Value) Decode(target any) error {
	return val.cueValue.Decode(target)
}
