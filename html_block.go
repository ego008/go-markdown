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

import "github.com/opennota/byteutil"

var (
	htmlBlocks = []string{
		"article",
		"aside",
		"blockquote",
		"body",
		"button",
		"canvas",
		"caption",
		"col",
		"colgroup",
		"dd",
		"div",
		"dl",
		"dt",
		"embed",
		"fieldset",
		"figcaption",
		"figure",
		"footer",
		"form",
		"h1",
		"h2",
		"h3",
		"h4",
		"h5",
		"h6",
		"header",
		"hgroup",
		"hr",
		"iframe",
		"li",
		"map",
		"object",
		"ol",
		"output",
		"p",
		"pre",
		"progress",
		"script",
		"section",
		"style",
		"table",
		"tbody",
		"td",
		"textarea",
		"tfoot",
		"th",
		"thead",
		"tr",
		"ul",
		"video",
	}

	htmlBlocksSet = make(map[string]bool)

	htmlSecond    [256]bool
	piOrComment   [256]bool
	slashOrLetter [256]bool
)

func init() {
	for _, tag := range htmlBlocks {
		htmlBlocksSet[tag] = true
	}
	for _, b := range "!/?abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		htmlSecond[b] = true
	}
	piOrComment['!'], piOrComment['?'] = true, true
	for _, b := range "/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		slashOrLetter[b] = true
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func matchTagName(s string) string {
	if len(s) < 2 {
		return ""
	}

	i := 0
	if s[0] == '/' {
		i++
	}
	start := i
	max := min(15+i, len(s))
	for i < max && byteutil.IsLetter(s[i]) {
		i++
	}
	if i >= len(s) {
		return ""
	}

	switch s[i] {
	case ' ', '\n', '/', '>':
		return byteutil.ToLower(s[start:i])
	default:
		return ""
	}
}

func ruleHTMLBlock(s *stateBlock, startLine, endLine int, silent bool) (_ bool) {
	if !s.md.HTML {
		return
	}

	shift := s.tShift[startLine]
	if shift < 0 {
		return
	}

	pos := s.bMarks[startLine] + shift
	max := s.eMarks[startLine]

	if shift > 3 {
		return
	}

	if pos+2 >= max {
		return
	}

	src := s.src

	if src[pos] != '<' {
		return
	}

	b := src[pos+1]
	if !htmlSecond[b] {
		return
	}

	if slashOrLetter[b] {
		tag := matchTagName(src[pos+1 : max])
		if tag == "" {
			return
		}
		if !htmlBlocksSet[tag] {
			return
		}
		if silent {
			return true
		}
	} else if piOrComment[b] {
		if silent {
			return true
		}
	} else {
		return
	}

	nextLine := startLine + 1
	for nextLine < s.lineMax && !s.isLineEmpty(nextLine) {
		nextLine++
	}

	s.line = nextLine
	s.pushToken(&HTMLBlock{
		Content: s.lines(startLine, nextLine, 0, true),
		Map:     [2]int{startLine, s.line},
	})

	return true
}
