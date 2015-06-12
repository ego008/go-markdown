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

import (
	"strings"
	"unicode"
)

type stackItem struct {
	token  int
	text   []rune
	pos    int
	single bool
	level  int
}

func nextQuoteIndex(s []rune, from int) int {
	for i := from; i < len(s); i++ {
		r := s[i]
		if r == '\'' || r == '"' {
			return i
		}
	}
	return -1
}

func replaceQuotes(tokens []Token, s *stateCore) {
	var stack []stackItem
	var changed map[int][]rune

	for i, tok := range tokens {
		thisLevel := tok.Level()
		_ = thisLevel

		j := len(stack) - 1
		for j >= 0 {
			if stack[j].level <= thisLevel {
				break
			}
			j--
		}
		stack = stack[:j+1]

		if tok, ok := tok.(*Text); ok {
			if !strings.ContainsAny(tok.Content, `"'`) {
				continue
			}

			text := []rune(tok.Content)
			pos := 0
			max := len(text)

		loop:
			for pos < max {
				index := nextQuoteIndex(text, pos)
				if index < 0 {
					break
				}

				canOpen := true
				canClose := true
				pos = index + 1
				isSingle := text[index] == '\''

				var lastChar rune
				if index > 0 {
					lastChar = text[index-1]
				}
				var nextChar rune
				if pos < max {
					nextChar = text[pos]
				}

				isLastSpaceOrStart := index == 0 || unicode.IsSpace(lastChar)
				isNextSpaceOrEnd := pos == max || unicode.IsSpace(nextChar)
				isLastPunct := !isLastSpaceOrStart && (isMarkdownPunct(lastChar) || unicode.IsPunct(lastChar))
				isNextPunct := !isNextSpaceOrEnd && (isMarkdownPunct(nextChar) || unicode.IsPunct(nextChar))

				if isNextSpaceOrEnd {
					canOpen = false
				} else if isNextPunct {
					if !(isLastSpaceOrStart || isLastPunct) {
						canOpen = false
					}
				}

				if isLastSpaceOrStart {
					canClose = false
				} else if isLastPunct {
					if !(isNextSpaceOrEnd || isNextPunct) {
						canClose = false
					}
				}

				if nextChar == '"' && text[index] == '"' {
					if lastChar >= '0' && lastChar <= '9' {
						canClose = false
						canOpen = false
					}
				}

				if canOpen && canClose {
					canOpen = false
					canClose = isNextPunct
				}

				if !canOpen && !canClose {
					if isSingle {
						text[index] = '’'
						if changed == nil {
							changed = make(map[int][]rune)
						}
						if _, ok := changed[i]; !ok {
							changed[i] = text
						}
					}
					continue
				}

				if canClose {
					for j := len(stack) - 1; j >= 0; j-- {
						item := stack[j]
						if item.level < thisLevel {
							break
						}
						if item.single == isSingle && item.level == thisLevel {
							if changed == nil {
								changed = make(map[int][]rune)
							}
							if isSingle {
								item.text[item.pos] = s.md.options.Quotes[2]
								text[index] = s.md.options.Quotes[3]
							} else {
								item.text[item.pos] = s.md.options.Quotes[0]
								text[index] = s.md.options.Quotes[1]
							}
							if _, ok := changed[i]; !ok {
								changed[i] = text
							}
							if ii := item.token; ii != i {
								if _, ok := changed[ii]; !ok {
									changed[ii] = item.text
								}
							}
							stack = stack[:j]
							continue loop
						}
					}
				}

				if canOpen {
					stack = append(stack, stackItem{
						token:  i,
						text:   text,
						pos:    index,
						single: isSingle,
						level:  thisLevel,
					})
				} else if canClose && isSingle {
					text[index] = '’'
					if changed == nil {
						changed = make(map[int][]rune)
					}
					if _, ok := changed[i]; !ok {
						changed[i] = text
					}
				}
			}
		}
	}

	if changed != nil {
		for i, text := range changed {
			tokens[i].(*Text).Content = string(text)
		}
	}
}

func ruleSmartQuotes(s *stateCore) {
	if !s.md.Typographer {
		return
	}

	tokens := s.tokens
	for i := len(tokens) - 1; i >= 0; i-- {
		tok := tokens[i]
		if tok, ok := tok.(*Inline); ok {
			replaceQuotes(tok.Children, s)
		}
	}
}
