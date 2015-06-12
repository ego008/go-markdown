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

var fence [256]bool

func init() {
	fence['~'], fence['`'] = true, true
}

func ruleFence(s *stateBlock, startLine, endLine int, silent bool) (_ bool) {
	shift := s.tShift[startLine]
	if shift < 0 {
		return
	}

	pos := s.bMarks[startLine] + shift
	max := s.eMarks[startLine]
	src := s.src

	if pos+3 > max {
		return
	}

	marker := src[pos]

	if !fence[marker] {
		return
	}

	mem := pos
	pos = s.skipBytes(pos, marker)
	len := pos - mem
	if len < 3 {
		return
	}

	params := strings.TrimSpace(src[pos:max])

	if strings.IndexByte(params, '`') >= 0 {
		return
	}

	if silent {
		return true
	}

	nextLine := startLine
	haveEndMarker := false

	for {
		nextLine++
		if nextLine >= endLine {
			break
		}

		mem = s.bMarks[nextLine] + s.tShift[nextLine]
		pos = mem
		max = s.eMarks[nextLine]

		if pos >= max {
			continue
		}

		if s.tShift[nextLine] < s.blkIndent {
			break
		}

		if src[pos] != marker {
			continue
		}

		if s.tShift[nextLine]-s.blkIndent > 3 {
			continue
		}

		pos = s.skipBytes(pos, marker)

		if pos-mem < len {
			continue
		}

		pos = s.skipSpaces(pos)
		if pos < max {
			continue
		}

		haveEndMarker = true

		break
	}

	s.line = nextLine
	if haveEndMarker {
		s.line++
	}

	s.pushToken(&Fence{
		Params:  params,
		Content: s.lines(startLine+1, nextLine, s.tShift[startLine], true),
		Map:     [2]int{startLine, nextLine},
	})

	return true
}
