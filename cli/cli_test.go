package cli

import (
	"strings"
	"testing"
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
			args: []string{"foo=bar", "bar=false", "baz=qux", "\n=", `\n=null`},
			expected: `{"foo":"bar","bar":false,"baz":"qux","\n":"","\\n":null}
`,
		},
		{
			name: "pretty",
			args: []string{"-p", "foo=bar", "bar=true", "baz=qux", "qux=null"},
			expected: `{
  "foo": "bar",
  "bar": true,
  "baz": "qux",
  "qux": null
}
`,
		},
		{
			name: "nested object",
			args: []string{"-p", `foo={"bar":{"foo":1,"baz":"qux","quux":["foo",{},[]]}}`},
			expected: `{
  "foo": {
    "bar": {
      "foo": 1,
      "baz": "qux",
      "quux": [
        "foo",
        {},
        []
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
			name: "yaml",
			args: []string{"-y", "foo=bar"},
			expected: `foo: bar
`,
		},
		{
			name: "yaml multiple",
			args: []string{"-y", "foo=bar", "bar=false", "baz=qux", "\n=", `\n=null`},
			expected: `foo: bar
bar: false
baz: qux
"\n": ""
\n: null
`,
		},
		{
			name: "yaml nested object",
			args: []string{"-y", `foo={"bar":{"foo":1,"baz":"qux","quux":["foo",{},[],{"bar":{}},{"baz":[]},0]}}`},
			expected: `foo:
  bar:
    foo: 1
    baz: qux
    quux:
      - foo
      - {}
      - []
      - bar: {}
      - baz: []
      - 0
`,
		},
		{
			name: "yaml array",
			args: []string{"-y", "-a", "foo", "bar", "baz", "false", "0x40"},
			expected: `- foo
- bar
- baz
- false
- 64
`,
		},
		{
			name: "yaml nested array",
			args: []string{"-y", "-a", `{"foo":{"bar":["a\nb\n",{"a\nb":"\na\nb\n"},{"baz":""}],"a\n":"\n"}}`},
			expected: `- foo:
    bar:
      - |
        a
        b
      - ? |-
          a
          b
        : |

          a
          b
      - baz: ""
    ? |
      a
    : "\n"
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
			name: "unknown flag",
			args: []string{"-b"},
			err:  `flag provided but not defined: -b`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var outStream, errStream strings.Builder
			cli := cli{
				inStream:  strings.NewReader(tc.input),
				outStream: &outStream,
				errStream: &errStream,
			}
			code := cli.run(tc.args)
			if tc.err == "" {
				if code != exitCodeOK {
					t.Errorf("code should be %d but got %d", exitCodeOK, code)
				}
				if got := outStream.String(); got != tc.expected {
					t.Errorf("output should be\n  %q\nbut got\n  %q", tc.expected, got)
				}
			} else {
				if code != exitCodeErr {
					t.Errorf("code should be %d but got %d", exitCodeErr, code)
				}
				if got := errStream.String(); !strings.Contains(got, tc.expected) {
					t.Errorf("error output should contain %q but got %q", tc.err, got)
				}
			}
		})
	}
}
