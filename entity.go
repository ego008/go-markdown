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

import "github.com/opennota/html"

func ruleEntity(s *stateInline, silent bool) (_ bool) {
	pos := s.pos
	src := s.src

	if src[pos] != '&' {
		return
	}

	max := s.posMax

	if pos+1 < max {
		if e, n := html.ParseEntity(src[pos:]); n > 0 {
			if !silent {
				s.pending.WriteString(e)
			}
			s.pos += n
			return true
		}
	}

	if !silent {
		s.pending.WriteByte('&')
	}
	s.pos++

	return true
}
