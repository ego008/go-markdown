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

func ruleAutolink(s *stateInline, silent bool) (_ bool) {
	pos := s.pos
	src := s.src

	if src[pos] != '<' {
		return
	}

	tail := src[pos:]

	if strings.IndexByte(tail, '>') < 0 {
		return
	}

	link := matchAutolink(tail)
	if link != "" {
		href := normalizeLink(link)
		if !validateLink(href) {
			return
		}

		if !silent {
			s.pushOpeningToken(&LinkOpen{Href: href})
			s.pushToken(&Text{Content: normalizeLinkText(link)})
			s.pushClosingToken(&LinkClose{})
		}

		s.pos += len(link) + 2

		return true
	}

	email := matchEmail(tail)
	if email != "" {
		href := normalizeLink("mailto:" + email)
		if !validateLink(href) {
			return
		}

		if !silent {
			s.pushOpeningToken(&LinkOpen{Href: href})
			s.pushToken(&Text{Content: email})
			s.pushClosingToken(&LinkClose{})
		}

		s.pos += len(email) + 2

		return true
	}

	return
}
