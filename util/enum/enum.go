package enum

import "errors"

type Wrapper struct {
	TypeName   string
	value2name map[int]string
	name2value map[string]int
}

type Item struct {
	Value int
	Name  string
}

const (
	UnknownValue = -1
)

func NewEnumWrapper(typeName string, enums ...Item) *Wrapper {
	e := new(Wrapper)
	e.TypeName = typeName

	e.value2name = map[int]string{}
	e.name2value = map[string]int{}
	for _, item := range enums {
		e.value2name[item.Value] = item.Name
		e.name2value[item.Name] = item.Value
	}
	return e
}

func (e *Wrapper) GetName(value int) string {
	if name, ok := e.value2name[value]; ok {
		return name
	}
	return ""
}

func (e *Wrapper) GetValue(name string) int {
	if value, ok := e.name2value[name]; ok {
		return value
	}
	return UnknownValue
}

func (e *Wrapper) MarshalJSONInner(value int) ([]byte, error) {
	return []byte(`"` + e.GetName(value) + `"`), nil
}

func (e *Wrapper) UnmarshalJSONInner(data []byte) (value int, err error) {
	if v, ok := e.name2value[string(data)]; ok {
		value = v
		return
	}
	return UnknownValue, errors.New("unsupported value")
}
