package utils

import (
	"fmt"
	"reflect"
	"unsafe"
)

type CachedPathStep struct {
	Offset    uintptr
	IsPointer bool
	FieldName string
	FieldType reflect.Type
	Children  map[string]*CachedPathStep
}

type PathCache struct {
	roots map[string]*CachedPathStep
}

func NewPathCache() *PathCache {
	return &PathCache{
		roots: make(map[string]*CachedPathStep),
	}
}

func (pc *PathCache) TakePtrOnField(obj any, path []string) (unsafe.Pointer, error) {
	val := reflect.ValueOf(obj)

	for val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	if val.Kind() == reflect.Struct {
		val = val.Addr()
	}

	if val.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("need pointer to struct, got %s", val.Kind())
	}

	start := unsafe.Pointer(val.Pointer())

	offset, ptrOffsets, err := pc.findOffsetForField(val.Interface(), path)
	if err != nil {
		return nil, err
	}

	return applyOffsets(start, offset, ptrOffsets), nil
}

func applyOffsets(start unsafe.Pointer, offset uintptr, ptrOffsets []uintptr) unsafe.Pointer {
	current := start
	for _, o := range ptrOffsets {
		ptrToField := unsafe.Add(current, o)
		current = *(*unsafe.Pointer)(ptrToField)
	}
	return unsafe.Add(current, offset)
}

func (pc *PathCache) findOffsetForField(obj any, path []string) (uintptr, []uintptr, error) {
	typ := reflect.TypeOf(obj)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return 0, nil, fmt.Errorf("expected struct, got %s", typ.Kind())
	}
	return pc.resolveOffsets(typ, path)
}

func (pc *PathCache) resolveOffsets(typ reflect.Type, path []string) (uintptr, []uintptr, error) {
	offsetsToPtrs := []uintptr{}
	existingOffset, step, cachedPtrs, found := pc.lookupCache(path)

	if found {
		return existingOffset, cachedPtrs, nil
	}
	if cachedPtrs != nil {
		offsetsToPtrs = cachedPtrs
	}

	var (
		currentType = typ
		offset      = existingOffset
	)

	if step != nil {
		currentType = step.FieldType
		path = trimPath(step, path)
	}

	for i, name := range path {
		field, ok := currentType.FieldByName(name)
		if !ok {
			return 0, offsetsToPtrs, fmt.Errorf("field %s not found in %s", name, currentType.Name())
		}

		isLast := i == len(path)-1
		isPtr := field.Type.Kind() == reflect.Ptr

		offset += field.Offset
		if isPtr {
			offsetsToPtrs = append(offsetsToPtrs, offset)
			offset = 0
			currentType = field.Type.Elem()
		} else if !isLast {
			if field.Type.Kind() != reflect.Struct {
				return 0, nil, fmt.Errorf("field %s in %s is not struct or pointer", name, currentType.Name())
			}
			currentType = field.Type
		}

		if step == nil {
			step = pc.cacheRoot(name, isPtr, field.Offset, field.Type)
		} else {
			step = pc.cacheChild(step, name, isPtr, field.Offset, field.Type)
		}
	}
	return offset, offsetsToPtrs, nil
}

func (pc *PathCache) lookupCache(path []string) (uintptr, *CachedPathStep, []uintptr, bool) {
	var (
		layer          = pc.roots
		step           *CachedPathStep
		offset         uintptr
		pointerOffsets []uintptr
	)

	for _, name := range path {
		next, ok := layer[name]
		if !ok {
			break
		}
		step = next
		layer = next.Children

		if next.IsPointer {
			offset += next.Offset
			pointerOffsets = append(pointerOffsets, offset)
			offset = 0
		} else {
			offset += next.Offset
		}
	}
	complete := step != nil && step.FieldName == path[len(path)-1]
	return offset, step, pointerOffsets, complete
}

func (pc *PathCache) cacheRoot(name string, isPointer bool, offset uintptr, typ reflect.Type) *CachedPathStep {
	step := &CachedPathStep{
		Offset:    offset,
		IsPointer: isPointer,
		FieldType: typ,
		FieldName: name,
		Children:  make(map[string]*CachedPathStep),
	}
	pc.roots[name] = step
	return step
}

func (pc *PathCache) cacheChild(parent *CachedPathStep, name string, isPointer bool, offset uintptr, typ reflect.Type) *CachedPathStep {
	step := &CachedPathStep{
		Offset:    offset,
		IsPointer: isPointer,
		FieldType: typ,
		FieldName: name,
		Children:  make(map[string]*CachedPathStep),
	}
	parent.Children[name] = step
	return step
}

func trimPath(step *CachedPathStep, path []string) []string {
	for i, p := range path {
		if p == step.FieldName {
			return path[i+1:]
		}
	}
	return path
}
