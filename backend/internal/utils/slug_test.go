package utils

import "testing"

func TestSlugify(t *testing.T) {
	cases := map[string]string{
		"Hello World":     "hello-world",
		"  Multiple   ":    "multiple",
		"Symbols!@#$%^&*": "symbols",
		"":                "book",
	}
	for input, expected := range cases {
		if got := Slugify(input); got != expected {
			t.Fatalf("slugify(%q) = %q, want %q", input, got, expected)
		}
	}
}
