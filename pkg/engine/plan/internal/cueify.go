package internal

import (
	"bytes"
	"encoding"
	"fmt"
	"go/ast"
	"reflect"
	"strings"
)

type OneOfType interface {
	OneOf() []any
}

var oneOfType = reflect.TypeOf((*OneOfType)(nil)).Elem()
var textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()

func newConvert(r TaskRegister) *convert {
	return &convert{
		r:    r,
		defs: map[reflect.Type]bool{},
	}
}

type convert struct {
	r    TaskRegister
	defs map[reflect.Type]bool
}

type opt struct {
	naming string
	extra  string
}

func (c *convert) toCueType(tpe reflect.Type, o opt) []byte {
	if o.naming == "" && tpe.PkgPath() != "" {
		if _, ok := c.defs[tpe]; !ok {
			c.defs[tpe] = true
			c.r.Register(reflect.New(tpe).Interface())
		}

		if o.extra != "" {
			return []byte(fmt.Sprintf(`#%s & { 
  %s 
}`, tpe.Name(), o.extra))
		}
		return []byte(fmt.Sprintf("#%s", tpe.Name()))
	}

	if tpe.Implements(textMarshalerType) {
		return []byte("string")
	}

	if tpe.Implements(oneOfType) {
		if ot, ok := reflect.New(tpe).Interface().(OneOfType); ok {
			types := ot.OneOf()
			b := bytes.NewBuffer(nil)

			for i := range types {
				t := reflect.TypeOf(types[i])
				if t.Kind() == reflect.Ptr {
					t = t.Elem()
				}
				if i > 0 {
					b.WriteString(" | ")
				}
				b.Write(c.toCueType(t, opt{extra: o.extra}))
			}

			return b.Bytes()
		}
	}

	switch tpe.Kind() {
	case reflect.Ptr:
		return []byte(fmt.Sprintf("%s | null", c.toCueType(tpe.Elem(), opt{extra: o.extra})))
	case reflect.Map:
		return []byte(fmt.Sprintf("[X=%s]: %s", c.toCueType(tpe.Key(), opt{extra: o.extra}), c.toCueType(tpe.Elem(), opt{extra: o.extra})))
	case reflect.Slice:
		return []byte(fmt.Sprintf("[...%s]", c.toCueType(tpe.Elem(), opt{extra: o.extra})))
	case reflect.Struct:
		b := bytes.NewBuffer(nil)

		_, _ = fmt.Fprintf(b, `{
`)

		walkFields(tpe, func(i *fieldInfo) {
			t := i.tpe

			// FIXME may support other inline
			if i.inline {
				if t.Kind() == reflect.Map {
					_, _ = fmt.Fprintf(b, `[!~"\\$wagon"]: %s`, c.toCueType(t.Elem(), opt{
						extra: i.cueExtra,
					}))
				}
				return
			}

			if i.optional {
				if t.Kind() == reflect.Ptr {
					t = t.Elem()
				}
				_, _ = fmt.Fprintf(b, "%s?: ", i.name)
			} else {
				_, _ = fmt.Fprintf(b, "%s: ", i.name)
			}

			cueType := c.toCueType(t, opt{
				extra: i.cueExtra,
			})

			if len(i.enum) > 0 {
				for i, e := range i.enum {
					if i > 0 {
						_, _ = fmt.Fprint(b, " | ")
					}
					_, _ = fmt.Fprintf(b, `%q`, e)
				}
			} else {
				_, _ = fmt.Fprintf(b, "%s", cueType)
			}

			if i.defaultValue != nil {
				switch string(cueType) {
				case "[]byte":
					_, _ = fmt.Fprintf(b, ` | *'%s'`, *i.defaultValue)
				case "string":
					_, _ = fmt.Fprintf(b, ` | *%q`, *i.defaultValue)
				default:
					_, _ = fmt.Fprintf(b, ` | *%v`, *i.defaultValue)
				}
			}

			if len(i.attrs) > 0 {
				_, _ = fmt.Fprintf(b, " @wagon(%s)", strings.Join(i.attrs, ","))
			}

			_, _ = fmt.Fprint(b, "\n")
		})

		_, _ = fmt.Fprintf(b, `}`)

		return b.Bytes()
	case reflect.Interface:
		return []byte("_")
	default:
		return []byte(tpe.Kind().String())
	}
}

type fieldInfo struct {
	name         string
	cueExtra     string
	idx          int
	tpe          reflect.Type
	optional     bool
	inline       bool
	defaultValue *string
	enum         []string
	attrs        []string
}

func (i *fieldInfo) EmptyDefaults() (string, bool) {
	if i.tpe.PkgPath() != "" {
		return "", false
	}

	switch i.tpe.Kind() {
	case reflect.Slice:
		return "", false
	case reflect.Map:
		return "", false
	case reflect.Interface:
		return "", false
	}
	return fmt.Sprintf("%v", reflect.New(i.tpe).Elem()), true
}

func (i *fieldInfo) HasAttr(expectAttr string) bool {
	for _, attr := range i.attrs {
		if attr == expectAttr {
			return true
		}
	}
	return false
}

func walkFields(s reflect.Type, each func(info *fieldInfo)) {
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		if !ast.IsExported(f.Name) {
			continue
		}

		info := &fieldInfo{}
		info.idx = i
		info.name = f.Name
		info.tpe = f.Type

		jsonTag, hasJsonTag := f.Tag.Lookup("json")
		if !hasJsonTag {
			if f.Anonymous && f.Type.Kind() == reflect.Struct {
				walkFields(f.Type, each)
			}
			continue
		}

		if strings.Contains(jsonTag, ",omitempty") {
			info.optional = true
		}

		if strings.Contains(jsonTag, ",inline") {
			info.inline = true
			info.name = ""
		}

		if cueExtra, hasCueExtra := f.Tag.Lookup("cueExtra"); hasCueExtra {
			info.cueExtra = cueExtra
		}

		wagonTag, hasWagonTag := f.Tag.Lookup("wagon")
		if jsonTag == "-" && !hasWagonTag {
			continue
		}

		if jsonName := strings.SplitN(jsonTag, ",", 2)[0]; jsonName != "" {
			info.name = jsonName
		}

		if hasWagonTag {
			attrs := strings.Split(wagonTag, ",")

			for _, n := range attrs {
				parts := strings.SplitN(n, "=", 2)
				if len(parts) == 2 && parts[0] == "name" {
					info.name = parts[1]
					continue
				}
				info.attrs = append(info.attrs, n)
			}
		}

		if defaultValue, ok := f.Tag.Lookup("default"); ok {
			info.defaultValue = &defaultValue
		}

		if enumValue, ok := f.Tag.Lookup("enum"); ok {
			info.enum = strings.Split(enumValue, ",")
		}

		each(info)
	}
}
