package query

import (
	"os"

	"gopkg.in/yaml.v3"
)

// ParseQuery parses a query from a YAML string.
func ParseQuery(queryStr string) (*Query, error) {
	var query Query
	err := yaml.Unmarshal([]byte(queryStr), &query)
	if err != nil {
		return nil, err
	}
	return &query, nil
}

// LoadQuery loads the query from a YAML file.
func LoadQuery(filepath string) (*Query, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var query Query
	err = yaml.Unmarshal(data, &query)
	if err != nil {
		return nil, err
	}
	return &query, nil
}
