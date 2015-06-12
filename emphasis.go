// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
// Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program.  If not, see <http://www.gnu.org/licenses/>.

package markdown

import (
	"unicode"
	"unicode/utf8"
)

func scanDelims(s *stateInline, start int) (canOpen bool, canClose bool, count int) {
	pos := start
	max := s.posMax
	src := s.src
	marker := src[start]

	lastChar, lastLen := utf8.DecodeLastRuneInString(src[:start])

	for pos < max && src[pos] == marker {
		pos++
	}
	count = pos - start

	nextChar, nextLen := utf8.DecodeRuneInString(src[pos:])

	isLastSpaceOrStart := lastLen == 0 || unicode.IsSpace(lastChar)
	isNextSpaceOrEnd := nextLen == 0 || unicode.IsSpace(nextChar)
	isLastPunct := !isLastSpaceOrStart && (isMarkdownPunct(lastChar) || unicode.IsPunct(lastChar))
	isNextPunct := !isNextSpaceOrEnd && (isMarkdownPunct(nextChar) || unicode.IsPunct(nextChar))

	leftFlanking := !isNextSpaceOrEnd && (!isNextPunct || isLastSpaceOrStart || isLastPunct)
	rightFlanking := !isLastSpaceOrStart && (!isLastPunct || isNextSpaceOrEnd || isNextPunct)

	if marker == '_' {
		canOpen = leftFlanking && (!rightFlanking || isLastPunct)
		canClose = rightFlanking && (!leftFlanking || isNextPunct)
	} else {
		canOpen = leftFlanking
		canClose = rightFlanking
	}

	return
}

var em [256]bool

func init() {
	em['*'], em['_'] = true, true
}

func ruleEmphasis(s *stateInline, silent bool) (_ bool) {
	src := s.src
	max := s.posMax
	start := s.pos
	marker := src[start]

	if !em[marker] {
		return
	}

	if silent {
		return
	}

	canOpen, _, startCount := scanDelims(s, start)
	s.pos += startCount
	if !canOpen {
		s.pending.WriteString(src[start:s.pos])
		return true
	}

	stack := []int{startCount}
	found := false

	for s.pos < max {
		if src[s.pos] == marker {
			canOpen, canClose, count := scanDelims(s, s.pos)

			if canClose {
				oldCount := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				newCount := count

				for oldCount != newCount {
					if newCount < oldCount {
						stack = append(stack, oldCount-newCount)
						break
					}

					newCount -= oldCount

					if len(stack) == 0 {
						break
					}

					s.pos += oldCount
					oldCount = stack[len(stack)-1]
					stack = stack[:len(stack)-1]
				}

				if len(stack) == 0 {
					startCount = oldCount
					found = true
					break
				}

				s.pos += count
				continue
			}

			if canOpen {
				stack = append(stack, count)
			}

			s.pos += count
			continue
		}

		s.md.inline.skipToken(s)
	}

	if !found {
		s.pos = start
		return
	}

	s.posMax = s.pos
	s.pos = start + startCount

	count := startCount
	for ; count > 1; count -= 2 {
		s.pushOpeningToken(&StrongOpen{})
	}
	if count > 0 {
		s.pushOpeningToken(&EmphasisOpen{})
	}

	s.md.inline.tokenize(s)

	if count%2 != 0 {
		s.pushClosingToken(&EmphasisClose{})
	}
	for count = startCount; count > 1; count -= 2 {
		s.pushClosingToken(&StrongClose{})
	}

	s.pos = s.posMax + startCount
	s.posMax = max

	return true
}
