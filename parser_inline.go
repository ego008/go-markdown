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

type inline struct {
}

type inlineRule func(*stateInline, bool) bool

var inlineRules []inlineRule

func init() {
	inlineRules = []inlineRule{
		ruleText,
		ruleNewline,
		ruleEscape,
		ruleBackticks,
		ruleStrikeThrough,
		ruleEmphasis,
		ruleLink,
		ruleImage,
		ruleAutolink,
		ruleHTMLInline,
		ruleEntity,
	}
}

func (i inline) parse(src string, md *Markdown, env *environment) []Token {
	if src == "" {
		return nil
	}

	var s stateInline
	s.src = src
	s.md = md
	s.env = env
	s.pos = 0
	s.posMax = len(src)
	s.level = 0
	s.pendingLevel = 0

	i.tokenize(&s)

	return s.tokens
}

func (inline) tokenize(s *stateInline) {
	max := s.posMax
	src := s.src
	maxNesting := s.md.MaxNesting

outer:
	for s.pos < max {
		if s.level < maxNesting {
			for _, rule := range inlineRules {
				if rule(s, false) {
					if s.pos >= max {
						break outer
					}
					continue outer
				}
			}
		}

		r, size := utf8.DecodeRuneInString(src[s.pos:])
		s.pending.WriteRune(r)
		s.pos += size
	}

	if s.pending.Len() > 0 {
		s.pushPending()
	}
}

func (inline) skipToken(s *stateInline) {
	pos := s.pos
	if s.cache != nil {
		if pos, ok := s.cache[pos]; ok {
			s.pos = pos
			return
		}
	} else {
		s.cache = make(map[int]int)
	}

	if s.level < s.md.MaxNesting {
		for _, r := range inlineRules {
			if r(s, true) {
				s.cache[pos] = s.pos
				return
			}
		}
	}

	s.pos++
	s.cache[pos] = s.pos
}
