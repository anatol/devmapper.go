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

func TestFixedArrayToString(t *testing.T) {
	t.Parallel()

	check := func(input []byte, expected string) {
		str := fixedArrayToString(input)
		if str != expected {
			t.Fatalf("Expected string %v, got %v", expected, str)
		}
	}

	check([]byte{}, "")
	check([]byte{'r'}, "r")
	check([]byte{'h', 'e', 'l', 'l', 'o', ',', ' '}, "hello, ")
	check([]byte{'h', '\x00', 'l', 'l', 'o', ',', ' '}, "h")
	check([]byte{'\x00'}, "")
}
