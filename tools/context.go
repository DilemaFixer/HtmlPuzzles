package tools

import "fmt"

type ValueType int

const (
	TypeString ValueType = iota
	TypeInt
	TypeUint
	TypeFloat
	TypeBool
	TypeObject
)

type Context struct {
	currentLayer int
	strings      map[string]string
	ints         map[string]int64
	uints        map[string]uint64
	floats       map[string]float64
	bools        map[string]bool
	objects      map[string]any
}

func NewContext() *Context {
	return &Context{
		currentLayer: 0,
		strings:      make(map[string]string),
		ints:         make(map[string]int64),
		uints:        make(map[string]uint64),
		floats:       make(map[string]float64),
		bools:        make(map[string]bool),
		objects:      make(map[string]any),
	}
}

func (c *Context) LayerUp() {
	c.currentLayer++
}

func (c *Context) LayerDown() {
	c.currentLayer--
}

func (c *Context) SetCurrentLayer(layer int) {
	c.currentLayer = layer
}

func (c *Context) GetCurrentLayer() int {
	return c.currentLayer
}

func (c *Context) makeLayerKey(key string) string {
	return fmt.Sprintf("%d:%s", c.currentLayer, key)
}

func (c *Context) Type(key string) ValueType {
	if _, ok := c.strings[key]; ok {
		return TypeString
	}
	if _, ok := c.ints[key]; ok {
		return TypeInt
	}
	if _, ok := c.uints[key]; ok {
		return TypeUint
	}
	if _, ok := c.floats[key]; ok {
		return TypeFloat
	}
	if _, ok := c.bools[key]; ok {
		return TypeBool
	}
	if _, ok := c.objects[key]; ok {
		return TypeObject
	}
	return TypeString
}

func (c *Context) Set(key string, value any) {
	switch v := value.(type) {
	case string:
		c.strings[key] = v
	case int:
		c.ints[key] = int64(v)
	case int32:
		c.ints[key] = int64(v)
	case int64:
		c.ints[key] = v
	case uint:
		c.uints[key] = uint64(v)
	case uint32:
		c.uints[key] = uint64(v)
	case uint64:
		c.uints[key] = v
	case float32:
		c.floats[key] = float64(v)
	case float64:
		c.floats[key] = v
	case bool:
		c.bools[key] = v
	default:
		c.objects[key] = v
	}
}

func (c *Context) SetLayered(key string, value any) {
	layerKey := c.makeLayerKey(key)
	c.Set(layerKey, value)
}

func (c *Context) SetString(key, value string) {
	c.strings[key] = value
}

func (c *Context) SetStringLayered(key, value string) {
	layerKey := c.makeLayerKey(key)
	c.strings[layerKey] = value
}

func (c *Context) SetInt(key string, value int64) {
	c.ints[key] = value
}

func (c *Context) SetIntLayered(key string, value int64) {
	layerKey := c.makeLayerKey(key)
	c.ints[layerKey] = value
}

func (c *Context) SetUint(key string, value uint64) {
	c.uints[key] = value
}

func (c *Context) SetUintLayered(key string, value uint64) {
	layerKey := c.makeLayerKey(key)
	c.uints[layerKey] = value
}

func (c *Context) SetFloat(key string, value float64) {
	c.floats[key] = value
}

func (c *Context) SetFloatLayered(key string, value float64) {
	layerKey := c.makeLayerKey(key)
	c.floats[layerKey] = value
}

func (c *Context) SetBool(key string, value bool) {
	c.bools[key] = value
}

func (c *Context) SetBoolLayered(key string, value bool) {
	layerKey := c.makeLayerKey(key)
	c.bools[layerKey] = value
}

func (c *Context) SetObject(key string, value any) {
	c.objects[key] = value
}

func (c *Context) SetObjectLayered(key string, value any) {
	layerKey := c.makeLayerKey(key)
	c.objects[layerKey] = value
}

func (c *Context) String(key string) string {
	return c.strings[key]
}

func (c *Context) StringLayered(key string) string {
	layerKey := c.makeLayerKey(key)
	return c.strings[layerKey]
}

func (c *Context) Int(key string) int64 {
	return c.ints[key]
}

func (c *Context) IntLayered(key string) int64 {
	layerKey := c.makeLayerKey(key)
	return c.ints[layerKey]
}

func (c *Context) Uint(key string) uint64 {
	return c.uints[key]
}

func (c *Context) UintLayered(key string) uint64 {
	layerKey := c.makeLayerKey(key)
	return c.uints[layerKey]
}

func (c *Context) Float(key string) float64 {
	return c.floats[key]
}

func (c *Context) FloatLayered(key string) float64 {
	layerKey := c.makeLayerKey(key)
	return c.floats[layerKey]
}

func (c *Context) Bool(key string) bool {
	return c.bools[key]
}

func (c *Context) BoolLayered(key string) bool {
	layerKey := c.makeLayerKey(key)
	return c.bools[layerKey]
}

func (c *Context) Object(key string) any {
	return c.objects[key]
}

func (c *Context) ObjectLayered(key string) any {
	layerKey := c.makeLayerKey(key)
	return c.objects[layerKey]
}

func (c *Context) GetString(key string) (string, bool) {
	v, ok := c.strings[key]
	return v, ok
}

func (c *Context) GetStringLayered(key string) (string, bool) {
	layerKey := c.makeLayerKey(key)
	v, ok := c.strings[layerKey]
	return v, ok
}

func (c *Context) GetInt(key string) (int64, bool) {
	v, ok := c.ints[key]
	return v, ok
}

func (c *Context) GetIntLayered(key string) (int64, bool) {
	layerKey := c.makeLayerKey(key)
	v, ok := c.ints[layerKey]
	return v, ok
}

func (c *Context) GetUint(key string) (uint64, bool) {
	v, ok := c.uints[key]
	return v, ok
}

func (c *Context) GetUintLayered(key string) (uint64, bool) {
	layerKey := c.makeLayerKey(key)
	v, ok := c.uints[layerKey]
	return v, ok
}

func (c *Context) GetFloat(key string) (float64, bool) {
	v, ok := c.floats[key]
	return v, ok
}

func (c *Context) GetFloatLayered(key string) (float64, bool) {
	layerKey := c.makeLayerKey(key)
	v, ok := c.floats[layerKey]
	return v, ok
}

func (c *Context) GetBool(key string) (bool, bool) {
	v, ok := c.bools[key]
	return v, ok
}

func (c *Context) GetBoolLayered(key string) (bool, bool) {
	layerKey := c.makeLayerKey(key)
	v, ok := c.bools[layerKey]
	return v, ok
}

func (c *Context) GetObject(key string) (any, bool) {
	v, ok := c.objects[key]
	return v, ok
}

func (c *Context) GetObjectLayered(key string) (any, bool) {
	layerKey := c.makeLayerKey(key)
	v, ok := c.objects[layerKey]
	return v, ok
}

func (c *Context) HasString(key string) bool {
	_, ok := c.strings[key]
	return ok
}

func (c *Context) HasStringLayered(key string) bool {
	layerKey := c.makeLayerKey(key)
	_, ok := c.strings[layerKey]
	return ok
}

func (c *Context) HasInt(key string) bool {
	_, ok := c.ints[key]
	return ok
}

func (c *Context) HasIntLayered(key string) bool {
	layerKey := c.makeLayerKey(key)
	_, ok := c.ints[layerKey]
	return ok
}

func (c *Context) HasUint(key string) bool {
	_, ok := c.uints[key]
	return ok
}

func (c *Context) HasUintLayered(key string) bool {
	layerKey := c.makeLayerKey(key)
	_, ok := c.uints[layerKey]
	return ok
}

func (c *Context) HasFloat(key string) bool {
	_, ok := c.floats[key]
	return ok
}

func (c *Context) HasFloatLayered(key string) bool {
	layerKey := c.makeLayerKey(key)
	_, ok := c.floats[layerKey]
	return ok
}

func (c *Context) HasBool(key string) bool {
	_, ok := c.bools[key]
	return ok
}

func (c *Context) HasBoolLayered(key string) bool {
	layerKey := c.makeLayerKey(key)
	_, ok := c.bools[layerKey]
	return ok
}

func (c *Context) HasObject(key string) bool {
	_, ok := c.objects[key]
	return ok
}

func (c *Context) HasObjectLayered(key string) bool {
	layerKey := c.makeLayerKey(key)
	_, ok := c.objects[layerKey]
	return ok
}

func (c *Context) Get(key string) (any, ValueType, bool) {
	if v, ok := c.strings[key]; ok {
		return v, TypeString, true
	}
	if v, ok := c.ints[key]; ok {
		return v, TypeInt, true
	}
	if v, ok := c.uints[key]; ok {
		return v, TypeUint, true
	}
	if v, ok := c.floats[key]; ok {
		return v, TypeFloat, true
	}
	if v, ok := c.bools[key]; ok {
		return v, TypeBool, true
	}
	if v, ok := c.objects[key]; ok {
		return v, TypeObject, true
	}
	return nil, TypeString, false
}

func (c *Context) GetLayered(key string) (any, ValueType, bool) {
	layerKey := c.makeLayerKey(key)
	return c.Get(layerKey)
}

func (c *Context) GetValueOfThisLayer(key string) any {
	layerKey := c.makeLayerKey(key)
	switch c.Type(layerKey) {
	case TypeString:
		return c.strings[layerKey]
	case TypeInt:
		return c.ints[layerKey]
	case TypeUint:
		return c.uints[layerKey]
	case TypeFloat:
		return c.floats[layerKey]
	case TypeBool:
		return c.bools[layerKey]
	case TypeObject:
		return c.objects[layerKey]
	default:
		return nil
	}
}

func (c *Context) DeleteString(key string) {
	delete(c.strings, key)
}

func (c *Context) DeleteStringLayered(key string) {
	layerKey := c.makeLayerKey(key)
	delete(c.strings, layerKey)
}

func (c *Context) DeleteInt(key string) {
	delete(c.ints, key)
}

func (c *Context) DeleteIntLayered(key string) {
	layerKey := c.makeLayerKey(key)
	delete(c.ints, layerKey)
}

func (c *Context) DeleteUint(key string) {
	delete(c.uints, key)
}

func (c *Context) DeleteUintLayered(key string) {
	layerKey := c.makeLayerKey(key)
	delete(c.uints, layerKey)
}

func (c *Context) DeleteFloat(key string) {
	delete(c.floats, key)
}

func (c *Context) DeleteFloatLayered(key string) {
	layerKey := c.makeLayerKey(key)
	delete(c.floats, layerKey)
}

func (c *Context) DeleteBool(key string) {
	delete(c.bools, key)
}

func (c *Context) DeleteBoolLayered(key string) {
	layerKey := c.makeLayerKey(key)
	delete(c.bools, layerKey)
}

func (c *Context) DeleteObject(key string) {
	delete(c.objects, key)
}

func (c *Context) DeleteObjectLayered(key string) {
	layerKey := c.makeLayerKey(key)
	delete(c.objects, layerKey)
}

func (c *Context) Delete(key string) {
	delete(c.strings, key)
	delete(c.ints, key)
	delete(c.uints, key)
	delete(c.floats, key)
	delete(c.bools, key)
	delete(c.objects, key)
}

func (c *Context) DeleteLayered(key string) {
	layerKey := c.makeLayerKey(key)
	c.Delete(layerKey)
}

func (c *Context) StringKeys() []string {
	keys := make([]string, 0, len(c.strings))
	for k := range c.strings {
		keys = append(keys, k)
	}
	return keys
}

func (c *Context) IntKeys() []string {
	keys := make([]string, 0, len(c.ints))
	for k := range c.ints {
		keys = append(keys, k)
	}
	return keys
}

func (c *Context) UintKeys() []string {
	keys := make([]string, 0, len(c.uints))
	for k := range c.uints {
		keys = append(keys, k)
	}
	return keys
}

func (c *Context) FloatKeys() []string {
	keys := make([]string, 0, len(c.floats))
	for k := range c.floats {
		keys = append(keys, k)
	}
	return keys
}

func (c *Context) BoolKeys() []string {
	keys := make([]string, 0, len(c.bools))
	for k := range c.bools {
		keys = append(keys, k)
	}
	return keys
}

func (c *Context) ObjectKeys() []string {
	keys := make([]string, 0, len(c.objects))
	for k := range c.objects {
		keys = append(keys, k)
	}
	return keys
}

func (c *Context) AllKeys() []string {
	allKeys := make(map[string]bool)
	for k := range c.strings {
		allKeys[k] = true
	}
	for k := range c.ints {
		allKeys[k] = true
	}
	for k := range c.uints {
		allKeys[k] = true
	}
	for k := range c.floats {
		allKeys[k] = true
	}
	for k := range c.bools {
		allKeys[k] = true
	}
	for k := range c.objects {
		allKeys[k] = true
	}

	keys := make([]string, 0, len(allKeys))
	for k := range allKeys {
		keys = append(keys, k)
	}
	return keys
}

func (c *Context) CopyStrings() map[string]string {
	copy := make(map[string]string, len(c.strings))
	for k, v := range c.strings {
		copy[k] = v
	}
	return copy
}

func (c *Context) CopyInts() map[string]int64 {
	copy := make(map[string]int64, len(c.ints))
	for k, v := range c.ints {
		copy[k] = v
	}
	return copy
}

func (c *Context) CopyUints() map[string]uint64 {
	copy := make(map[string]uint64, len(c.uints))
	for k, v := range c.uints {
		copy[k] = v
	}
	return copy
}

func (c *Context) CopyFloats() map[string]float64 {
	copy := make(map[string]float64, len(c.floats))
	for k, v := range c.floats {
		copy[k] = v
	}
	return copy
}

func (c *Context) CopyBools() map[string]bool {
	copy := make(map[string]bool, len(c.bools))
	for k, v := range c.bools {
		copy[k] = v
	}
	return copy
}

func (c *Context) CopyObjects() map[string]any {
	copy := make(map[string]any, len(c.objects))
	for k, v := range c.objects {
		copy[k] = v
	}
	return copy
}

func (c *Context) MergeStrings(other map[string]string) {
	for k, v := range other {
		c.strings[k] = v
	}
}

func (c *Context) MergeInts(other map[string]int64) {
	for k, v := range other {
		c.ints[k] = v
	}
}

func (c *Context) MergeUints(other map[string]uint64) {
	for k, v := range other {
		c.uints[k] = v
	}
}

func (c *Context) MergeFloats(other map[string]float64) {
	for k, v := range other {
		c.floats[k] = v
	}
}

func (c *Context) MergeBools(other map[string]bool) {
	for k, v := range other {
		c.bools[k] = v
	}
}

func (c *Context) MergeObjects(other map[string]any) {
	for k, v := range other {
		c.objects[k] = v
	}
}
