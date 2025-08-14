package reflection

import (
	"fmt"
	"reflect"
)

type ChachedPathStep struct {
	Offset    uintptr
	IsPointer bool
	FieldName string
	FieldType reflect.Type
	Childrens map[string]*ChachedPathStep
}

var g_PathsChache = make(map[string]*ChachedPathStep, 0)
var g_CachedBreanchesCount uint = 0

func FindOffsetForField(obj any, path []string) (uintptr, []uintptr, reflect.Type, error) {
	_type := reflect.TypeOf(obj)

	if _type.Kind() == reflect.Ptr {
		_type = _type.Elem()
	}

	if _type.Kind() != reflect.Struct {
		return 0, nil, nil, fmt.Errorf("Field offset searching error: type of input data is not Struct")
	}

	return findOffsetForField(_type, path)
}

func findOffsetForField(_type reflect.Type, path []string) (uintptr, []uintptr, reflect.Type, error) {
	var offsetsToPointers []uintptr
	existingOffset, existingStep, offsetsToPointers, foundComplete := checkFieldOffsetInCache(path)

	if offsetsToPointers == nil {
		offsetsToPointers = make([]uintptr, 0)
	}

	if foundComplete {
		return existingOffset, offsetsToPointers, existingStep.FieldType, nil
	}

	var currentType reflect.Type
	var offset uintptr = existingOffset
	var fieldType reflect.Type
	lastExistingStepWasNil := existingStep == nil

	if existingStep != nil {
		currentType = existingStep.FieldType
		path = cutPathPartThetExistInChache(existingStep, path)
	} else {
		currentType = _type
	}

	for i, step := range path {
		field, isExist := currentType.FieldByName(step)
		if !isExist {
			return 0, offsetsToPointers, nil, fmt.Errorf("Field offset searching error: Not found field %s in struct %s", step, currentType.Name())
		}

		fieldKind := field.Type.Kind()
		isLastStep := i == len(path)-1
		if isLastStep {
			fieldType = field.Type
		}

		if fieldKind == reflect.Ptr {
			offset += field.Offset
			offsetsToPointers = append(offsetsToPointers, offset)
			offset = 0
			currentType = field.Type.Elem()
		} else {
			offset += field.Offset
			if !isLastStep {
				currentType = field.Type
			}
		}

		if lastExistingStepWasNil {
			existingStep = cacheToGlobal(step, fieldKind == reflect.Ptr, field.Offset, field.Type)
			lastExistingStepWasNil = false
		} else {
			existingStep = cacheAsChildren(existingStep, step, fieldKind == reflect.Ptr, field.Offset, field.Type)
		}

		if !isLastStep && fieldKind != reflect.Ptr && fieldKind != reflect.Struct {
			return 0, nil, nil, fmt.Errorf("Field offset searching error: Field %s in struct %s is not struct or ptr on it. Library no suppor this types for routing", step, currentType.Name())
		}
	}

	return offset, offsetsToPointers, fieldType, nil
}

func checkFieldOffsetInCache(path []string) (uintptr, *ChachedPathStep, []uintptr, bool) {
	if g_CachedBreanchesCount == 0 {
		return cacheMiss()
	}

	pointersOffset := make([]uintptr, 0)
	var lastExistingStep *ChachedPathStep = nil
	var existingOffset uintptr = 0
	layer := g_PathsChache

	for _, step := range path {
		stepStruct, isExist := layer[step]
		if !isExist {
			break
		}

		layer = stepStruct.Childrens
		lastExistingStep = stepStruct
		if stepStruct.IsPointer {
			existingOffset += stepStruct.Offset
			pointersOffset = append(pointersOffset, existingOffset)
			existingOffset = 0
		} else {
			existingOffset += stepStruct.Offset
		}
	}

	foundComplete := lastExistingStep != nil &&
		len(path) > 0 &&
		lastExistingStep.FieldName == path[len(path)-1]
	return existingOffset, lastExistingStep, pointersOffset, foundComplete
}

func cacheMiss() (uintptr, *ChachedPathStep, []uintptr, bool) {
	return uintptr(0), nil, nil, false
}

func cutPathPartThetExistInChache(step *ChachedPathStep, path []string) []string {
	for i, part := range path {
		if step.FieldName == part {
			return path[i+1:]
		}
	}
	return path
}

func cacheAsChildren(perent *ChachedPathStep, name string, isPointer bool, offset uintptr, _type reflect.Type) *ChachedPathStep {
	perent.Childrens[name] = &ChachedPathStep{
		Offset:    offset,
		IsPointer: isPointer,
		FieldType: _type,
		FieldName: name,
		Childrens: make(map[string]*ChachedPathStep),
	}
	return perent.Childrens[name]
}

func cacheToGlobal(name string, isPointer bool, offset uintptr, _type reflect.Type) *ChachedPathStep {
	g_PathsChache[name] = &ChachedPathStep{
		Offset:    offset,
		IsPointer: isPointer,
		FieldType: _type,
		FieldName: name,
		Childrens: make(map[string]*ChachedPathStep),
	}
	g_CachedBreanchesCount++
	return g_PathsChache[name]
}

func ClearCache() {
	g_PathsChache = make(map[string]*ChachedPathStep)
	g_CachedBreanchesCount = 0
}

func CacheCount() uint {
	return g_CachedBreanchesCount
}
