package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var update = flag.Bool("update", false, "update .golden files")

// Test_MicroLogger tests MicroLogger output.
//
// It uses golden file as reference and when changes to template are
// intentional, they can be updated by providing -update flag for go test.
//
//	go test . -run Test_MicroLogger -update
//
func Test_MicroLogger(t *testing.T) {
	testCases := []struct {
		ctx context.Context
		kvs []string
	}{
		// Case 0 emits a debug log. The default filter allows info logs. There
		// is no log line in the golden file.
		{
			ctx: context.Background(),
			kvs: []string{
				"lvl", "deb",
				"mes", "foo",
			},
		},
		// Case 1 emits an info log. The default filter allows info logs. There
		// is a log line in the golden file.
		{
			ctx: context.Background(),
			kvs: []string{
				"lvl", "inf",
				"mes", "foo",
			},
		},
		// Case 2 emits a warning log. The default filter allows info logs. There
		// is a log line in the golden file.
		{
			ctx: context.Background(),
			kvs: []string{
				"lvl", "war",
				"mes", "foo",
			},
		},
		// Case 3 emits an error log. The default filter allows info logs. There
		// is a log line in the golden file.
		{
			ctx: context.Background(),
			kvs: []string{
				"lvl", "err",
				"mes", "foo",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var err error

			w := &bytes.Buffer{}

			var l Interface
			{
				c := Config{
					Timer: func() string {
						return "time" // static for the golden files
					},
					Writer: w,
				}

				l, err = New(c)
				if err != nil {
					t.Fatal(err)
				}
			}

			l.Log(tc.ctx, tc.kvs...)

			// Compute the actually logged output and change the given paths of
			// the caller to avoid prefixes like "/Users/username/go/src/" so
			// this test can be executed on different machines. See the golden
			// files in the testdata folder for specific examples.
			var actual []byte
			if w.Len() > 0 {
				r := regexp.MustCompile(`("caller"\s*:\s*")\S+(/[^/"]+.go:\d+")`)
				b := r.ReplaceAll(w.Bytes(), []byte("$1--REPLACED--$2"))

				buf := &bytes.Buffer{}
				err := json.Indent(buf, b, "", "\t")
				if err != nil {
					t.Fatal(err)
				}
				actual = buf.Bytes()
			}

			p := filepath.Join("testdata", fileName(i))
			if *update {
				err := ioutil.WriteFile(p, actual, 0644) // nolint:gosec
				if err != nil {
					t.Fatal(err)
				}
			}

			expected, err := ioutil.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(actual, expected) {
				t.Fatalf("\n\n%s\n", cmp.Diff(actual, expected))
			}
		})
	}
}

func fileName(i int) string {
	return "case-" + strconv.Itoa(i) + ".golden"
}
