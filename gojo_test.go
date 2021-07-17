package gojo_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/itchyny/gojo"
)

func TestGojoMap(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected string
		err      string
	}{
		{
			name:     "default",
			args:     []string{"foo=bar"},
			expected: `{"foo":"bar"}`,
		},
		{
			name:     "multiple",
			args:     []string{"foo=bar", "bar=false", "baz=qux", "\n=", `\n=null`},
			expected: `{"foo":"bar","bar":false,"baz":"qux","\n":"","\\n":null}`,
		},
		{
			name:     "numbers",
			args:     []string{"a=123", "c=3.14", "d=3e10", "b=-128", "e=[1,2,3]", "f=0xffdc", "g=0XFF", "h=037"},
			expected: `{"a":123,"c":3.14,"d":30000000000,"b":-128,"e":[1,2,3],"f":65500,"g":255,"h":37}`,
		},
		{
			name:     "nested object",
			args:     []string{`foo={"bar":{"baz":"qux","quux":[{"y":1,"x":2}]}}`, `bar=[{"y":1,"x":2},{"x":3,"y":4}]`},
			expected: `{"foo":{"bar":{"baz":"qux","quux":[{"y":1,"x":2}]}},"bar":[{"y":1,"x":2},{"x":3,"y":4}]}`,
		},
		{
			name:     "nested keys",
			args:     []string{"foo[]=1", "foo[]=bar", "bar[baz]=10", "qux][=20", "a[b]c=d"},
			expected: `{"foo":[1,"bar"],"bar":{"baz":10},"qux][":20,"a[b]c":"d"}`,
		},
		{
			name:     "deep keys",
			args:     []string{"foo[]=1", "foo[]=bar", "bar[baz]=10", "qux][=20", "a[b]c=d"},
			expected: `{"foo":[1,"bar"],"bar":{"baz":10},"qux][":20,"a[b]c":"d"}`,
		},
		{
			name:     "merge to json object",
			args:     []string{`a={"b":{"c":10,"d":[20]}}`, "a[b][d][]=30", "a[b][e]=40"},
			expected: `{"a":{"b":{"c":10,"d":[20,30],"e":40}}}`,
		},
		{
			name:     "merge to json array",
			args:     []string{`a=[{"b":10}]`, "a[][c][]=20"},
			expected: `{"a":[{"b":10},{"c":[20]}]}`,
		},
		{
			name:     "read from file",
			args:     []string{"foo=@testdata/file.txt"},
			expected: `{"foo":"a\nb\nc\nd\ne"}`,
		},
		{
			name:     "json object of file",
			args:     []string{"foo=:testdata/file1.json"},
			expected: `{"foo":{"x":1,"z":2,"y":3}}`,
		},
		{
			name:     "json array of file",
			args:     []string{"foo=:testdata/file2.json"},
			expected: `{"foo":[{"y":1,"x":2},{"x":3,"y":4},{"y":4,"z":5}]}`,
		},
		{
			name:     "json number of file",
			args:     []string{"foo=:testdata/file3.json"},
			expected: `{"foo":42}`,
		},
		{
			name:     "json string of file",
			args:     []string{"foo=:testdata/file4.json"},
			expected: `{"foo":"hello, world"}`,
		},
		{
			name:     "base64 of file",
			args:     []string{"foo=%testdata/file.txt"},
			expected: `{"foo":"YQpiCmMKZAplCg=="}`,
		},
		{
			name: "parse error",
			args: []string{"foo"},
			err:  `failed to parse: "foo"`,
		},
		{
			name: "expected object",
			args: []string{"foo[]=1", "foo[bar]=2"},
			err:  "expected an object: foo: [1]",
		},
		{
			name: "expected object deep",
			args: []string{"foo[bar][]=1", "foo[bar][baz]=2"},
			err:  "expected an object: foo.bar: [1]",
		},
		{
			name: "expected array",
			args: []string{"foo[bar]=1", "foo[]=2"},
			err:  `expected an array: foo: {"bar":1}`,
		},
		{
			name: "expected array deep",
			args: []string{"foo[bar][baz][qux]=1", "foo[bar][baz][]=2"},
			err:  `expected an array: foo.bar.baz: {"qux":1}`,
		},
		{
			name: "json error",
			args: []string{"foo=:testdata/file5.json"},
			err:  "invalid character 'x' after top-level value",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := gojo.Map(tc.args)
			if tc.err == "" {
				if err == nil {
					if got := diff(tc.expected, got); got != "" {
						t.Error(got)
					}
				} else {
					t.Errorf("error should be nil but got: %s", err)
				}
			} else {
				if err == nil || !strings.Contains(err.Error(), tc.err) {
					t.Errorf("should return error: %q, got %v", tc.err, err)
				}
			}
		})
	}
}

func TestGojoArray(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected string
		err      error
	}{
		{
			name:     "array",
			args:     []string{"foo", "false", "baz", "null", "3.14", "0xff"},
			expected: `["foo",false,"baz",null,3.14,255]`,
		},
		{
			name:     "nested array",
			args:     []string{"[[[[1]]]]"},
			expected: "[[[[[1]]]]]",
		},
		{
			name:     "read from file",
			args:     []string{"@testdata/file.txt"},
			expected: `["a\nb\nc\nd\ne"]`,
		},
		{
			name:     "json object of file",
			args:     []string{":testdata/file1.json"},
			expected: `[{"x":1,"z":2,"y":3}]`,
		},
		{
			name:     "json array of file",
			args:     []string{":testdata/file2.json"},
			expected: `[[{"y":1,"x":2},{"x":3,"y":4},{"y":4,"z":5}]]`,
		},
		{
			name:     "json number of file",
			args:     []string{":testdata/file3.json"},
			expected: "[42]",
		},
		{
			name:     "json string of file",
			args:     []string{":testdata/file4.json"},
			expected: `["hello, world"]`,
		},
		{
			name:     "base64 of file",
			args:     []string{"%testdata/file.txt"},
			expected: `["YQpiCmMKZAplCg=="]`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := gojo.Array(tc.args)
			if got := diff(tc.expected, got); got != "" {
				t.Error(got)
			}
		})
	}
}

func diff(expected string, got interface{}) string {
	bs, _ := json.Marshal(got)
	if string(bs) == expected {
		return ""
	}
	return fmt.Sprintf("diff:\nexpected: %s\n     got: %s", expected, bs)
}
