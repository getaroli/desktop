package main

import (
	"strings"
	"testing"
)

func TestRunRejectsUnknownCommand(t *testing.T) {
	err := run([]string{"not-a-command"})
	if err == nil || !strings.Contains(err.Error(), "comando desconhecido") {
		t.Fatalf("unexpected error: %v", err)
	}
}
