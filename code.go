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

func ruleCode(s *stateBlock, startLine, endLine int, _ bool) (_ bool) {
	if s.tShift[startLine]-s.blkIndent < 4 {
		return
	}

	nextLine := startLine + 1
	last := nextLine

	for nextLine < endLine {
		if s.isLineEmpty(nextLine) {
			nextLine++
			continue
		}

		if s.tShift[nextLine]-s.blkIndent > 3 {
			nextLine++
			last = nextLine
			continue
		}

		break
	}

	s.line = nextLine
	s.pushToken(&CodeBlock{
		Content: s.lines(startLine, last, 4+s.blkIndent, true),
		Map:     [2]int{startLine, s.line},
	})

	return true
}
