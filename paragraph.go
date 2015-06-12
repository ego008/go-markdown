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

import "strings"

func ruleParagraph(s *stateBlock, startLine, _ int, _ bool) bool {
	nextLine := startLine + 1
	endLine := s.lineMax

outer:
	for ; nextLine < endLine && !s.isLineEmpty(nextLine); nextLine++ {
		shift := s.tShift[nextLine]
		if shift < 0 || shift-s.blkIndent > 3 {
			continue
		}

		for _, r := range []blockRule{
			ruleFence,
			ruleBlockQuote,
			ruleHR,
			ruleList,
			ruleHeading,
			ruleHTMLBlock,
			ruleTable,
		} {
			if r(s, nextLine, endLine, true) {
				break outer
			}
		}
	}

	content := strings.TrimSpace(s.lines(startLine, nextLine, s.blkIndent, false))

	s.line = nextLine

	s.pushOpeningToken(&ParagraphOpen{
		Map: [2]int{startLine, s.line},
	})
	s.pushToken(&Inline{
		Content: content,
		Map:     [2]int{startLine, s.line},
	})
	s.pushClosingToken(&ParagraphClose{
		Map: [2]int{startLine, s.line},
	})

	return true
}
