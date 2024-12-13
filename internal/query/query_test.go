package query

import (
	"testing"
)

func TestLoadQuery(t *testing.T) {
	query, err := LoadQuery("../../example/queries/simple.yaml")
	if err != nil {
		t.Fatalf("Failed to load query: %v", err)
	}

	err = query.Validate()
	if err != nil {
		t.Fatalf("Query did not pass validation: %v", err)
	}
}
