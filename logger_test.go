package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/xh3b4sd/logger/meta"
	"github.com/xh3b4sd/tracer"
)

// go test -update
var update = flag.Bool("update", false, "update .golden files")

func Test_Logger_Log(t *testing.T) {
	testCases := []struct {
		ctx context.Context
		kvs []string
		cal func() string
	}{
		// Case 000 emits a debug log. The default filter allows info logs. There
		// is no log line in the golden file.
		{
			ctx: context.Background(),
			kvs: []string{
				"level", "debug",
				"message", "foo",
			},
		},
		// Case 001 emits an info log. The default filter allows info logs. There
		// is a log line in the golden file.
		{
			ctx: context.Background(),
			kvs: []string{
				"level", "info",
				"message", "foo bar baz",
			},
		},
		// Case 002 emits a warning log. The default filter allows info logs. There
		// is a log line in the golden file. The reserved key code is included
		// together with its value.
		{
			ctx: context.Background(),
			kvs: []string{
				"level", "warning",
				"message", "foo",
				"code", "red",
			},
		},
		// Case 003 emits an error log. The default filter allows info logs. There
		// is a log line in the golden file.
		{
			ctx: context.Background(),
			kvs: []string{
				"level", "error",
				"message", "foo bar baz",
			},
			cal: EmptyCaller,
		},
		// Case 004 emits an uneven amount of key-value pairs. Nothing but an error
		// message will be logged in the golden file.
		{
			ctx: context.Background(),
			kvs: []string{
				"level", "error",
				"message",
			},
		},
		// Case 005 emits an error log. The log line is annotated with additional
		// meta information. There is a log line in the golden file.
		{
			ctx: func() context.Context {
				ctx := context.Background()

				ctx = meta.Add(ctx, "key", "val")
				ctx = meta.Add(ctx, "mar", "noo")
				ctx = meta.Add(ctx, "foo", "bar")
				ctx = meta.Add(ctx, "zoo", "fre")
				ctx = meta.Add(ctx, "tri", "fan")

				return ctx
			}(),
			kvs: []string{
				"level", "error",
				"message", "foo",
			},
		},
		// Case 006 emits a warning log. The default filter allows info logs. There
		// is a log line in the golden file containing a stack trace of the provided
		// error. Empty values are ignored.
		{
			ctx: context.Background(),
			kvs: []string{
				"level", "warning",
				"message", "foo",
				"light", "darkness",
				"zoo", "",
				"stack", tracer.Stack(tracer.Mask(&tracer.Error{Kind: "testError"})),
			},
		},
		// Case 007 is like 6 but without caller.
		{
			ctx: context.Background(),
			kvs: []string{
				"level", "warning",
				"message", "foo",
				"code", "",
				"stack", tracer.Stack(tracer.Mask(tracer.Mask(&tracer.Error{Kind: "testError"}))),
			},
			cal: EmptyCaller,
		},
		// Case 008 emits a warning log like produced by a *tracer.Error. A bug in
		// the default formatter once happened where invalid JSON was produced due
		// to a missing comma when consecutive fields were empty and then filled.
		// Below code is empty and right after desc is filled.
		{
			ctx: context.Background(),
			kvs: []string{
				"level", "warning",
				"message", "foo",
				"code", "",
				"desc", "whatever",
				"kind", "notfoundError",
			},
		},
		// Case 009 same as above but with an additional empty custom item after
		// desc.
		{
			ctx: context.Background(),
			kvs: []string{
				"level", "warning",
				"message", "foo",
				"code", "",
				"desc", "whatever",
				"engl", "",
				"kind", "notfoundError",
			},
		},
		// Case 010 same as above but with an additional empty supported item after
		// desc.
		{
			ctx: context.Background(),
			kvs: []string{
				"level", "warning",
				"message", "foo",
				"code", "",
				"desc", "whatever",
				"docs", "",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var err error

			w := &bytes.Buffer{}

			var l Interface
			{
				l = New(Config{
					Caller: tc.cal,
					Timer: func() string {
						return "time" // static for the golden files
					},
					Writer: w,
				})
			}

			l.Log(tc.ctx, tc.kvs...)

			// Compute the actually logged output and change the given paths of
			// the caller to avoid prefixes like "/Users/username/go/src/" so
			// this test can be executed on different machines. See the golden
			// files in the testdata folder for specific examples.
			var actual []byte
			if w.Len() > 0 {
				p, err := os.Getwd()
				if err != nil {
					t.Fatal(err)
				}
				b := bytes.ReplaceAll(w.Bytes(), []byte(p), []byte("--REPLACED--"))

				buf := &bytes.Buffer{}
				err = json.Indent(buf, b, "", "  ")
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
				err := os.WriteFile(p, actual, 0644) // nolint:gosec
				if err != nil {
					t.Fatal(err)
				}
			}

			expected, err := os.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(actual, expected) {
				t.Fatalf("\n\n%s\n", cmp.Diff(expected, actual))
			}
		})
	}
}

func fileName(i int) string {
	return "case-" + fmt.Sprintf("%03d", i) + ".golden.json"
}

func isJSONError(err error) bool {
	_, ok := err.(*json.SyntaxError)
	return ok
}
