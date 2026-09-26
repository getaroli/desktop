package power

import "testing"

func TestBatteryRejectsUnknownAction(t *testing.T) {
	if err := Battery([]string{"turbo"}); err == nil {
		t.Fatal("unknown battery action was accepted")
	}
}
