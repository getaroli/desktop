package profiles

import "testing"

func TestProfileByName(t *testing.T) {
	p, err := profileByName("creator")
	if err != nil || len(p.Plugins) == 0 {
		t.Fatalf("creator profile = %#v, %v", p, err)
	}
	if _, err := profileByName("does-not-exist"); err == nil {
		t.Fatal("unknown profile was accepted")
	}
}
