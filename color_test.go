// Copyright © 2021 Hidetatsu Yaginuma. All rights reserved.
package skew

import "testing"

func TestColor(t *testing.T) {
	if expected, got := "[33mtest[0m", yellow("test"); expected != got {
		t.Errorf("yellow: expected: %s, got: %s", expected, got)
	}

	if expected, got := "[32mtest[0m", green("test"); expected != got {
		t.Errorf("green: expected: %s, got: %s", expected, got)
	}

	if expected, got := "[31mtest[0m", red("test"); expected != got {
		t.Errorf("red: expected: %s, got: %s", expected, got)
	}
}
