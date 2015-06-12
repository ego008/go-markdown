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

func parseLinkLabel(s *stateInline, start int, disableNested bool) int { //XXX get rid of s param
	src := s.src
	labelEnd := -1
	max := s.posMax
	oldPos := s.pos

	s.pos = start + 1
	level := 1
	found := false

	for s.pos < max {
		marker := src[s.pos]

		if marker == ']' {
			level--
			if level == 0 {
				found = true
				break
			}
		}

		prevPos := s.pos

		s.md.inline.skipToken(s)

		if marker == '[' {
			if prevPos == s.pos-1 {
				level++
			} else if disableNested {
				s.pos = oldPos
				return -1
			}
		}
	}

	if found {
		labelEnd = s.pos
	}

	s.pos = oldPos

	return labelEnd
}

func parseLinkDestination(s string, pos, max int) (url string, endpos int, ok bool) {
	start := pos
	if pos < max && s[pos] == '<' {
		pos++
		for pos < max {
			b := s[pos]
			if b == '\n' {
				return
			}
			if b == '>' {
				endpos = pos + 1
				url = unescapeAll(s[start+1 : pos])
				ok = true
				return
			}
			if b == '\\' && pos+1 < max {
				pos += 2
				continue
			}

			pos++
		}

		return
	}

	level := 0
	for pos < max {
		b := s[pos]

		if b == ' ' {
			break
		}

		if b < 0x20 || b == 0x7f {
			break
		}

		if b == '\\' && pos+1 < max {
			pos += 2
			continue
		}

		if b == '(' {
			level++
			if level > 1 {
				break
			}
		}

		if b == ')' {
			level--
			if level < 0 {
				break
			}
		}

		pos++
	}

	if start == pos {
		return
	}

	url = unescapeAll(s[start:pos])
	endpos = pos
	ok = true

	return
}

var tmark [256]bool

func init() {
	tmark['"'], tmark['\''], tmark['('] = true, true, true
}

func parseLinkTitle(s string, pos, max int) (title string, nlines, endpos int, ok bool) {
	lines := 0
	start := pos

	if pos >= max {
		return
	}

	marker := s[pos]

	if !tmark[marker] {
		return
	}

	pos++

	if marker == '(' {
		marker = ')'
	}

	for pos < max {
		switch s[pos] {
		case marker:
			endpos = pos + 1
			nlines = lines
			title = unescapeAll(s[start+1 : pos])
			ok = true
			return
		case '\n':
			lines++
		case '\\':
			if pos+1 < max {
				pos++
				if s[pos] == '\n' {
					lines++
				}
			}
		}
		pos++
	}

	return
}
