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

// Package markdown provides CommonMark-compliant markdown parser and renderer.
package markdown

import (
	"bytes"
	"io"
)

type Markdown struct {
	options
	block         block
	inline        inline
	renderOptions RenderOptions
}

type RenderOptions struct {
	XHTML      bool   // render as XHTML instead of HTML
	Breaks     bool   // convert \n in paragraphs into <br>
	LangPrefix string // CSS language class prefix for fenced blocks
	Nofollow   bool   // add rel="nofollow" to the links
}

type options struct {
	HTML        bool    // allow raw HTML in the markup
	Tables      bool    // GFM tables
	Linkify     bool    // autoconvert URL-like text to links
	Typographer bool    // enable some typographic replacements
	Quotes      [4]rune // double/single quotes replacement pairs
	MaxNesting  int     // maximum nesting level
}

type environment struct {
	References map[string]map[string]string
}

type coreRule func(*stateCore)

func New(opts ...option) *Markdown {
	m := &Markdown{
		options: options{
			Tables:      true,
			Linkify:     true,
			Typographer: true,
			Quotes:      [4]rune{'“', '”', '‘', '’'},
			MaxNesting:  20,
		},
		renderOptions: RenderOptions{LangPrefix: "language-"},
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *Markdown) Parse(src []byte) []Token {
	if len(src) == 0 {
		return nil
	}

	s := &stateCore{
		md:  m,
		env: &environment{},
	}
	s.tokens = m.block.parse(src, m, s.env)

	for _, r := range []coreRule{
		ruleInline,
		ruleLinkify,
		ruleReplacements,
		ruleSmartQuotes,
	} {
		r(s)
	}
	return s.tokens
}

func (m *Markdown) Render(w io.Writer, src []byte) error {
	if len(src) == 0 {
		return nil
	}

	return NewRenderer(w).Render(m.Parse(src), m.renderOptions)
}

func (m *Markdown) RenderTokens(w io.Writer, tokens []Token) error {
	if len(tokens) == 0 {
		return nil
	}

	return NewRenderer(w).Render(tokens, m.renderOptions)
}

func (m *Markdown) RenderToString(src []byte) string {
	if len(src) == 0 {
		return ""
	}

	var buf bytes.Buffer
	NewRenderer(&buf).Render(m.Parse(src), m.renderOptions)
	return buf.String()
}

func (m *Markdown) RenderTokensToString(tokens []Token) string {
	if len(tokens) == 0 {
		return ""
	}

	var buf bytes.Buffer
	NewRenderer(&buf).Render(tokens, m.renderOptions)
	return buf.String()
}
