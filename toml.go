package config

import (
	"fmt"
	toml "github.com/pelletier/go-toml"
	"strings"
)

type Toml struct {
	value interface{}
}

func LoadTomlFile(filename string) (*Toml, error) {
	cfg, err := toml.LoadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Toml{cfg}, nil
}

func LoadToml(content string) (*Toml, error) {
	cfg, err := toml.Load(content)
	if err != nil {
		return nil, err
	}
	return &Toml{cfg}, nil
}

func (t *Toml) Type() Type {
	return TomlType
}

func (t *Toml) Value() interface{} {
	return t.value
}

func (t *Toml) Access(key string) Node {
	tree, ok := t.value.(*toml.Tree)
	if !ok {
		return nil
	}
	value := tree.GetPath(strings.Split(key, Delimiter))
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case *toml.Tree:
		return &Toml{v}
	case []*toml.Tree:
		if len(v) == 0 {
			return nil
		}
		return &Toml{v[len(v)-1]}
	default:
		return &Toml{v}
	}
	return nil
}

func (t *Toml) AccessArray(key string) []Node {
	tree, ok := t.value.(*toml.Tree)
	if !ok {
		return nil
	}
	values := tree.GetArrayPath(strings.Split(key, Delimiter))
	if values == nil {
		return nil
	}
	r := []Node{}
	switch v := values.(type) {

	case []*toml.Tree:
		for _, tree := range v {
			r = append(r, &Toml{tree})
		}
		return r

	case []interface{}:
		for _, item := range v {
			r = append(r, &Toml{item})
		}
		return r

	case []bool:
		for _, item := range v {
			r = append(r, &Toml{item})
		}
		return r

	case []string:
		for _, item := range v {
			r = append(r, &Toml{item})
		}
		return r

	case []int64:
		for _, item := range v {
			r = append(r, &Toml{item})
		}
		return r

	case []float64:
		for _, item := range v {
			r = append(r, &Toml{item})
		}
		return r

	}

	return nil
}

func (t *Toml) AccessMap(key string) map[string]Node {
	tree, ok := t.value.(*toml.Tree)
	if !ok {
		return nil
	}
	value := tree.GetPath(strings.Split(key, Delimiter))
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case *toml.Tree:
		r := map[string]Node{}
		for _, key := range v.Keys() {
			r[key] = &Toml{v.Get(key)}
		}
		return r
	case []*toml.Tree:
		return nil
	}
	return nil
}

func (t *Toml) String() string {
	switch v := t.value.(type) {
	case *toml.Tree:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}
