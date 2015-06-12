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

func ruleImage(s *stateInline, silent bool) (_ bool) {
	pos := s.pos
	max := s.posMax

	if pos+2 >= max {
		return
	}

	src := s.src
	if src[pos] != '!' {
		return
	}

	if src[pos+1] != '[' {
		return
	}

	labelStart := pos + 2
	labelEnd := parseLinkLabel(s, pos+1, false)
	if labelEnd < 0 {
		return
	}

	var href, title, label string
	oldPos := pos
	pos = labelEnd + 1
	if pos < max && src[pos] == '(' {
		pos = skipws(src, pos+1, max)
		if pos >= max {
			return
		}

		url, endpos, ok := parseLinkDestination(src, pos, s.posMax)
		if ok {
			url = normalizeLink(url)
			if validateLink(url) {
				href = url
				pos = endpos
			}
		}

		start := pos
		pos = skipws(src, pos, max)
		if pos >= max {
			return
		}

		title, _, endpos, ok = parseLinkTitle(src, pos, s.posMax)
		if pos < max && start != pos && ok {
			pos = skipws(src, endpos, max)
		}

		if pos >= max || src[pos] != ')' {
			s.pos = oldPos
			return
		}

		pos++

	} else {
		if s.env.References == nil {
			return
		}

		pos = skipws(src, pos, max)

		if pos < max && src[pos] == '[' {
			start := pos + 1
			pos = parseLinkLabel(s, pos, false)
			if pos >= 0 {
				label = src[start:pos]
				pos++
			} else {
				pos = labelEnd + 1
			}
		} else {
			pos = labelEnd + 1
		}

		if label == "" {
			label = src[labelStart:labelEnd]
		}

		ref, ok := s.env.References[normalizeReference(label)]
		if !ok {
			s.pos = oldPos
			return
		}

		href = ref["href"]
		title = ref["title"]
	}

	if !silent {
		s.pos = labelStart
		s.posMax = labelEnd

		src := src[labelStart:labelEnd]

		var newState stateInline
		newState.src = src
		newState.md = s.md
		newState.env = s.env
		newState.posMax = len(src)
		newState.md.inline.tokenize(&newState)

		s.pushToken(&Image{
			Src:    href,
			Title:  title,
			Tokens: newState.tokens,
		})
	}

	s.pos = pos
	s.posMax = max

	return true
}
