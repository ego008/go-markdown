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

import "unicode/utf8"

const runeErrorStr = string(utf8.RuneError)

var special [256]bool

func init() {
	for _, b := range "\x00\t\n\r " {
		special[b] = true
	}
}

func normalizeAndIndex(src []byte) (s string, bMarks []int, eMarks []int, tShift []int) {
	buf := make([]byte, len(src)*4)
	i := 0
	j := 0
	pos := 0
	skipNextLf := false
	lineStart := 0
	lastTabPos := 0
	indent := 0
	indentFound := false
	start := 0

	for pos < len(src) {
		r, size := utf8.DecodeRune(src[pos:])
		pos += size

		if skipNextLf {
			skipNextLf = false
			if r == '\n' {
				continue
			}
		}

		if !(r <= 0x20 && special[r]) {
			j += utf8.EncodeRune(buf[j:], r)
			indentFound = true
			i++
			continue
		}

		switch r {
		case ' ':
			buf[j] = ' '
			j++

			if !indentFound {
				indent++
			}
		case '\r':
			skipNextLf = true
			fallthrough
		case '\n':
			bMarks = append(bMarks, start)
			eMarks = append(eMarks, j)
			tShift = append(tShift, indent)
			indentFound = false
			indent = 0

			buf[j] = '\n'
			j++
			lineStart = i + 1
			lastTabPos = 0

			start = j
		case '\t':
			k := (i - lineStart - lastTabPos) % 4
			j += copy(buf[j:], "    "[k:])
			lastTabPos = i - lineStart + 1

			if !indentFound {
				indent += 4 - k
			}
		case '\x00':
			j += copy(buf[j:], runeErrorStr)
			indentFound = true
		}

		i++
	}

	if j > 0 && buf[j-1] != '\n' {
		bMarks = append(bMarks, start)
		eMarks = append(eMarks, j)
		tShift = append(tShift, indent)
	}

	s = string(buf[:j])
	return
}
