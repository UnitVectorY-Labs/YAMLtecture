package configuration

import (
	"testing"
)

func TestLoadYAMLFiles(t *testing.T) {
	config, err := LoadConfig("../../example/simple/architecture.yaml")
	if err != nil {
		t.Fatalf("Failed to load YAML files: %v", err)
	}

	// Assert nodes are loaded correctly
	expectedNodes := map[string]Node{
		"cluster": {
			ID:   "cluster",
			Type: "Infrastructure",
			Attributes: map[string]interface{}{
				"name": "Container Hosting",
			},
		},
		"service_foo": {
			ID:     "service_foo",
			Type:   "Microservice",
			Parent: "cluster",
			Attributes: map[string]interface{}{
				"name":     "Foo Service",
				"language": "Java",
			},
		},
		"service_bar": {
			ID:     "service_bar",
			Type:   "Microservice",
			Parent: "cluster",
			Attributes: map[string]interface{}{
				"name":     "Bar Service",
				"language": "Go",
			},
		},
	}
	if len(config.Nodes) != len(expectedNodes) {
		t.Errorf("Expected %d nodes, got %d", len(expectedNodes), len(config.Nodes))
	}
	for id, expected := range expectedNodes {
		actual, exists := config.Nodes[id]
		if !exists {
			t.Errorf("Expected node with ID '%s' not found", id)
			continue
		}
		if actual.Type != expected.Type {
			t.Errorf("Expected node '%s' type '%s', got '%s'", id, expected.Type, actual.Type)
		}
		if actual.Parent != expected.Parent {
			t.Errorf("Expected node '%s' parent '%s', got '%s'", id, expected.Parent, actual.Parent)
		}
		for key, val := range expected.Attributes {
			if actual.Attributes[key] != val {
				t.Errorf("Expected node '%s' attribute '%s'='%v', got '%v'", id, key, val, actual.Attributes[key])
			}
		}
	}

	// Assert links are loaded correctly
	expectedLinks := []Link{
		{
			Source: "service_foo",
			Target: "service_bar",
			Type:   "API",
			Attributes: map[string]interface{}{
				"payload": "example",
			},
		},
	}
	if len(config.Links) != len(expectedLinks) {
		t.Errorf("Expected %d links, got %d", len(expectedLinks), len(config.Links))
	}
	for i, expected := range expectedLinks {
		actual := config.Links[i]
		if actual.Source != expected.Source || actual.Target != expected.Target || actual.Type != expected.Type {
			t.Errorf("Expected link %v, got %v", expected, actual)
		}
		for key, val := range expected.Attributes {
			if actual.Attributes[key] != val {
				t.Errorf("Expected link attribute '%s'='%v', got '%v'", key, val, actual.Attributes[key])
			}
		}
	}
}
