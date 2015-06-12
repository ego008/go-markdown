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

	"github.com/opennota/byteutil"
	"github.com/opennota/linkify"
)

func isLinkOpen(s string) bool {
	return byteutil.IsLetter(s[1])
}

func isLinkClose(s string) bool {
	return s[1] == '/'
}

func ruleLinkify(s *stateCore) {
	blockTokens := s.tokens

	if !s.md.Linkify {
		return
	}

	for _, tok := range blockTokens {
		if tok, ok := tok.(*Inline); ok {
			tokens := tok.Children

			htmlLinkLevel := 0

			for i := len(tokens) - 1; i >= 0; i-- {
				currentTok := tokens[i]

				if _, ok := currentTok.(*LinkClose); ok {
					i--
					for tokens[i].Level() != currentTok.Level() {
						if _, ok := tokens[i].(*LinkOpen); ok {
							break
						}
						i--
					}
					continue
				}

				if currentTok, ok := currentTok.(*HTMLInline); ok {
					if isLinkOpen(currentTok.Content) && htmlLinkLevel > 0 {
						htmlLinkLevel--
					}
					if isLinkClose(currentTok.Content) {
						htmlLinkLevel++
					}
				}
				if htmlLinkLevel > 0 {
					continue
				}

				if currentTok, ok := currentTok.(*Text); ok {
					text := currentTok.Content
					links := linkify.Links(text)
					if len(links) == 0 {
						continue
					}

					var nodes []Token
					level := currentTok.Lvl
					lastPos := 0

					for _, ln := range links {
						urlText := text[ln.Start:ln.End]
						url := urlText
						if ln.Schema == "" {
							url = "http://" + url
						} else if ln.Schema == "mailto:" && !strings.HasPrefix(url, "mailto:") {
							url = "mailto:" + url
						}
						url = normalizeLink(url)
						if !validateLink(url) {
							continue
						}

						urlText = normalizeLinkText(urlText)

						pos := ln.Start

						if pos > lastPos {
							tok := Text{
								Content: text[lastPos:pos],
								Lvl:     level,
							}
							nodes = append(nodes, &tok)
						}

						nodes = append(nodes, &LinkOpen{
							Href: url,
							Lvl:  level,
						})
						nodes = append(nodes, &Text{
							Content: urlText,
							Lvl:     level + 1,
						})
						nodes = append(nodes, &LinkClose{
							Lvl: level,
						})

						lastPos = ln.End
					}

					if lastPos < len(text) {
						tok := Text{
							Content: text[lastPos:],
							Lvl:     level,
						}
						nodes = append(nodes, &tok)
					}

					children := make([]Token, len(tokens)+len(nodes)-1)
					copy(children, tokens[:i])
					copy(children[i:], nodes)
					copy(children[i+len(nodes):], tokens[i+1:])
					tok.Children = children
					tokens = children
				}
			}
		}
	}
}
