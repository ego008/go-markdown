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

func ruleBackticks(s *stateInline, silent bool) (_ bool) {
	pos := s.pos
	src := s.src

	if src[pos] != '`' {
		return
	}

	start := pos
	pos++
	max := s.posMax

	for pos < max && src[pos] == '`' {
		pos++
	}

	marker := src[start:pos]

	end := pos

	for {
		for start = end; start < max && src[start] != '`'; start++ {
			// do nothing
		}
		if start >= max {
			break
		}
		end = start + 1

		for end < max && src[end] == '`' {
			end++
		}

		if end-start == len(marker) {
			if !silent {
				s.pushToken(&CodeInline{
					Content: normalizeInlineCode(src[pos:start]),
				})
			}
			s.pos = end
			return true
		}
	}

	if !silent {
		s.pending.WriteString(marker)
	}

	s.pos += len(marker)

	return true
}
