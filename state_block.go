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

const (
	ptRoot = iota
	ptList
	ptBlockQuote
)

type stateBlock struct {
	stateCore

	bMarks     []int // offsets of the line beginnings
	eMarks     []int // offsets of the line endings
	tShift     []int // indents for each line
	blkIndent  int   // required block content indent (in a list etc.)
	line       int   // line index in the source string
	lineMax    int   // number of lines
	tight      bool  // loose or tight mode for lists
	parentType byte  // parent block type
	level      int
}

func (s *stateBlock) isLineEmpty(n int) bool {
	return s.bMarks[n]+s.tShift[n] >= s.eMarks[n]
}

func (s *stateBlock) skipEmptyLines(from int) int {
	for from < s.lineMax && s.isLineEmpty(from) {
		from++
	}
	return from
}

func (s *stateBlock) skipSpaces(pos int) int {
	src := s.src
	for pos < len(src) && src[pos] == ' ' {
		pos++
	}
	return pos
}

func (s *stateBlock) skipBytes(pos int, b byte) int {
	src := s.src
	for pos < len(src) && src[pos] == b {
		pos++
	}
	return pos
}

func (s *stateBlock) skipBytesBack(pos int, b byte, min int) int {
	for pos > min {
		pos--
		if s.src[pos] != b {
			return pos + 1
		}
	}
	return pos
}

func (s *stateBlock) lines(begin, end, indent int, keepLastLf bool) string {
	if begin == end {
		return ""
	}

	src := s.src

	if begin+1 == end {
		shift := s.tShift[begin]
		if shift < 0 {
			shift = 0
		} else if shift > indent {
			shift = indent
		}
		first := s.bMarks[begin] + shift

		last := s.eMarks[begin]
		if keepLastLf && last < len(src) {
			last++
		}

		return src[first:last]
	}

	size := 0
	var firstFirst int
	var previousLast int
	adjoin := true
	for line := begin; line < end; line++ {
		shift := s.tShift[line]
		if shift < 0 {
			shift = 0
		} else if shift > indent {
			shift = indent
		}
		first := s.bMarks[line] + shift
		last := s.eMarks[line]
		if line+1 < end || (keepLastLf && last < len(src)) {
			last++
		}
		size += last - first
		if line == begin {
			firstFirst = first
		} else if previousLast != first {
			adjoin = false
		}
		previousLast = last
	}

	if adjoin {
		return src[firstFirst:previousLast]
	}

	buf := make([]byte, size)
	i := 0
	for line := begin; line < end; line++ {
		shift := s.tShift[line]
		if shift < 0 {
			shift = 0
		} else if shift > indent {
			shift = indent
		}
		first := s.bMarks[line] + shift
		last := s.eMarks[line]
		if line+1 < end || (keepLastLf && last < len(src)) {
			last++
		}

		i += copy(buf[i:], src[first:last])
	}

	return string(buf)
}

func (s *stateBlock) pushToken(tok Token) {
	tok.SetLevel(s.level)
	s.tokens = append(s.tokens, tok)
}

func (s *stateBlock) pushOpeningToken(tok Token) {
	tok.SetLevel(s.level)
	s.level++
	s.tokens = append(s.tokens, tok)
}

func (s *stateBlock) pushClosingToken(tok Token) {
	s.level--
	tok.SetLevel(s.level)
	s.tokens = append(s.tokens, tok)
}
