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
	"github.com/xh3b4sd/logger/meta"
	"github.com/xh3b4sd/tracer"
)

var update = flag.Bool("update", false, "update .golden files")

// Test_Logger_Log tests the logger behaviour based on its output. The tests use
// golden file references. In case the golden files change something is broken.
// In case intentional changes get introduced the golden files have to be
// updated. In case the golden files have to be adjusted, simply provide the
// -update flag when running the tests.
//
//     go test . -run Test_Logger_Log -update
//
func Test_Logger_Log(t *testing.T) {
	testCases := []struct {
		ctx context.Context
		kvs []string
	}{
		// Case 0 emits a debug log. The default filter allows info logs. There
		// is no log line in the golden file.
		{
			ctx: context.Background(),
			kvs: []string{
				"level", "debug",
				"message", "foo",
			},
		},
		// Case 1 emits an info log. The default filter allows info logs. There
		// is a log line in the golden file.
		{
			ctx: context.Background(),
			kvs: []string{
				"level", "info",
				"message", "foo",
			},
		},
		// Case 2 emits a warning log. The default filter allows info logs. There
		// is a log line in the golden file.
		{
			ctx: context.Background(),
			kvs: []string{
				"level", "warning",
				"message", "foo",
			},
		},
		// Case 3 emits an error log. The default filter allows info logs. There
		// is a log line in the golden file.
		{
			ctx: context.Background(),
			kvs: []string{
				"level", "error",
				"message", "foo",
			},
		},
		// Case 4 emits an uneven amount of key-value pairs. Nothing but an
		// error message will be logged in the golden file.
		{
			ctx: context.Background(),
			kvs: []string{
				"level", "error",
				"message",
			},
		},
		// Case 5 emits an error log. The log line is annotated with additional
		// meta information. There is a log line in the golden file.
		{
			ctx: func() context.Context {
				ctx := context.Background()

				ctx = meta.Add(ctx, "foo", "bar")
				ctx = meta.Add(ctx, "key", "val")

				return ctx
			}(),
			kvs: []string{
				"level", "error",
				"message", "foo",
			},
		},
		// Case 6 emits a warning log. The default filter allows info logs.
		// There is a log line in the golden file containing a stack trace of
		// the provided error.
		{
			ctx: context.Background(),
			kvs: []string{
				"level", "warning",
				"message", "foo",
				"stack", tracer.JSON(&tracer.Error{Kind: "testError"}),
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
				if isJSONError(err) {
					actual = w.Bytes()
				} else if err != nil {
					t.Fatal(err)
				} else {
					actual = buf.Bytes()
				}
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

func isJSONError(err error) bool {
	_, ok := err.(*json.SyntaxError)
	return ok
}
