package fileformat

import (
	"fmt"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in        string
		shouldErr bool
	}{
		{"does-not-exist", true},
		{"bad-runbook-missing-info.yml", true},
		{"bad-runbook-schema.yml", true},
		{"make-mod-del.yml", false},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			expandedPath := fmt.Sprintf("../../examples/runbooks/%s", test.in)
			_, err := Validate(expandedPath)
			if test.shouldErr && err == nil {
				t.Errorf("error expected for test case: %s", test.in)
			} else if !test.shouldErr && err != nil {
				t.Errorf("no error expected for test case: %s \n err: %s", test.in, err)
			}
		})
	}
}
