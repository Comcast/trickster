/*
 * Copyright 2018 Comcast Cable Communications Management, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ruler

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/Comcast/trickster/internal/proxy/urls"
)

const testURLRaw = "https://example.com:8480/path1/path2?param1=value&param2=value&param1=value2"

var testURL, _ = url.Parse(testURLRaw)

var testRL0 = rewriteList{
	[]string{"header", "set", "Cache-Control", "max-age=60"},
	[]string{"header", "append", "Cache-Control", "max-age=300"},
	[]string{"header", "append", "Cache-Control", "private"},
	[]string{"header", "append", "Cache-Control", "private"},
	[]string{"header", "set", "Test-Header", "Trickster"},
	[]string{"header", "replace", "Cache-Control", "300", "60"},
	[]string{"header", "delete", "Test-Header"},
	[]string{"header", "delete", "Cache-Control", "private"},
	[]string{"header", "append", "Cache-Control", "smax-age=30"},
	[]string{"param", "set", "param1", "foo"},
	[]string{"param", "append", "param1", "value2"},
	[]string{"param", "set", "param2", "${trickster}"},
	[]string{"param", "replace", "param1", "foo", "bar"},
	[]string{"param", "replace", "paramX", "foo", "bar"},
	[]string{"param", "delete", "param2"},
	[]string{"param", "delete", "param1", "value2"},
	[]string{"param", "append", "param1", "too"},
	[]string{"param", "append", "param1", "too"},
	[]string{"param", "append", "param3", "trickster"},
}

var testRL1 = rewriteList{
	[]string{"path", "set", "my/path/is/here"},
	[]string{"path", "set", "was", "2"},
	[]string{"path", "replace", "he", "the", "3"},
	[]string{"path", "replace", "the", "the"}, // test depth -1
}

var testRL2 = rewriteList{
	[]string{"params", "set", "param1=foo&param2=trickster&param3=foo&param1=too"},
	[]string{"params", "replace", "foo", "bar", "1"},
}

var testRL3 = rewriteList{
	[]string{"method", "set", "POST"},
	[]string{"host", "set", "example.com:9090"},
	[]string{"host", "replace", "example.com", "tricksterproxy.io"},
	[]string{"port", "delete"},
	[]string{"port", "set", "8000"},
	[]string{"port", "replace", "000", "480"},
}

var testRI0 = rewriteInstructions{
	&rwiKeyBasedSetter{key: "Cache-Control", value: "max-age=60"},
	&rwiKeyBasedAppender{key: "Cache-Control", value: "max-age=300"},
	&rwiKeyBasedAppender{key: "Cache-Control", value: "private"},
	&rwiKeyBasedAppender{key: "Cache-Control", value: "private"},
	&rwiKeyBasedSetter{key: "Test-Header", value: "Trickster"},
	&rwiKeyBasedReplacer{key: "Cache-Control", search: "300", replacement: "60"},
	&rwiKeyBasedDeleter{key: "Test-Header"},
	&rwiKeyBasedDeleter{key: "Cache-Control", value: "private"},
	&rwiKeyBasedAppender{key: "Cache-Control", value: "smax-age=30"},
	&rwiKeyBasedSetter{key: "param1", value: "foo"},
	&rwiKeyBasedAppender{key: "param1", value: "value2"},
	&rwiKeyBasedSetter{key: "param2", value: "${trickster}", hasTokens: true},
	&rwiKeyBasedReplacer{key: "param1", search: "foo", replacement: "bar"},
	&rwiKeyBasedReplacer{key: "paramX", search: "foo", replacement: "bar"},
	&rwiKeyBasedDeleter{key: "param2"},
	&rwiKeyBasedDeleter{key: "param1", value: "value2"},
	&rwiKeyBasedAppender{key: "param1", value: "too"},
	&rwiKeyBasedAppender{key: "param1", value: "too"},
	&rwiKeyBasedAppender{key: "param3", value: "trickster"},
}

var testRI1 = rewriteInstructions{
	&rwiPathSetter{value: "my/path/is/here", depth: -1},
	&rwiPathSetter{value: "was", depth: 2},
	&rwiPathReplacer{search: "he", replacement: "the", depth: 3},
	&rwiPathReplacer{search: "the", replacement: "the", depth: -1},
}

var testRI2 = rewriteInstructions{
	&rwiBasicSetter{value: "param1=foo&param2=trickster&param3=foo&param1=too"},
	&rwiBasicReplacer{search: "foo", replacement: "bar", depth: 1},
}

var testRI3 = rewriteInstructions{
	&rwiBasicSetter{value: "POST"},
	&rwiBasicSetter{value: "example.com:9090"},
	&rwiBasicReplacer{search: "example.com", replacement: "tricksterproxy.io", depth: -1},
	&rwiPortDeleter{},
	&rwiBasicSetter{value: "8000"},
	&rwiBasicReplacer{search: "000", replacement: "480", depth: -1},
}

func TestParseRewriteList(t *testing.T) {

	var tests = []struct {
		rl          rewriteList
		expected    rewriteInstructions
		expectedErr error
	}{
		// run 0: key-based instructions
		{
			rl:          testRL0,
			expected:    testRI0,
			expectedErr: nil,
		},

		// run 1: path-based instructions
		{
			rl:          testRL1,
			expected:    testRI1,
			expectedErr: nil,
		},

		// run 2: basic instructions - params coverage
		{
			rl:          testRL2,
			expected:    testRI2,
			expectedErr: nil,
		},

		// run 3: basic instructions - method, host, port coverage
		{
			rl:          testRL3,
			expected:    testRI3,
			expectedErr: nil,
		},

		// runs 4-: error cases
		{
			// 4 - key-based set error case A
			rl: rewriteList{
				[]string{"header", "set"},
			},
			expectedErr: errBadParams,
		},

		{
			// 5 - key-based set error case B
			rl: rewriteList{
				[]string{"header", "TESTING"},
			},
			expectedErr: errBadParams,
		},

		{
			// 6 - key-based replace error case A
			rl: rewriteList{
				[]string{"header", "replace"},
			},
			expectedErr: errBadParams,
		},
		{
			// 7 - key-based delete error case A
			rl: rewriteList{
				[]string{"header", "delete"},
			},
			expectedErr: errBadParams,
		},
		{
			// 8 - path setter error case A
			rl: rewriteList{
				[]string{"path", "set"},
			},
			expectedErr: errBadParams,
		},
		{
			// 9 - path replacer error case A
			rl: rewriteList{
				[]string{"path", "replace"},
			},
			expectedErr: errBadParams,
		},
		{
			// 10 - basic setter error case A
			rl: rewriteList{
				[]string{"params", "set"},
			},
			expectedErr: errBadParams,
		},
		{
			// 11 - basic setter error case A
			rl: rewriteList{
				[]string{"params", "replace"},
			},
			expectedErr: errBadParams,
		},
		{
			// 12 - basic replacer error case B
			rl: rewriteList{
				[]string{"params", "replace", "foo", "bar", "not-an-integer"},
			},
			expectedErr: errBadDepthParse,
		},
		{
			// 13 - path replacer error case B
			rl: rewriteList{
				[]string{"path", "replace", "foo", "bar", "not-an-integer"},
			},
			expectedErr: errBadDepthParse,
		},
		{
			// 14 - path setter error case B
			rl: rewriteList{
				[]string{"path", "set", "foo", "not-an-integer"},
			},
			expectedErr: errBadDepthParse,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {

			got, err := parseRewriteList(test.rl)
			if err != test.expectedErr {
				t.Errorf("expected error %s got %s", test.expectedErr, err)
			}

			if len(got) != len(test.expected) {
				t.Errorf("expected %d got %d", len(test.expected), len(got))
			}

			if got.String() != test.expected.String() {
				t.Errorf("\ngot      %s\nexpected %s", got.String(), test.expected.String())
			}

		})
	}
}

func TestDictFuncsNilRequest(t *testing.T) {

	f := dicts["header"]
	d := f(nil)
	if d != nil {
		t.Error(errors.New("expected nil value"))
	}

	f = dicts["param"]
	d = f(nil)
	if d != nil {
		t.Error(errors.New("expected nil value"))
	}

}

func TestExecuteRewriteInstructions(t *testing.T) {

	var exh0 = http.Header{"Cache-Control": []string{"max-age=60, smax-age=30"}}
	eu0, _ := url.Parse("https://example.com:8480/path1/path2?param1=bar&param1=too&param3=trickster")
	ri0, _ := parseRewriteList(testRL0)

	eu1, _ := url.Parse("https://example.com:8480/my/path/was/there?param1=value&param2=value&param1=value2")
	ri1, _ := parseRewriteList(testRL1)

	eu2, _ := url.Parse("https://example.com:8480/path1/path2?param1=bar&param2=trickster&param3=foo&param1=too")
	ri2, _ := parseRewriteList(testRL2)

	eu3, _ := url.Parse("https://tricksterproxy.io:8480/path1/path2?param1=value&param2=value&param1=value2")
	ri3, _ := parseRewriteList(testRL3)

	var tests = []struct {
		in       *http.Request
		ri       rewriteInstructions
		expected *http.Request
	}{
		// run 0: key-based instructions
		{
			in:       &http.Request{Method: "GET", URL: urls.Clone(testURL), Header: make(http.Header)},
			ri:       ri0,
			expected: &http.Request{Method: "GET", URL: eu0, Header: exh0},
		},
		// run 1: path-based instructions
		{
			in:       &http.Request{Method: "GET", URL: urls.Clone(testURL), Header: make(http.Header)},
			ri:       ri1,
			expected: &http.Request{Method: "GET", URL: eu1, Header: make(http.Header)},
		},
		// run 2: params (not key-based) instructions
		{
			in:       &http.Request{Method: "GET", URL: urls.Clone(testURL), Header: make(http.Header)},
			ri:       ri2,
			expected: &http.Request{Method: "GET", URL: eu2, Header: make(http.Header)},
		},
		// run 3: host/port-based instructions
		{
			in:       &http.Request{Method: "GET", URL: urls.Clone(testURL), Header: make(http.Header)},
			ri:       ri3,
			expected: &http.Request{Method: "POST", URL: eu3, Header: make(http.Header)},
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			test.ri.Execute(test.in)
			if !reqLazyEqual(test.in, test.expected) {
				t.Errorf("\ngot:\n%s\n\nexpected:\n%s", reqString(test.in), reqString(test.expected))
			}
		})
	}
}

func reqLazyEqual(r1, r2 *http.Request) bool {

	if r1 == nil && r2 == nil {
		return true
	}
	if r1 == nil || r2 == nil {
		return false
	}

	return reqString(r1) == reqString(r2)

}

func TestHasTokens(t *testing.T) {

	ris := rewriteInstructions{
		&rwiPathSetter{},
		&rwiPathReplacer{},
		&rwiKeyBasedDeleter{},
		&rwiKeyBasedReplacer{},
		&rwiKeyBasedSetter{},
		&rwiBasicSetter{},
		&rwiBasicReplacer{},
		&rwiPortDeleter{},
		&rwiKeyBasedAppender{},
	}

	for _, ri := range ris {
		if ri.HasTokens() {
			t.Error(errors.New("expected false got true"))
		}
	}

}

func TestNilRequestGetters(t *testing.T) {
	for _, f := range scalarGets {
		v := f(nil)
		if v != "" {
			t.Errorf("expected empty string got %s", v)
		}
	}
}

func TestMiscRequestGetters(t *testing.T) {

	r := &http.Request{Method: "GET", URL: testURL}
	fm := scalarGets["method"]
	fh := scalarGets["hostname"]

	v := fm(r)
	if v != "GET" {
		t.Errorf("expected %s got %s", "GET", v)
	}

	v = fh(r)
	if v != "example.com" {
		t.Errorf("expected %s got %s", "example.com", v)
	}

}

func TestMiscRequestSetters(t *testing.T) {

	r := &http.Request{Method: "GET", URL: testURL}
	fp := scalarSets["port"]
	fh := scalarSets["hostname"]

	fp(nil, "")

	fp(r, "8480")
	fh(r, "tricksterproxy.io")

	if r.URL.Host != "tricksterproxy.io:8480" {
		t.Errorf("expected %s got %s", "tricksterproxy.io:8480", r.URL.Host)
	}

	var s rewriteInstruction

	s = &rwiKeyBasedSetter{}
	err := s.Parse([]string{"foo", "foo", "foo", "foo"})
	if err != errBadParams {
		t.Error(errors.New("expected bad params error"))
	}

	s = &rwiKeyBasedReplacer{}
	err = s.Parse([]string{"foo", "foo", "foo", "foo", "foo"})
	if err != errBadParams {
		t.Error(errors.New("expected bad params error"))
	}

	s = &rwiKeyBasedDeleter{}
	err = s.Parse([]string{"foo", "foo", "foo", "foo"})
	if err != errBadParams {
		t.Error(errors.New("expected bad params error"))
	}

	s = &rwiBasicReplacer{}
	err = s.Parse([]string{"foo", "foo", "foo", "foo"})
	if err != errBadParams {
		t.Error(errors.New("expected bad params error"))
	}

	s = &rwiBasicSetter{}
	err = s.Parse([]string{"foo", "foo", "foo"})
	if err != errBadParams {
		t.Error(errors.New("expected bad params error"))
	}

	s = &rwiKeyBasedAppender{}
	err = s.Parse([]string{"foo", "foo", "foo"})
	if err != errBadParams {
		t.Error(errors.New("expected bad params error"))
	}

	err = s.Parse([]string{"foo", "foo", "foo", "foo"})
	if err != errBadParams {
		t.Error(errors.New("expected bad params error"))
	}
}

func reqString(r *http.Request) string {

	if r == nil || r.URL == nil {
		return ""
	}

	sb := strings.Builder{}

	var q string
	if r.URL.RawQuery != "" {
		q = "?" + r.URL.RawQuery
	}

	sb.WriteString(r.Method + " " + r.URL.Path + q + " " + r.Proto + "\n")
	sb.WriteString("Host: " + r.URL.Host + "\n")

	if r.Header != nil {
		for k := range r.Header {
			sb.WriteString(k + ": " + r.Header.Get(k) + "\n")
		}
	}

	sb.WriteString("\n")

	return sb.String()
}
