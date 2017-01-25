package ansiwrap

import (
	"fmt"
	"testing"
)

type testCase struct {
	in, out string
	w       int
	f, r    int
}

var tcs = []testCase{
	{in: "Some text.", out: "Some\ntext.", w: 4},  // Simple wrap
	{in: "Some text.", out: "Some text.", w: 20},  // No wrap needed
	{in: "Some text.", out: "Some\ntext.", w: 7},  // wrap in middle of word
	{in: "Some\ntext.", out: "Some text.", w: 16}, // existing break removed

	{in: "Mayflower", out: "Mayflower", w: 6}, // wrap in middle of single word

	{in: "neat™ test", out: "neat™ test", w: 10}, // unicode counted correctly

	// ANSI escape codes are handled correctly
	{in: "\033[44mColored\033[0m text", out: "\033[44mColored\033[0m text", w: 12},
	{in: "\033[34mColored\033[0m text", out: "\033[34mColored\033[0m\ntext", w: 7},

	// Indent test cases
	{in: "first indent", out: "  first indent", w: 80, f: 2},
	{in: "first wr", out: "  first\nwr", w: 8, f: 2},
	{in: "the rest indent", out: "the rest\n  indent", w: 8, r: 2},
	{in: "the rest indent multi ln", out: "the rest\n  indent\n  multi\n  ln", w: 8, r: 2},
}

func TestBalanced(t *testing.T) {
	for _, tc := range tcs {
		var out string
		if tc.f == 0 && tc.r == 0 {
			out = Balanced(tc.in, tc.w)
		} else {
			out = BalancedIndent(tc.in, tc.w, tc.f, tc.r)
		}

		if tc.out != out {
			t.Errorf("Bad wrapping len %d.\nWanted:\n%s\nGot:\n%s", tc.w, tc.out, out)
		}
	}
}

func TestGreedy(t *testing.T) {
	for _, tc := range tcs {
		var out string
		if tc.f == 0 && tc.r == 0 {
			out = Greedy(tc.in, tc.w)
		} else {
			out = GreedyIndent(tc.in, tc.w, tc.f, tc.r)
		}

		if tc.out != out {
			t.Errorf("Bad wrapping len %d.\nWanted:\n%s\nGot:\n%s", tc.w, tc.out, out)
		}
	}
}

func TestWrap(t *testing.T) {
	tcs := []struct {
		in, out string
		w       int
	}{
		{in: "don't be so greedy", out: "don't\nbe so\ngreedy", w: 8},
		{
			in:  "this test case should be greedy to wrap nicer",
			out: "this test case should be greedy to wrap\nnicer",
			w:   41,
		},
	}

	for _, tc := range tcs {
		out := Wrap(tc.in, tc.w)
		if tc.out != out {
			t.Errorf("Bad wrapping len %d.\nWanted:\n%s\nGot:\n%s", tc.w, tc.out, out)
		}
	}

}

func ExampleWrapIndent_firstIndent() {
	fmt.Println("Example:")
	fmt.Println(WrapIndent("firstIndent can create an indent", 12, 4, 0))
	// Output: Example:
	//     firstIndent can
	// create an indent
}

func ExampleWrapIndent_restIndent() {
	fmt.Println(WrapIndent("restIndent can create a hanging indent", 12, 0, 2))
	// Output: restIndent
	//   can create
	//   a hanging
	//   indent
}

func ExampleBalanced() {
	fmt.Println(Balanced("balanced output is much less ragged than greedy output", 40))
	// Output: balanced output is much less
	// ragged than greedy output
}

func ExampleGreedy() {
	fmt.Println(Greedy("balanced output is much less ragged than greedy output", 40))
	// Output: balanced output is much less ragged than
	// greedy output
}
