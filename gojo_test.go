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
			args: []string{"foo=bar", "bar=false", "baz=qux"},
			expected: `{"bar":false,"baz":"qux","foo":"bar"}
`,
		},
		{
			name:   "pretty",
			args:   []string{"foo=bar", "bar=true", "baz=qux"},
			pretty: true,
			expected: `{
  "bar": true,
  "baz": "qux",
  "foo": "bar"
}
`,
		},
		{
			name:   "numbers",
			args:   []string{"a=123", "b=3.14", "c=3e10", "d=-128"},
			pretty: true,
			expected: `{
  "a": 123,
  "b": 3.14,
  "c": 30000000000,
  "d": -128
}
`,
		},
		{
			name: "nested object",
			args: []string{`foo={"bar":{"baz":"qux","quux":["foo"]}}`},
			expected: `{"foo":{"bar":{"baz":"qux","quux":["foo"]}}}
`,
		},
		{
			name:   "nested object pretty",
			args:   []string{`foo={"bar":{"baz":"qux","quux":["foo"]}}`},
			pretty: true,
			expected: `{
  "foo": {
    "bar": {
      "baz": "qux",
      "quux": [
        "foo"
      ]
    }
  }
}
`,
		},
		{
			name:  "array",
			args:  []string{"foo", "false", "baz"},
			array: true,
			expected: `["foo",false,"baz"]
`,
		},
		{
			name:   "array pretty",
			args:   []string{"foo", "true", "baz"},
			array:  true,
			pretty: true,
			expected: `[
  "foo",
  true,
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
