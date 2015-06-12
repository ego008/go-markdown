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

func ruleBlockQuote(s *stateBlock, startLine, endLine int, silent bool) (_ bool) {
	shift := s.tShift[startLine]
	if shift < 0 {
		return
	}

	pos := s.bMarks[startLine] + shift
	src := s.src

	if src[pos] != '>' {
		return
	}

	if silent {
		return true
	}

	pos++
	max := s.eMarks[startLine]

	if pos < max && src[pos] == ' ' {
		pos++
	}

	oldIndent := s.blkIndent
	s.blkIndent = 0

	oldBMarks := []int{s.bMarks[startLine]}
	s.bMarks[startLine] = pos

	if pos < max {
		pos = s.skipSpaces(pos)
	}
	lastLineEmpty := pos >= max

	oldTShift := []int{s.tShift[startLine]}
	s.tShift[startLine] = pos - s.bMarks[startLine]

	nextLine := startLine + 1
outer:
	for ; nextLine < endLine; nextLine++ {
		shift := s.tShift[nextLine]
		if shift < oldIndent {
			break
		}

		pos = s.bMarks[nextLine] + shift
		max = s.eMarks[nextLine]

		if pos >= max {
			break
		}

		if src[pos] == '>' {
			pos++
			if pos < max && src[pos] == ' ' {
				pos++
			}

			oldBMarks = append(oldBMarks, s.bMarks[nextLine])
			s.bMarks[nextLine] = pos

			if pos < max {
				pos = s.skipSpaces(pos)
			}
			lastLineEmpty = pos >= max

			oldTShift = append(oldTShift, s.tShift[nextLine])
			s.tShift[nextLine] = pos - s.bMarks[nextLine]

			continue
		}

		if lastLineEmpty {
			break
		}

		for _, r := range []blockRule{
			ruleFence,
			ruleHR,
			ruleList,
			ruleHeading,
			ruleHTMLBlock,
		} {
			if r(s, nextLine-1, endLine, true) {
				break outer
			}
			if r(s, nextLine, endLine, true) {
				break outer
			}
		}

		oldBMarks = append(oldBMarks, s.bMarks[nextLine])
		oldTShift = append(oldTShift, s.tShift[nextLine])

		s.tShift[nextLine] = -1
	}

	oldParentType := s.parentType
	s.parentType = ptBlockQuote
	tok := &BlockquoteOpen{
		Map: [2]int{startLine, 0},
	}
	s.pushOpeningToken(tok)

	s.md.block.tokenize(s, startLine, nextLine)

	s.pushClosingToken(&BlockquoteClose{})
	s.parentType = oldParentType
	tok.Map[1] = s.line

	for i := 0; i < len(oldTShift); i++ {
		s.bMarks[startLine+i] = oldBMarks[i]
		s.tShift[startLine+i] = oldTShift[i]
	}
	s.blkIndent = oldIndent

	return true
}
