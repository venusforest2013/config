package config

import (
	"time"
)

type Value struct {
	Node
}

func (v *Value) Has(key string) bool {
	if key != "" {
		return (v.Access(key) != nil)
	}
	return false
}

func (v *Value) HasArray(key string) bool {
	if key != "" {
		return (v.AccessArray(key) != nil)
	}
	return false
}

func (v *Value) HasMap(key string) bool {
	if key != "" {
		return (v.AccessMap(key) != nil)
	}
	return false
}

func (v *Value) Bool(key string, def bool) bool {
	if key != "" {
		if t := v.Access(key); t != nil {
			if r, ok := t.Value().(bool); ok {
				return r
			}
		}
	}
	return def
}

func (v *Value) Str(key, def string) string {
	if key != "" {
		if t := v.Access(key); t != nil {
			if r, ok := t.Value().(string); ok {
				return r
			}
		}
	}
	return def
}

func (v *Value) Int64(key string, def int64) int64 {
	if key != "" {
		if t := v.Access(key); t != nil {
			if r, ok := t.Value().(int64); ok {
				return r
			}
		}
	}
	return def
}

func (v *Value) Float64(key string, def float64) float64 {
	if key != "" {
		if t := v.Access(key); t != nil {
			if r, ok := t.Value().(float64); ok {
				return r
			}
		}
	}
	return def
}

func (v *Value) Duration(key string, def time.Duration) time.Duration {
	if key != "" {
		if t := v.Access(key); t != nil {
			if str, ok := t.Value().(string); ok {
				if duration, err := time.ParseDuration(str); err == nil {
					return duration
				}
			}
		}
	}
	return def
}

func (v *Value) BoolArray(key string) []bool {
	if key != "" {
		if t := v.AccessArray(key); t != nil {
			r := make([]bool, 0, len(t))
			for _, node := range t {
				if it, ok := node.Value().(bool); ok {
					r = append(r, it)
				}
			}
			return r
		}
	}
	return nil
}

func (v *Value) StrArray(key string) []string {
	if key != "" {
		if t := v.AccessArray(key); t != nil {
			r := make([]string, 0, len(t))
			for _, node := range t {
				if it, ok := node.Value().(string); ok {
					r = append(r, it)
				}
			}
			return r
		}
	}
	return nil
}

func (v *Value) Int64Array(key string) []int64 {
	if key != "" {
		if t := v.AccessArray(key); t != nil {
			r := make([]int64, 0, len(t))
			for _, node := range t {
				if it, ok := node.Value().(int64); ok {
					r = append(r, it)
				}
			}
			return r
		}
	}
	return nil
}

func (v *Value) Float64Array(key string) []float64 {
	if key != "" {
		if t := v.AccessArray(key); t != nil {
			r := make([]float64, 0, len(t))
			for _, node := range t {
				if it, ok := node.Value().(float64); ok {
					r = append(r, it)
				}
			}
			return r
		}
	}
	return nil
}

func (v *Value) ValueArray(key string) []*Value {
	if t := v.AccessArray(key); t != nil {
		r := make([]*Value, 0, len(t))
		for _, node := range t {
			r = append(r, &Value{node})
		}
		return r
	}
	return nil
}

func (v *Value) ValueMap(key string) map[string]*Value {
	if t := v.AccessMap(key); t != nil {
		r := map[string]*Value{}
		for k, node := range t {
			r[k] = &Value{node}
		}
		return r
	}
	return nil
}
