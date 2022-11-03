package plan

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/token"
)

func WrapValue(v cue.Value) *Value {
	return &Value{v: v}
}

type Value struct {
	v cue.Value
}

func (val *Value) Path() cue.Path {
	return val.v.Path()
}

func (val *Value) Value() cue.Value {
	return val.v
}

func (val *Value) Lookup(s string) *Value {
	return WrapValue(val.v.LookupPath(cue.ParsePath(s)))
}

func (val *Value) Decode(target any) error {
	if err := val.v.Err(); err != nil {
		return err
	}
	return val.v.Decode(target)
}

func (val *Value) Pos() token.Pos {
	return val.v.Pos()
}
