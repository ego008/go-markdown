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

type block struct {
}

type blockRule func(*stateBlock, int, int, bool) bool

func (b block) parse(src []byte, md *Markdown, env *environment) []Token {
	str, bMarks, eMarks, tShift := normalizeAndIndex(src)
	bMarks = append(bMarks, len(str))
	eMarks = append(eMarks, len(str))
	tShift = append(tShift, 0)
	var s stateBlock
	s.bMarks = bMarks
	s.eMarks = eMarks
	s.tShift = tShift
	s.lineMax = len(bMarks) - 1
	s.src = str
	s.md = md
	s.env = env

	b.tokenize(&s, s.line, s.lineMax)

	return s.tokens
}

func (block) tokenize(s *stateBlock, startLine, endLine int) {
	line := startLine
	hasEmptyLines := false
	maxNesting := s.md.MaxNesting

	for line < endLine {
		line = s.skipEmptyLines(line)
		s.line = line
		if line >= endLine {
			break
		}

		if s.tShift[line] < s.blkIndent {
			break
		}

		if s.level >= maxNesting {
			s.line = endLine
			break
		}

		for _, r := range []blockRule{
			ruleCode,
			ruleFence,
			ruleBlockQuote,
			ruleHR,
			ruleList,
			ruleReference,
			ruleHeading,
			ruleLHeading,
			ruleHTMLBlock,
			ruleTable,
			ruleParagraph,
		} {
			if r(s, line, endLine, false) {
				break
			}
		}

		s.tight = !hasEmptyLines

		if s.isLineEmpty(s.line - 1) {
			hasEmptyLines = true
		}

		line = s.line

		if line < endLine && s.isLineEmpty(line) {
			hasEmptyLines = true
			line++

			if line < endLine && s.parentType == ptList && s.isLineEmpty(line) {
				break
			}
			s.line = line
		}
	}
}
