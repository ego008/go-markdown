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

func ruleStrikeThrough(s *stateInline, silent bool) (_ bool) {
	start := s.pos
	max := s.posMax
	src := s.src

	if src[start] != '~' {
		return
	}

	if silent {
		return
	}

	canOpen, canClose, delims := scanDelims(s, start)
	startCount := delims
	if !canOpen {
		s.pos += startCount
		s.pending.WriteString(src[start:s.pos])
		return true
	}

	stack := startCount / 2
	if stack <= 0 {
		return
	}
	s.pos = start + startCount

	var found bool
	for s.pos < max {
		if src[s.pos] == '~' {
			canOpen, canClose, delims = scanDelims(s, s.pos)
			count := delims
			tagCount := count / 2
			if canClose {
				if tagCount >= stack {
					s.pos += count - 2
					found = true
					break
				}
				stack -= tagCount
				s.pos += count
				continue
			}

			if canOpen {
				stack += tagCount
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
	s.pos = start + 2

	s.pushOpeningToken(&StrikethroughOpen{})

	s.md.inline.tokenize(s)

	s.pushClosingToken(&StrikethroughClose{})

	s.pos = s.posMax + 2
	s.posMax = max

	return true
}
