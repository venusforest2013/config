package config

import (
	"testing"
	"time"
)

var (
	testValue *Value
)

func TestValue_Init(t *testing.T) {
	TestToml_LoadFile(t)

	testValue = &Value{testToml}
}

func TestValue_Has(t *testing.T) {
	if !testValue.Has("test.server") {
		t.Fatal()
	}
	if testValue.Has("string.str3") {
		t.Fatal()
	}
}

func TestValue_HasArray(t *testing.T) {
	if !testValue.HasArray("products") {
		t.Fatal()
	}
	if testValue.HasArray("boolean_array") {
		t.Fatal()
	}
	if testValue.HasArray("boolean_array.key.true") {
		t.Fatal()
	}
}

func TestValue_HasMap(t *testing.T) {
	if !testValue.HasMap("test.server") {
		t.Fatal()
	}
	if testValue.HasMap("test.client.tt") {
		t.Fatal()
	}
	if testValue.HasMap("test.server.timeout.key1") {
		t.Fatal()
	}
}

func TestValue_Value(t *testing.T) {
	if !testValue.Bool("boolean.True", false) {
		t.Fatal()
	}
	if testValue.Bool("boolean.False", true) {
		t.Fatal()
	}
	if testValue.Bool("boolean.nil", false) {
		t.Fatal()
	}
	if v := testValue.Float64("test.client.name.float_val", 1.0); v != float64(3.1415926) {
		t.Fatal(v)
	}
	if v := testValue.Float64("test.client.name.float2_val", 1.23); v != float64(1.23) {
		t.Fatal(v)
	}
	if v := testValue.Duration("test.golang.timeout", 100*time.Millisecond); v != 300*time.Millisecond {
		t.Fatal(v)
	}
	if v := testValue.Duration("test.golang.timeout2", 100*time.Millisecond); v != 100*time.Millisecond {
		t.Fatal(v)
	}
}

func TestValue_ValueArray(t *testing.T) {
	v := testValue.BoolArray("boolean_array.key")
	if len(v) != 4 || v[2] {
		t.Fatal(v)
	}

	if testValue.BoolArray("boolean_array.key2") != nil {
		t.Fatal()
	}

	if arr := testValue.ValueArray("string_array.key2"); arr == nil {
		t.Fatal(arr)
	} else {
		if len(arr) != 3 {
			t.Fatal(arr)
		}
		if tv, ok := arr[1].Value().(string); !ok || tv != "pear" {
			t.Fatal(tv)
		}
		if tv, ok := arr[2].Value().(string); !ok || tv != "banana" {
			t.Fatal(tv)
		}
	}
}

func TestValue_ValueMap(t *testing.T) {
	vmap := testValue.ValueMap("test.server")
	if vmap == nil || len(vmap) != 2 {
		t.Fatal(vmap)
	}

	item, ok := vmap["name"]
	if !ok {
		t.Fatal(vmap)
	}
	if item.Str("key1", "str") != "str1" {
		t.Fatal(item)
	}
}
