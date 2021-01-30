package gojo

import (
	"strings"
	"testing"
)

func TestGojoRun(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		array    bool
		pretty   bool
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
			args: []string{"foo=bar", "bar=false", "baz=qux", "\n=quux", `\n=foo`},
			expected: `{"foo":"bar","bar":false,"baz":"qux","\n":"quux","\\n":"foo"}
`,
		},
		{
			name:   "pretty",
			args:   []string{"foo=bar", "bar=true", "baz=qux", "qux=null"},
			pretty: true,
			expected: `{
  "foo": "bar",
  "bar": true,
  "baz": "qux",
  "qux": null
}
`,
		},
		{
			name:   "numbers",
			args:   []string{"a=123", "c=3.14", "d=3e10", "b=-128", "e=[1,2,3]", "f=0xffdc", "g=0XFF", "h=037"},
			pretty: true,
			expected: `{
  "a": 123,
  "c": 3.14,
  "d": 30000000000,
  "b": -128,
  "e": [
    1,
    2,
    3
  ],
  "f": 65500,
  "g": 255,
  "h": 37
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
			args:   []string{`foo={"bar":{"foo":1,"baz":"qux","quux":["foo", []]}}`},
			pretty: true,
			expected: `{
  "foo": {
    "bar": {
      "foo": 1,
      "baz": "qux",
      "quux": [
        "foo",
        []
      ]
    }
  }
}
`,
		},
		{
			name:  "array",
			args:  []string{"foo", "false", "baz", "null"},
			array: true,
			expected: `["foo",false,"baz",null]
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
		{
			name: "nested keys",
			args: []string{`foo[]=1`, `foo[]=bar`, `bar[baz]=10`, `qux][=20`, `a[b]c=d`},
			expected: `{"foo":[1,"bar"],"bar":{"baz":10},"qux][":20,"a[b]c":"d"}
`,
		},
		{
			name: "deep keys",
			args: []string{`a[b][c]d=e`, `a[b][c][d]=f`, `b[c][d][][]=f`, `b[c][d][]=g`, `c[][][]=d`, `[][][]=h`, `[][1][2]=i`},
			expected: `{"a[b][c]d":"e","a":{"b":{"c":{"d":"f"}}},"b":{"c":{"d":[["f"],"g"]}},"c":[[["d"]]],"":[[["h"]],{"1":{"2":"i"}}]}
`,
		},
		{
			name: "merge to json",
			args: []string{`a={"b":{"c":10,"d":[20]}}`, `a[b][d][]=30`, `a[b][e]=40`},
			expected: `{"a":{"b":{"c":10,"d":[20,30],"e":40}}}
`,
		},
		{
			name: "parse error",
			args: []string{`foo`},
			err:  `failed to parse: "foo"`,
		},
		{
			name: "expected object",
			args: []string{`foo[]=1`, `foo[bar]=2`},
			err:  `expected an object: foo: [1]`,
		},
		{
			name: "expected object deep",
			args: []string{`foo[bar][]=1`, `foo[bar][baz]=2`},
			err:  `expected an object: foo.bar: [1]`,
		},
		{
			name: "expected array",
			args: []string{`foo[bar]=1`, `foo[]=2`},
			err:  `expected an array: foo: {"bar":1}`,
		},
		{
			name: "expected array deep",
			args: []string{`foo[bar][baz][qux]=1`, `foo[bar][baz][]=2`},
			err:  `expected an array: foo.bar.baz: {"qux":1}`,
		},
		{
			name: "read from file",
			args: []string{`foo=@testdata/file.txt`},
			expected: `{"foo":"a\nb\nc\nd\ne"}
`,
		},
		{
			name: "base64 of file",
			args: []string{`foo=%testdata/file.txt`},
			expected: `{"foo":"YQpiCmMKZAplCg=="}
`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var out strings.Builder
			opts := []Option{
				Args(tc.args),
				Output(&out),
			}
			if tc.array {
				opts = append(opts, Array())
			}
			if tc.pretty {
				opts = append(opts, Pretty())
			}
			err := New(opts...).Run()
			if tc.err == "" {
				if err != nil {
					t.Errorf("error should be nil but got: %s", err)
				} else if got := out.String(); got != tc.expected {
					t.Errorf("output should be %q but got %q", tc.expected, got)
				}
			} else {
				if err == nil {
					t.Errorf("error should be nil but got: %s", err)
				} else if got := out.String(); got != tc.expected {
					t.Errorf("output should be %q but got %q", tc.expected, got)
				} else if got := err.Error(); got != tc.err {
					t.Errorf("error should be %q but got %q", tc.err, got)
				}
			}
		})
	}
}
