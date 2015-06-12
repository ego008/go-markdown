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

type option func(m *Markdown)

func HTML(b bool) option {
	return func(m *Markdown) {
		m.HTML = b
	}
}

func Linkify(b bool) option {
	return func(m *Markdown) {
		m.Linkify = b
	}
}

func Typographer(b bool) option {
	return func(m *Markdown) {
		m.Typographer = b
	}
}

func Quotes(s string) option {
	return func(m *Markdown) {
		for i, r := range s {
			m.Quotes[i] = r
		}
	}
}

func MaxNesting(n int) option {
	return func(m *Markdown) {
		m.MaxNesting = n
	}
}

func XHTMLOutput(b bool) option {
	return func(m *Markdown) {
		m.renderOptions.XHTML = b
	}
}

func Breaks(b bool) option {
	return func(m *Markdown) {
		m.renderOptions.Breaks = b
	}
}

func LangPrefix(p string) option {
	return func(m *Markdown) {
		m.renderOptions.LangPrefix = p
	}
}

func Nofollow(b bool) option {
	return func(m *Markdown) {
		m.renderOptions.Nofollow = b
	}
}

func Tables(b bool) option {
	return func(m *Markdown) {
		m.Tables = b
	}
}
