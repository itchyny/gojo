package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCliRun(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		input    string
		expected string
		err      string
	}{
		{
			name: "default",
			args: []string{"foo=bar"},
			expected: `{"foo":"bar"}
`,
		},
		{
			name: "multiple",
			args: []string{"foo=bar", "bar=baz", "qux=quux"},
			expected: `{"foo":"bar","bar":"baz","qux":"quux"}
`,
		},
		{
			name: "pretty",
			args: []string{"-p", "foo=bar", "bar=baz", "qux=quux"},
			expected: `{
  "foo": "bar",
  "bar": "baz",
  "qux": "quux"
}
`,
		},
		{
			name: "nested object",
			args: []string{"-p", `foo={"bar":{"foo":1,"baz":"qux","quux":["foo"]}}`},
			expected: `{
  "foo": {
    "bar": {
      "foo": 1,
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
			name: "array",
			args: []string{"-a", "foo", "bar", "baz", "false", "0x40"},
			expected: `["foo","bar","baz",false,64]
`,
		},
		{
			name: "array pretty",
			args: []string{"-a", "-p", "foo", "0xf0", `{"foo":{"bar":30}}`},
			expected: `[
  "foo",
  240,
  {
    "foo": {
      "bar": 30
    }
  }
]
`,
		},
		{
			name: "hyphen hyphen",
			args: []string{"-a", "--", "-p", "foo"},
			expected: `["-p","foo"]
`,
		},
		{
			name: "input",
			args: []string{"-a", "-p"},
			input: `foo
{"bar":100}
`,
			expected: `[
  "foo",
  {
    "bar": 100
  }
]
`,
		},
		{
			name: "input error",
			input: `foo
bar
`,
			err: `failed to parse: "foo"`,
		},
		{
			name: "parse error",
			args: []string{"foo"},
			err:  `failed to parse: "foo"`,
		},
		{
			name: "unkown flag",
			args: []string{"-b"},
			err:  `flag provided but not defined: -b`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			outStream := new(bytes.Buffer)
			errStream := new(bytes.Buffer)
			cli := cli{
				inStream:  strings.NewReader(tc.input),
				outStream: outStream,
				errStream: errStream,
			}
			code := cli.run(tc.args)
			if tc.err == "" {
				assert.Equal(t, exitCodeOK, code)
				assert.Equal(t, tc.expected, outStream.String())
			} else {
				assert.Equal(t, exitCodeErr, code)
				assert.Contains(t, errStream.String(), tc.err)
			}
		})
	}
}
