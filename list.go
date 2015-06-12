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
	"strconv"

	"github.com/opennota/byteutil"
)

var (
	bullet     [256]bool
	afterdigit [256]bool
)

func init() {
	bullet['*'], bullet['+'], bullet['-'] = true, true, true
	afterdigit[')'], afterdigit['.'] = true, true
}

func skipBulletListMarker(s *stateBlock, startLine int) int {
	pos := s.bMarks[startLine] + s.tShift[startLine]
	src := s.src

	if !bullet[src[pos]] {
		return -1
	}

	pos++
	max := s.eMarks[startLine]

	if pos < max && src[pos] != ' ' {
		return -1
	}

	return pos
}

func skipOrderedListMarker(s *stateBlock, startLine int) int {
	pos := s.bMarks[startLine] + s.tShift[startLine]
	max := s.eMarks[startLine]

	if pos+1 >= max {
		return -1
	}

	src := s.src
	b := src[pos]

	if !byteutil.IsDigit(b) {
		return -1
	}

	for {
		if pos >= max {
			return -1
		}

		b = src[pos]
		pos++

		if byteutil.IsDigit(b) {
			continue
		}

		if afterdigit[b] {
			break
		}

		return -1
	}

	if pos < max && src[pos] != ' ' {
		return -1
	}

	return pos
}

func markParagraphsTight(s *stateBlock, idx int) {
	level := s.level + 2
	tokens := s.tokens

	for i := idx + 2; i < len(tokens)-2; i++ {
		if tokens[i].Level() == level {
			if tok, ok := tokens[i].(*ParagraphOpen); ok {
				tok.Tight = true
				i += 2
				tokens[i].(*ParagraphClose).Tight = true
			}
		}
	}
}

func ruleList(s *stateBlock, startLine, endLine int, silent bool) (_ bool) {
	shift := s.tShift[startLine]
	if shift < 0 {
		return
	}

	isOrdered := false
	posAfterMarker := skipOrderedListMarker(s, startLine)
	if posAfterMarker > 0 {
		isOrdered = true
	} else {
		posAfterMarker = skipBulletListMarker(s, startLine)
		if posAfterMarker < 0 {
			return
		}
	}

	src := s.src
	markerChar := src[posAfterMarker-1]

	if silent {
		return true
	}

	tokenIdx := len(s.tokens)

	var listMap *[2]int
	if isOrdered {
		start := s.bMarks[startLine] + shift
		markerValue, _ := strconv.Atoi(src[start : posAfterMarker-1])

		tok := &OrderedListOpen{
			Order: markerValue,
			Map:   [2]int{startLine, 0},
		}
		s.pushOpeningToken(tok)
		listMap = &tok.Map
	} else {
		tok := &BulletListOpen{
			Map: [2]int{startLine, 0},
		}
		s.pushOpeningToken(tok)
		listMap = &tok.Map
	}

	nextLine := startLine
	prevEmptyEnd := false

	tight := true
outer:
	for nextLine < endLine {
		contentStart := s.skipSpaces(posAfterMarker)
		max := s.eMarks[nextLine]

		var indentAfterMarker int
		if contentStart >= max {
			indentAfterMarker = 1
		} else {
			indentAfterMarker = contentStart - posAfterMarker
		}

		if indentAfterMarker > 4 {
			indentAfterMarker = 1
		}

		indent := posAfterMarker - s.bMarks[nextLine] + indentAfterMarker

		tok := &ListItemOpen{
			Map: [2]int{startLine, 0},
		}
		s.pushOpeningToken(tok)
		itemMap := &tok.Map

		oldIndent := s.blkIndent
		oldTight := s.tight
		oldTShift := s.tShift[startLine]
		oldParentType := s.parentType
		s.tShift[startLine] = contentStart - s.bMarks[startLine]
		s.blkIndent = indent
		s.tight = true
		s.parentType = ptList

		s.md.block.tokenize(s, startLine, endLine)

		if !s.tight || prevEmptyEnd {
			tight = false
		}
		prevEmptyEnd = s.line-startLine > 1 && s.isLineEmpty(s.line-1)
		if prevEmptyEnd {
			lastToken := s.tokens[len(s.tokens)-1]
			if _, ok := lastToken.(*BlockquoteClose); ok {
				prevEmptyEnd = false
			}
		}

		s.blkIndent = oldIndent
		s.tShift[startLine] = oldTShift
		s.tight = oldTight
		s.parentType = oldParentType

		s.pushClosingToken(&ListItemClose{})

		startLine = s.line
		nextLine = startLine
		(*itemMap)[1] = nextLine
		contentStart = s.bMarks[startLine]

		if nextLine >= endLine {
			break
		}

		if s.isLineEmpty(nextLine) {
			break
		}

		if s.tShift[nextLine] < s.blkIndent {
			break
		}

		for _, r := range []blockRule{
			ruleFence,
			ruleBlockQuote,
			ruleHR,
		} {
			if r(s, nextLine, endLine, true) {
				break outer
			}
		}

		if isOrdered {
			posAfterMarker = skipOrderedListMarker(s, nextLine)
			if posAfterMarker < 0 {
				break
			}
		} else {
			posAfterMarker = skipBulletListMarker(s, nextLine)
			if posAfterMarker < 0 {
				break
			}
		}

		if markerChar != src[posAfterMarker-1] {
			break
		}
	}

	if isOrdered {
		s.pushClosingToken(&OrderedListClose{})
	} else {
		s.pushClosingToken(&BulletListClose{})
	}
	(*listMap)[1] = nextLine

	s.line = nextLine

	if tight {
		markParagraphsTight(s, tokenIdx)
	}

	return true
}
