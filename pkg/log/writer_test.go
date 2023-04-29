package log

import (
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestStreamWriter(t *testing.T) {
	var buf strings.Builder
	out := NewStreamWriter(&buf)

	out.Print("test")
	out.Printf("case %d\n", 1)
	out.Println("end of test")

	got := buf.String()
	want := "testcase 1\nend of test\n"

	assert.Assert(t, out.Enabled())
	assert.Equal(t, got, want)
}

func TestDisabledStreamWriter(t *testing.T) {
	out := NewDisabledStreamWriter()

	out.Print("test")
	out.Printf("case %d\n", 1)
	out.Println("end of test")

	assert.Assert(t, !out.Enabled())
}
