package query

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Query struct {
	Nodes Nodes `yaml:"nodes"`
	Links Links `yaml:"links"`
}

type Nodes struct {
	Filters []Filter `yaml:"filters"`
}

type Links struct {
	Filters []Filter `yaml:"filters"`
}

type Filter struct {
	Condition Condition `yaml:"condition"`
}

type Condition struct {
	Field      string      `yaml:"field"`
	Operator   string      `yaml:"operator"`
	Value      string      `yaml:"value"`
	Conditions []Condition `yaml:"conditions"`
}

// YamlString returns the YAML representation of the query
func (q *Query) YamlString() string {
	// Marshall the config to a string
	data, err := yaml.Marshal(q)
	if err != nil {
		return fmt.Sprintf("error marshalling query: %v", err)
	}
	return string(data)
}
