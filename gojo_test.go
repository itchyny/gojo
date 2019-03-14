package gojo

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGojoRun(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		pretty   bool
		expected string
	}{
		{
			name: "default",
			args: []string{"foo=bar"},
			expected: `{"foo":"bar"}
`,
		},
		{
			name: "multiple",
			args: []string{"foo=bar", "bar=baz", "baz=qux"},
			expected: `{"bar":"baz","baz":"qux","foo":"bar"}
`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			assert.NoError(t, New(tc.args, out).Run())
			assert.Equal(t, tc.expected, out.String())
		})
	}
}
