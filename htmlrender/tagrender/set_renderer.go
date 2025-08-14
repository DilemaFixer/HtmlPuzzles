package tagrender

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	render "github.com/DilemaFixer/HtmlPuzzles/htmlrender"
	reflection "github.com/DilemaFixer/HtmlPuzzles/reflection"
)

type SetRenderer struct {
	way string
}

func NewSetRenderer(way string) *SetRenderer {
	return &SetRenderer{
		way: way,
	}
}

func (r *SetRenderer) Render(ctx render.Context) (*render.RenderedNode, error) {
	value, isExist := ctx.GetStringLayered(render.WayKey)
	if !isExist {
		return nil, fmt.Errorf("Error handling seter tag: %s key exist in current context layer", render.WayKey)
	}

	path := strings.Split(value, ".")

	offset, ptrsOffset, fieldType, err := reflection.FindOffsetForField(ctx.Obj, path)
	if err != nil {
		return nil, err
	}

	objValue := reflect.ValueOf(ctx.Obj)
	var objPtr unsafe.Pointer

	if objValue.Kind() == reflect.Ptr {
		objPtr = unsafe.Pointer(objValue.Pointer())
	} else {
		objPtr = unsafe.Pointer(objValue.UnsafeAddr())
	}

	fieldPtr, err := handlingPointersOnWay(objPtr, ptrsOffset)
	if err != nil {
		return nil, err
	}

	fieldPtr = unsafe.Add(fieldPtr, offset)

	str, err := convertPointerToValue(fieldPtr, fieldType)
	if err != nil {
		return nil, err
	}

	return &render.RenderedNode{
		Html: str,
	}, nil
}

func handlingPointersOnWay(obj unsafe.Pointer, ptrsOffset []uintptr) (unsafe.Pointer, error) {
	if obj == nil {
		return nil, fmt.Errorf("obj pointer is nil")
	}

	currentPtr := obj
	for i, offset := range ptrsOffset {
		if currentPtr == nil {
			return nil, fmt.Errorf("encountered nil pointer at offset index %d", i)
		}

		fieldPtr := unsafe.Add(currentPtr, offset)
		nextPtr := *(*unsafe.Pointer)(fieldPtr)

		if nextPtr == nil {
			return nil, fmt.Errorf("nil pointer found at offset %d (index %d)", offset, i)
		}

		currentPtr = nextPtr
	}

	return currentPtr, nil
}

func convertPointerToValue(ptr unsafe.Pointer, fieldType reflect.Type) (string, error) {
	if ptr == nil {
		return "", nil
	}
	if fieldType.Kind() == reflect.Ptr {
		ptrToPtr := (*unsafe.Pointer)(ptr)
		if *ptrToPtr == nil {
			return "", nil
		}
		return convertPointerToValue(*ptrToPtr, fieldType.Elem())
	}
	switch fieldType.Kind() {
	case reflect.String:
		strPtr := (*string)(ptr)
		return *strPtr, nil
	case reflect.Int:
		intPtr := (*int)(ptr)
		return fmt.Sprintf("%d", *intPtr), nil
	case reflect.Int8:
		int8Ptr := (*int8)(ptr)
		return fmt.Sprintf("%d", *int8Ptr), nil
	case reflect.Int16:
		int16Ptr := (*int16)(ptr)
		return fmt.Sprintf("%d", *int16Ptr), nil
	case reflect.Int32:
		int32Ptr := (*int32)(ptr)
		return fmt.Sprintf("%d", *int32Ptr), nil
	case reflect.Int64:
		int64Ptr := (*int64)(ptr)
		return fmt.Sprintf("%d", *int64Ptr), nil
	case reflect.Uint:
		uintPtr := (*uint)(ptr)
		return fmt.Sprintf("%d", *uintPtr), nil
	case reflect.Uint8:
		uint8Ptr := (*uint8)(ptr)
		return fmt.Sprintf("%d", *uint8Ptr), nil
	case reflect.Uint16:
		uint16Ptr := (*uint16)(ptr)
		return fmt.Sprintf("%d", *uint16Ptr), nil
	case reflect.Uint32:
		uint32Ptr := (*uint32)(ptr)
		return fmt.Sprintf("%d", *uint32Ptr), nil
	case reflect.Uint64:
		uint64Ptr := (*uint64)(ptr)
		return fmt.Sprintf("%d", *uint64Ptr), nil
	case reflect.Float32:
		float32Ptr := (*float32)(ptr)
		return fmt.Sprintf("%g", *float32Ptr), nil
	case reflect.Float64:
		float64Ptr := (*float64)(ptr)
		return fmt.Sprintf("%g", *float64Ptr), nil
	case reflect.Bool:
		boolPtr := (*bool)(ptr)
		return fmt.Sprintf("%t", *boolPtr), nil
	default:
		return "", fmt.Errorf("unsupported type: %v", fieldType.Kind())
	}
}

func (r *SetRenderer) AddChildren(*render.Renderer) {}
