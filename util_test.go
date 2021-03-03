package devmapper

import "testing"

func TestRundup(t *testing.T) {
	compare := func(expected, got int) {
		if expected != got {
			t.Fatalf("expected %v, got %v", expected, got)
		}
	}

	compare(0, roundUp(0, 8))
	compare(8, roundUp(1, 8))
}
