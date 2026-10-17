package seed

import "testing"

func TestCatalog(t *testing.T) {
	courses, err := Catalog()
	if err != nil {
		t.Fatalf("Catalog() failed: %v", err)
	}
	if len(courses) != 56 {
		t.Fatalf("expected 56 courses, got %d", len(courses))
	}
}

