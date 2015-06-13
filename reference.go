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

func ruleReference(s *stateBlock, startLine, _ int, silent bool) (_ bool) {
	pos := s.bMarks[startLine] + s.tShift[startLine]
	src := s.src

	if src[pos] != '[' {
		return
	}

	pos++
	max := s.eMarks[startLine]

	for pos < max {
		if src[pos] == ']' && src[pos-1] != '\\' {
			if pos+1 == max {
				return
			}
			if src[pos+1] != ':' {
				return
			}
			break
		}
		pos++
	}

	nextLine := startLine + 1
	endLine := s.lineMax
outer:
	for ; nextLine < endLine && !s.isLineEmpty(nextLine); nextLine++ {
		if s.tShift[nextLine]-s.blkIndent > 3 {
			continue
		}

		for _, r := range []blockRule{
			ruleFence,
			ruleBlockQuote,
			ruleHR,
			ruleList,
			ruleHeading,
			ruleHTMLBlock,
			ruleTable,
		} {
			if r(s, nextLine, endLine, true) {
				break outer
			}
		}
	}

	str := strings.TrimSpace(s.lines(startLine, nextLine, s.blkIndent, false))
	max = len(str)
	lines := 0
	var labelEnd int
	for pos = 1; pos < max; pos++ {
		b := str[pos]
		if b == '[' {
			return
		} else if b == ']' {
			labelEnd = pos
			break
		} else if b == '\n' {
			lines++
		} else if b == '\\' {
			pos++
			if pos < max && str[pos] == '\n' {
				lines++
			}
		}
	}

	if labelEnd <= 0 || labelEnd+1 >= max || str[labelEnd+1] != ':' {
		return
	}

	for pos = labelEnd + 2; pos < max; pos++ {
		b := str[pos]
		if b == '\n' {
			lines++
		} else if b != ' ' {
			break
		}
	}

	href, endpos, ok := parseLinkDestination(str, pos, max)
	if !ok {
		return
	}
	href = normalizeLink(href)
	if !validateLink(href) {
		return
	}
	pos = endpos

	savedPos := pos
	savedLineNo := lines

	start := pos
	for ; pos < max; pos++ {
		b := str[pos]
		if b == '\n' {
			lines++
		} else if b != ' ' {
			break
		}
	}

	title, nlines, endpos, ok := parseLinkTitle(str, pos, max)
	if pos < max && start != pos && ok {
		pos = endpos
		lines += nlines
	} else {
		pos = savedPos
		lines = savedLineNo
	}

	for pos < max && str[pos] == ' ' {
		pos++
	}

	if pos < max && str[pos] != '\n' {
		return
	}

	label := normalizeReference(str[1:labelEnd])
	if label == "" {
		return false
	}

	if silent {
		return true
	}

	if s.env.References == nil {
		s.env.References = make(map[string]map[string]string)
	}
	if _, ok := s.env.References[label]; !ok {
		s.env.References[label] = map[string]string{
			"title": title,
			"href":  href,
		}
	}

	s.line = startLine + lines + 1

	return true
}
