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
	"bytes"
	"net/url"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/opennota/html"
	"github.com/opennota/urlesc"
)

var mdpunct [256]bool

func init() {
	for _, b := range "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~" {
		mdpunct[b] = true
	}
}

func isMarkdownPunct(r rune) bool {
	if r > 0x7e {
		return false
	}
	return mdpunct[r]
}

func normalizeLink(rawurl string) string {
	parsed, err := url.Parse(rawurl)
	if err != nil {
		return ""
	}

	return urlesc.Escape(parsed)
}

func normalizeLinkText(text string) string {
	unescaped, _ := url.QueryUnescape(text)
	if unescaped == "" || !utf8.ValidString(unescaped) {
		return text
	}
	return unescaped
}

var badProtos = []string{"file", "javascript", "vbscript"}

var rGoodData = regexp.MustCompile(`^data:image/(gif|png|jpeg|webp);`)

func removeSpecial(s string) string {
	i := 0
	for i < len(s) && !(s[i] <= 0x20 || s[i] == 0x7f) {
		i++
	}
	if i >= len(s) {
		return s
	}
	buf := make([]byte, len(s))
	j := 0
	for i := 0; i < len(s); i++ {
		if !(s[i] <= 0x20 || s[i] == 0x7f) {
			buf[j] = s[i]
			j++
		}
	}
	return string(buf[:j])
}

func validateLink(url string) bool {
	str := html.ReplaceEntities(url)
	str = strings.TrimSpace(str)
	str = strings.ToLower(str)

	if strings.IndexByte(str, ':') >= 0 {
		proto := strings.SplitN(str, ":", 2)[0]
		proto = removeSpecial(proto)
		for _, p := range badProtos {
			if proto == p {
				return false
			}
		}
		if proto == "data" && !rGoodData.MatchString(str) {
			return false
		}
	}

	return true
}

func unescapeAll(s string) string {
	anyChanges := false
	i := 0
	for i < len(s)-1 {
		b := s[i]
		if b == '\\' {
			if mdpunct[s[i+1]] {
				anyChanges = true
				break
			}
		} else if b == '&' {
			if _, n := html.ParseEntity(s[i:]); n > 0 {
				anyChanges = true
				break
			}
		}
		i++
	}

	if !anyChanges {
		return s
	}

	buf := make([]byte, len(s))
	copy(buf[:i], s)
	j := i
	for i < len(s) {
		b := s[i]
		if b == '\\' {
			if i+1 < len(s) {
				b = s[i+1]
				if mdpunct[b] {
					buf[j] = b
					j++
				} else {
					buf[j] = '\\'
					j++
					buf[j] = b
					j++
				}
				i += 2
				continue
			}
		} else if b == '&' {
			if e, n := html.ParseEntity(s[i:]); n > 0 {
				if len(e) > n && len(buf) == len(s) {
					newBuf := make([]byte, cap(buf)*2)
					copy(newBuf[:j], buf)
					buf = newBuf
				}
				j += copy(buf[j:], e)
				i += n
				continue
			}
		}
		buf[j] = b
		j++
		i++
	}

	return string(buf[:j])
}

func normalizeInlineCode(s string) string {
	if s == "" {
		return ""
	}

	byteFuckery := false
	i := 0
	for i < len(s)-1 {
		b := s[i]
		if b == '\n' {
			byteFuckery = true
			break
		}
		if b == ' ' {
			i++
			b = s[i]
			if b == ' ' || b == '\n' {
				i--
				byteFuckery = true
				break
			}
		}
		i++
	}

	if !byteFuckery {
		return strings.TrimSpace(s)
	}

	buf := make([]byte, len(s))
	copy(buf[:i], s)
	buf[i] = ' '
	i++
	j := i
	lastSpace := true
	for i < len(s) {
		b := s[i]
		switch b {
		case ' ', '\n':
			if lastSpace {
				break
			}

			buf[j] = ' '
			lastSpace = true
			j++
		default:
			buf[j] = b
			lastSpace = false
			j++
		}

		i++
	}

	return string(bytes.TrimSpace(buf[:j]))
}

func normalizeReference(s string) string {
	var buf bytes.Buffer
	lastSpace := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !lastSpace {
				buf.WriteByte(' ')
				lastSpace = true
			}
			continue
		}

		buf.WriteRune(unicode.To(unicode.LowerCase, r))
		lastSpace = false
	}

	return string(bytes.TrimSpace(buf.Bytes()))
}

func skipws(s string, pos, max int) int {
	for pos < max && ws[s[pos]] {
		pos++
	}
	return pos
}
