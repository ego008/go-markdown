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

var under [256]bool

func init() {
	under['-'], under['='] = true, true
}

func ruleLHeading(s *stateBlock, startLine, endLine int, silent bool) (_ bool) {
	nextLine := startLine + 1

	if nextLine >= endLine {
		return
	}

	shift := s.tShift[nextLine]
	if shift < s.blkIndent {
		return
	}

	if shift-s.blkIndent > 3 {
		return
	}

	pos := s.bMarks[nextLine] + shift
	max := s.eMarks[nextLine]

	if pos >= max {
		return
	}

	src := s.src
	marker := src[pos]

	if !under[marker] {
		return
	}

	pos = s.skipBytes(pos, marker)

	pos = s.skipSpaces(pos)

	if pos < max {
		return
	}

	pos = s.bMarks[startLine] + s.tShift[startLine]

	s.line = nextLine + 1

	hLevel := 1
	if marker == '-' {
		hLevel++
	}

	s.pushOpeningToken(&HeadingOpen{
		HLevel: hLevel,
		Map:    [2]int{startLine, s.line},
	})
	s.pushToken(&Inline{
		Content: strings.TrimSpace(src[pos:s.eMarks[startLine]]),
		Map:     [2]int{startLine, s.line - 1},
	})
	s.pushClosingToken(&HeadingClose{HLevel: hLevel})

	return true
}
