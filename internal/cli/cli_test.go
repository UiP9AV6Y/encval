package cli

import (
	"testing"

	"github.com/alecthomas/kong"
	"gotest.tools/v3/assert"

	"github.com/UiP9AV6Y/encval/internal/cli/cmd"
)

func TestCliParse(t *testing.T) {
	out, err := New()
	assert.NilError(t, err)
	assert.Assert(t, out != nil)

	type globalOptionsAssertion func(*cmd.GlobalOptions)
	type nodeAssertion func(*kong.Node)

	type testCase struct {
		name        string
		input       []string
		assertOpts  globalOptionsAssertion
		assertNode  nodeAssertion
		assertError bool
	}

	run := func(t *testing.T, tc testCase) {
		got, err := out.parse(tc.input)

		if tc.assertError {
			assert.Assert(t, err != nil)
		} else {
			assert.NilError(t, err)
			assert.Assert(t, got != nil)

			if tc.assertNode != nil {
				node := got.Selected()
				assert.Assert(t, node != nil)
				tc.assertNode(node)
			}

			if tc.assertOpts != nil {
				opts := out.Context()
				assert.Assert(t, opts != nil)
				tc.assertOpts(opts)
			}
		}
	}
	testCases := []testCase{
		{
			name:        "empty",
			input:       []string{},
			assertError: true,
		},
		{
			name:  "version",
			input: []string{"version"},
			assertNode: func(n *kong.Node) {
				assert.Equal(t, n.Name, "version")
			},
		},
		{
			name:  "verbose",
			input: []string{"-vvv", "version"},
			assertOpts: func(o *cmd.GlobalOptions) {
				assert.Equal(t, o.Verbose, 3)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run(t, tc)
		})
	}
}
