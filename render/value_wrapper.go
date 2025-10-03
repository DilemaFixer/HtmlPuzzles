package render

import (
	"fmt"
	"unsafe"
)

type RenderValueAs int

const (
	Non RenderValueAs = iota
	String
	Uint
	Float
	Integer
	Boolean
)

const (
	TString  = "string"
	TInt     = "int"
	TUint    = "uint"
	TFloat   = "float"
	TBoolean = "bool"
)

type Value interface {
	GetValue(ctx *Context) (interface{}, RenderValueAs, error)
}

type ValueFromContext struct {
	path []string
	typ  string
}

type StaticValue struct {
	value string
}

func NewValueFromContext(path []string, typ string) Value {
	return &ValueFromContext{path: path, typ: typ}
}

func NewStaticValue(value string) Value {
	return &StaticValue{value: value}
}

func (sv *StaticValue) GetValue(ctx *Context) (interface{}, RenderValueAs, error) {
	return sv.value, String, nil
}

func (v *ValueFromContext) GetValue(ctx *Context) (interface{}, RenderValueAs, error) {
	if len(v.path) == 0 {
		return nil, Non, fmt.Errorf("no value for type %s", v.typ)
	}

	if len(v.path) == 1 {
		return v.getPrimitive(ctx, v.path[0])
	}

	obj, exists := ctx.GetObject(v.path[0])
	if !exists {
		return nil, Non, fmt.Errorf("object %q not exist", v.path[0])
	}

	ptr, err := ctx.Offset.TakePtrOnField(obj, v.path[1:])
	if err != nil {
		return nil, Non, err
	}

	return v.getFromObject(ptr)
}

func (v *ValueFromContext) getPrimitive(ctx *Context, key string) (interface{}, RenderValueAs, error) {
	switch v.typ {
	case TInt:
		if value, ok := ctx.GetInt(key); ok {
			return value, Integer, nil
		}
		return nil, Non, fmt.Errorf("value %q (int64) not exist", key)

	case TUint:
		if value, ok := ctx.GetUint(key); ok {
			return value, Uint, nil
		}
		return nil, Non, fmt.Errorf("value %q (uint64) not exist", key)

	case TFloat:
		if value, ok := ctx.GetFloat(key); ok {
			return value, Float, nil
		}
		return nil, Non, fmt.Errorf("value %q (float64) not exist", key)

	case TString:
		if value, ok := ctx.GetString(key); ok {
			return value, String, nil
		}
		return nil, Non, fmt.Errorf("value %q (string) not exist", key)

	case TBoolean:
		if value, ok := ctx.GetBool(key); ok {
			return value, Boolean, nil
		}
		return nil, Non, fmt.Errorf("value %q (bool) not exist", key)

	default:
		return nil, Non, fmt.Errorf("unsupported type: %s", v.typ)
	}
}

func (v *ValueFromContext) getFromObject(ptr unsafe.Pointer) (interface{}, RenderValueAs, error) {
	switch v.typ {
	case TInt:
		return *(*int64)(ptr), Integer, nil
	case TUint:
		return *(*uint64)(ptr), Uint, nil
	case TFloat:
		return *(*float64)(ptr), Float, nil
	case TString:
		return *(*string)(ptr), String, nil
	case TBoolean:
		return *(*bool)(ptr), Boolean, nil
	default:
		return nil, Non, fmt.Errorf("unsupported type: %s", v.typ)
	}
}
