package sequence

import (
	"fmt"
	"iter"
	"strings"
)

func Enum[E any](s iter.Seq[E]) iter.Seq2[int, E] {
	return func(yield func(int, E) bool) {
		var i int
		for v := range s {
			if !yield(i, v) {
				break
			}
			i++
		}
	}
}

func ValueAt[E any](s iter.Seq[E], index int) (E, bool) {
	for i, v := range Enum(s) {
		if i == index {
			return v, true
		}
	}
	var zeroValue E
	return zeroValue, false
}

func String[E any](s iter.Seq[E]) string {
	return Format(s, " ")
}

func Format[E any](s iter.Seq[E], sep string) string {
	var buf strings.Builder
	buf.WriteRune('[')
	for i, v := range Enum(s) {
		if i > 0 {
			buf.WriteString(sep)
		}
		if _, err := fmt.Fprint(&buf, v); err != nil {
			panic(fmt.Errorf("sequence.Format: write value: %w", err))
		}
	}
	buf.WriteRune(']')
	return buf.String()
}
