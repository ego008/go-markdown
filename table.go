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

var hdr [256]bool

func init() {
	for _, b := range "-:| " {
		hdr[b] = true
	}
}

func getLine(s *stateBlock, line int) string {
	pos := s.bMarks[line] + s.blkIndent
	max := s.eMarks[line]
	if pos >= max {
		return ""
	}
	return s.src[pos:max]
}

func escapedSplit(s string) (result []string) {
	escapes := 0
	lastPos := 0
	backticked := false
	lastBackTick := 0
	pos := 0

	if len(s) > 0 && s[len(s)-1] == '|' {
		s = s[:len(s)-1]
	}
	if len(s) > 0 && s[0] == '|' {
		pos++
		lastPos++
	}

	for pos < len(s) {
		b := s[pos]
		switch {
		case b == '`' && escapes%2 == 0:
			backticked = !backticked
			lastBackTick = pos
		case b == '|' && escapes%2 == 0 && !backticked:
			result = append(result, s[lastPos:pos])
			lastPos = pos + 1
		case b == '\\':
			escapes++
		default:
			escapes = 0
		}

		pos++

		if pos == len(s) && backticked {
			backticked = false
			pos = lastBackTick + 1
		}
	}

	result = append(result, s[lastPos:])

	return
}

func ruleTable(s *stateBlock, startLine, endLine int, silent bool) (_ bool) {
	if !s.md.Tables {
		return
	}

	if startLine+2 > endLine {
		return
	}

	nextLine := startLine + 1

	if s.tShift[nextLine] < s.blkIndent {
		return
	}

	pos := s.bMarks[nextLine] + s.tShift[nextLine]
	if pos >= s.eMarks[nextLine] {
		return
	}

	src := s.src
	if !hdr[src[pos]] {
		return
	}

	lineText := getLine(s, startLine+1)
	if !isHeaderLine(lineText) {
		return
	}

	rows := strings.Split(lineText, "|")
	if len(rows) < 2 {
		return
	}
	var aligns []Align
	for i := 0; i < len(rows); i++ {
		t := strings.TrimSpace(rows[i])
		if t == "" {
			continue
		}

		if t[len(t)-1] == ':' {
			if t[0] == ':' {
				aligns = append(aligns, AlignCenter)
			} else {
				aligns = append(aligns, AlignRight)
			}
		} else if t[0] == ':' {
			aligns = append(aligns, AlignLeft)
		} else {
			aligns = append(aligns, AlignNone)
		}
	}

	lineText = strings.TrimSpace(getLine(s, startLine))
	if strings.IndexByte(lineText, '|') == -1 {
		return
	}

	rows = escapedSplit(lineText)
	if len(aligns) != len(rows) {
		return
	}

	if silent {
		return true
	}

	tableTok := &TableOpen{
		Map: [2]int{startLine, 0},
	}
	s.pushOpeningToken(tableTok)
	s.pushOpeningToken(&TheadOpen{
		Map: [2]int{startLine, startLine + 1},
	})
	s.pushOpeningToken(&TrOpen{
		Map: [2]int{startLine, startLine + 1},
	})

	for i := 0; i < len(rows); i++ {
		s.pushOpeningToken(&ThOpen{
			Align: aligns[i],
			Map:   [2]int{startLine, startLine + 1},
		})
		s.pushToken(&Inline{
			Content: strings.TrimSpace(rows[i]),
			Map:     [2]int{startLine, startLine + 1},
		})
		s.pushClosingToken(&ThClose{})
	}

	s.pushClosingToken(&TrClose{})
	s.pushClosingToken(&TheadClose{})

	tbodyTok := &TbodyOpen{
		Map: [2]int{startLine + 2, 0},
	}
	s.pushOpeningToken(tbodyTok)

	for nextLine = startLine + 2; nextLine < endLine; nextLine++ {
		shift := s.tShift[nextLine]
		if shift >= 0 && shift < s.blkIndent {
			break
		}

		lineText = strings.TrimSpace(getLine(s, nextLine))
		if strings.IndexByte(lineText, '|') == -1 {
			break
		}
		rows = escapedSplit(lineText)
		if len(rows) < len(aligns) {
			rows = append(rows, make([]string, len(aligns)-len(rows))...)
		} else if len(rows) > len(aligns) {
			rows = rows[:len(aligns)]
		}

		s.pushOpeningToken(&TrOpen{})
		for i := 0; i < len(rows); i++ {
			s.pushOpeningToken(&TdOpen{Align: aligns[i]})
			s.pushToken(&Inline{
				Content: strings.TrimSpace(rows[i]),
			})
			s.pushClosingToken(&TdClose{})
		}
		s.pushClosingToken(&TrClose{})
	}

	s.pushClosingToken(&TbodyClose{})
	s.pushClosingToken(&TableClose{})

	tableTok.Map[1] = nextLine
	tbodyTok.Map[1] = nextLine
	s.line = nextLine

	return true
}
