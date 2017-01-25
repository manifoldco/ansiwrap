// Package ansiwrap implements ANSI terminal escape code aware text wrapping.
package ansiwrap

import (
	"math"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Wrap wraps the given string str so that each line is at most width printable
// characters wide.
//
// Newlines are inserted at unicode space characters.
// If the distance between two space characters is larger than width, no
// newline is inserted.
//
// str is interpreted as UTF-8.
//
// Both non-printable unicode characters, and ANSI terminal escape sequences
// are ignored for counting purposes.
//
// Wrap selects the best function between Balanced and Greedy, based on the
// length of str and width.
//
// The Balanced algorithm wraps lineswith minimal raggedness using the
// 'Divide & Conquer' algorithm described at http://xxyxyz.org/line-breaking/
//
// The Greedy algorithm greedily fills a line as close to width as possible
// before continuing to the next line.
func Wrap(str string, width int) string {
	return WrapIndent(str, width, 0, 0)
}

// WrapIndent wraps lines following the same rules as Wrap, in addition
// to exposing optional indent values firstIndent and restIndent.
//
// firstIndent will indent the first line by that number of spaces.
// restIndent will indent all remaining lines by that number of spaces.
// Indent values are taken into account when calculating wrapping.
func WrapIndent(str string, width, firstIndent, restIndent int) string {
	rc := RuneCount(str)
	if rc+firstIndent < width*2 {
		return GreedyIndent(str, width, firstIndent, restIndent)
	}
	return BalancedIndent(str, width, firstIndent, restIndent)
}

// Balanced is equivalent to Wrap, but always using the Balanced algorithm.
func Balanced(str string, width int) string {
	return BalancedIndent(str, width, 0, 0)
}

// BalancedIndent is equivalent to WrapIndent, but always using the Balanced
// algorithm.
func BalancedIndent(str string, width, firstIndent, restIndent int) string {
	words := strings.Fields(str)
	count := len(words)

	offsets := make([]int, count+1)
	minima := make([]int, count+1)
	breaks := make([]int, count+1)
	for i, word := range words {
		offsets[i+1] = offsets[i] + RuneCount(word)
		minima[i+1] = math.MaxInt32
	}

	cost := func(i, j int) int {
		indent := restIndent
		if i == 0 {
			indent = firstIndent
		}

		adjWidth := width - indent

		w := offsets[j] - offsets[i] + j - i - 1
		if w > adjWidth {
			return math.MaxInt32 - 1
		}

		return minima[i] + (adjWidth-w)*(adjWidth-w)
	}

	search := func(i0, j0, i1, j1 int) {
		stack := [][4]int{{i0, j0, i1, j1}}
		for len(stack) > 0 {
			l := stack[len(stack)-1]
			i0, j0, i1, j1 := l[0], l[1], l[2], l[3]
			stack = stack[:len(stack)-1]
			if j0 < j1 {
				j := (j0 + j1) / 2
				for i := i0; i < i1; i++ {
					c := cost(i, j)
					if c <= minima[j] {
						minima[j] = c
						breaks[j] = i
					}
				}
				stack = append(stack, [4]int{breaks[j], j + 1, i1, j1})
				stack = append(stack, [4]int{i0, j0, breaks[j] + 1, j})
			}
		}
	}

	n := count + 1
	var i uint
	offset := 0

OuterLoop:
	for {
		r := min(n, 1<<(i+1))
		edge := 1<<i + offset
		search(offset, edge, edge, r+offset)
		x := minima[r-1+offset]
		for j := 1 << i; j < r-1; j++ {
			y := cost(j+offset, r-1+offset)
			if y <= x {
				n -= j
				i = 0
				offset += j
				continue OuterLoop
			}
		}
		if r == n {
			break
		}
		i++
	}

	var lines []string
	var y int
	for j := count; j > 0; j = y {
		y = breaks[j]

		indent := restIndent
		if y <= 0 {
			indent = firstIndent
		}

		lines = append([]string{strings.Repeat(" ", indent) + strings.Join(words[y:j], " ")}, lines...)
	}
	return strings.Join(lines, "\n")
}

// Greedy is equivalent to Wrap, but always using the Greedy algorithm.
func Greedy(str string, width int) string {
	return GreedyIndent(str, width, 0, 0)
}

// GreedyIndent is equivalent to WrapIndent, but always using the Greedy
// algorithm.
func GreedyIndent(str string, width, firstIndent, restIndent int) string {
	words := strings.Fields(str)
	var lines []string
	for i := 0; i < len(words); {
		count := 0

		indent := restIndent
		if len(lines) == 0 {
			indent = firstIndent
		}

		line := strings.Repeat(" ", indent)

		for i < len(words) {
			w := words[i]
			rc := RuneCount(w)
			if count != 0 { // not the first word on line; join with a space.
				w = " " + w
				rc++
			}

			// Break if we're past the line width. The first word on a line can
			// be longer; how else would we fit it?
			if indent+count+rc > width && count != 0 {
				break
			}

			line += w
			count += rc
			i++
		}
		lines = append(lines, line)

	}

	return strings.Join(lines, "\n")
}

// RuneCount counts the number of printable runes in str.
// In addition to ignoring non-printable unicode characters, it ignores all
// ANSI escape sequences.
func RuneCount(str string) int {
	l := 0
	b := []byte(str)
	inSequence := false
	for len(b) > 0 {
		if b[0] == '\033' {
			inSequence = true
			b = b[1:]
			continue
		}

		r, rl := utf8.DecodeRune(b)
		b = b[rl:]

		if inSequence {
			if r == 'm' {
				inSequence = false
			}

			continue
		}

		if !unicode.IsPrint(r) {
			continue
		}

		l++
	}

	return l
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
