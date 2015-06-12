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

import "bytes"

type stateInline struct {
	stateCore

	pos          int
	posMax       int
	level        int
	pending      bytes.Buffer
	pendingLevel int

	cache map[int]int
}

func (s *stateInline) pushToken(tok Token) {
	if s.pending.Len() > 0 {
		s.pushPending()
	}
	tok.SetLevel(s.level)
	s.pendingLevel = s.level
	s.tokens = append(s.tokens, tok)
}

func (s *stateInline) pushOpeningToken(tok Token) {
	if s.pending.Len() > 0 {
		s.pushPending()
	}
	tok.SetLevel(s.level)
	s.level++
	s.pendingLevel = s.level
	s.tokens = append(s.tokens, tok)
}

func (s *stateInline) pushClosingToken(tok Token) {
	if s.pending.Len() > 0 {
		s.pushPending()
	}
	s.level--
	tok.SetLevel(s.level)
	s.pendingLevel = s.level
	s.tokens = append(s.tokens, tok)
}

func (s *stateInline) pushPending() {
	s.tokens = append(s.tokens, &Text{
		Content: s.pending.String(),
		Lvl:     s.pendingLevel,
	})
	s.pending.Reset()
}
