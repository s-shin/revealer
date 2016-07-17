package main

import "testing"

func TestReplaceEmoji(t *testing.T) {
	expected := "images/emoji/unicode/1f604.png"
	replaced := ReplaceEmoji(":smile:", func(path string) string {
		if path != expected {
			t.Errorf("expected %s but %s", expected, path)
		}
		return path
	})
	if replaced != expected {
		t.Errorf("expected %s but %s", expected, replaced)
	}
}
