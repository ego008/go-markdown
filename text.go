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

var term [256]bool

func init() {
	for _, b := range "\n!#$%&*+-:<=>@[\\]^_`{}~" {
		term[b] = true
	}
}

func ruleText(s *stateInline, silent bool) (_ bool) {
	pos := s.pos
	max := s.posMax
	src := s.src

	for pos < max && !term[src[pos]] {
		pos++
	}
	if pos == s.pos {
		return
	}

	if !silent {
		s.pending.WriteString(src[s.pos:pos])
	}

	s.pos = pos

	return true
}
