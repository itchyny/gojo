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
		array    bool
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
		{
			name:   "pretty",
			args:   []string{"foo=bar", "bar=baz", "baz=qux"},
			pretty: true,
			expected: `{
  "bar": "baz",
  "baz": "qux",
  "foo": "bar"
}
`,
		},
		{
			name:  "array",
			args:  []string{"foo", "bar", "baz"},
			array: true,
			expected: `["foo","bar","baz"]
`,
		},
		{
			name:   "array pretty",
			args:   []string{"foo", "bar", "baz"},
			array:  true,
			pretty: true,
			expected: `[
  "foo",
  "bar",
  "baz"
]
`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			opts := []Option{
				Args(tc.args),
				OutStream(out),
			}
			if tc.array {
				opts = append(opts, Array())
			}
			if tc.pretty {
				opts = append(opts, Pretty())
			}
			assert.NoError(t, New(opts...).Run())
			assert.Equal(t, tc.expected, out.String())
		})
	}
}
