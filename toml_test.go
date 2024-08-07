package config

import (
	"testing"
)

var (
	testToml *Toml

	testTomlFile = "./test.toml"
)

func init() {
	toml, err := LoadTomlFile(testTomlFile)
	if err != nil {
		panic(err)
	}
	testToml = toml
}

func TestToml_LoadFile(t *testing.T) {
	toml, err := LoadTomlFile(testTomlFile)
	if err != nil {
		t.Fatal(err)
	}
	testToml = toml
	if testToml.Type() != TomlType {
		t.Fatal(nil)
	}
}

func TestToml_Access(t *testing.T) {
	var (
		node Node
	)

	node = testToml.Access("proc")
	if node == nil {
		t.Fatal(nil)
	}
	if v, ok := node.Value().(int64); !ok || v != 10 {
		t.Fatal(nil)
	}

	node = testToml.Access("string.str2")
	if node == nil {
		t.Fatal(nil)
	}
	if v, ok := node.Value().(string); !ok || v != "world" {
		t.Fatal(nil)
	}

	node = testToml.Access("string.str3")
	if node != nil {
		t.Fatal(nil)
	}

	node = testToml.Access("products")
	if node == nil {
		t.Fatal(nil)
	}
}

func TestToml_AccessArray(t *testing.T) {
	var (
		node  Node
		nodes []Node
	)
	nodes = testToml.AccessArray("products")
	if len(nodes) != 2 {
		t.Fatal(nodes)
	}

	node = nodes[0].Access("name")
	if v, ok := node.Value().(string); !ok || v != "Hammer" {
		t.Fatal(v)
	}
	node = nodes[1].Access("name")
	if v, ok := node.Value().(string); !ok || v != "Nail" {
		t.Fatal(v)
	}

	nodes = testToml.AccessArray("boolean_array.key")
	if len(nodes) != 4 {
		t.Fatal(nodes)
	}
}

func TestToml_AccessMap(t *testing.T) {
	var (
		nodeMap map[string]Node
	)
	nodeMap = testToml.AccessMap("boolean")
	if len(nodeMap) != 2 {
		t.Fatal(nodeMap)
	}
	if v, ok := nodeMap["True"]; !ok {
		t.Fatal(v)
	} else {
		if tv, ok := v.Value().(bool); !ok || !tv {
			t.Fatal(tv)
		}
	}
	if v, ok := nodeMap["False"]; !ok {
		t.Fatal(v)
	} else {
		if tv, ok := v.Value().(bool); !ok || tv {
			t.Fatal(tv)
		}
	}

	nodeMap = testToml.AccessMap("products.name")
	if nodeMap != nil {
		t.Fatal(nodeMap)
	}

	nodeMap = testToml.AccessMap("test.server")
	if nodeMap == nil {
		t.Fatal(nodeMap)
	}
	if len(nodeMap) != 2 {
		t.Fatal(nodeMap)
	}
	if v, ok := nodeMap["timeout"]; !ok {
		t.Fatal(v)
	} else {
		node := v.Access("key2")
		if node == nil {
			t.Fatal(v)
		}
		if tv, ok := node.Value().(int64); !ok || tv != 123 {
			t.Fatal(tv)
		}
	}

}
