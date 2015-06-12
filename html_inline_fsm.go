package markdown

import (
	"strings"

	"github.com/opennota/byteutil"
)

var (
	ws  [256]bool
	cs1 [256]bool
	cs2 [256]bool
	cs3 [256]bool
	cs4 [256]bool
)

func init() {
	for _, b := range "!/?abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		htmlSecond[b] = true
	}
	ws[' '], ws['\n'] = true, true
	for _, b := range "-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		cs1[b] = true
	}
	for _, b := range ":_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		cs2[b] = true
	}
	for _, b := range ":._-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		cs3[b] = true
	}
	for i := 0; i <= 32; i++ {
		cs4[byte(i)] = true
	}
	for _, b := range "\"'=<>`" {
		cs4[b] = true
	}
}

func matchHTML(s string) string {
	if len(s) < 3 || s[0] != '<' {
		return ""
	}

	st := 0

	for i := 1; i < len(s); {
		b := s[i]
		i++

		switch st {
		case 0: // initial state
			switch {
			case byteutil.IsLetter(b):
				st = 1
			case b == '/':
				st = 2
			case b == '!':
				st = 3
			case b == '?':
				st = 4
			default:
				return ""
			}

		case 1: // opening tag <DIV
			switch {
			case cs1[b]:
				break
			case ws[b]:
				st = 5
			case b == '/':
				st = 6
			case b == '>':
				return s[:i]
			default:
				return ""
			}

		case 2: // closing tag
			switch {
			case byteutil.IsLetter(b):
				st = 14
			default:
				return ""
			}

		case 3: // comment or declaration
			switch {
			case b == '-':
				st = 17
			case byteutil.IsUppercaseLetter(b):
				st = 18
			case b == '[':
				st = 19
			default:
				return ""
			}

		case 4: // processing instruction
			switch b {
			case '?':
				st = 16
			}

		case 5: // <DIV SPACE
			switch {
			case ws[b]:
				break
			case b == '/':
				st = 6
			case b == '>':
				return s[:i]
			case cs2[b]:
				st = 7
			default:
				return ""
			}

		case 6: // <BR/
			switch b {
			case '>':
				return s[:i]
			default:
				return ""
			}

		case 7: // <A H
			switch {
			case cs3[b]:
				break
			case b == '=':
				st = 9
			case ws[b]:
				st = 8
			case b == '/':
				st = 6
			case b == '>':
				return s[:i]
			default:
				return ""
			}

		case 8: // <A HREF SPACE
			switch {
			case ws[b]:
				break
			case b == '=':
				st = 9
			case b == '>':
				return s[:i]
			case cs2[b]:
				st = 7
			default:
				return ""
			}

		case 9: // <A HREF=
			switch {
			case ws[b]:
				break
			case b == '"':
				st = 10
			case b == '\'':
				st = 11
			case cs4[b]:
				return ""
			default:
				st = 12
			}

		case 10: // <A HREF="
			switch b {
			case '"':
				st = 13
			}

		case 11: // <A HREF='
			switch b {
			case '\'':
				st = 13
			}

		case 12: // <A HREF=H
			switch {
			case ws[b]:
				st = 5
			case b == '/':
				st = 6
			case b == '>':
				return s[:i]
			case cs4[b]:
				return ""
			default:
				st = 12
			}

		case 13: // <A HREF="http://google.com"
			switch {
			case ws[b]:
				st = 5
			case b == '/':
				st = 6
			case b == '>':
				return s[:i]
			default:
				return ""
			}

		case 14: // </I
			switch {
			case cs1[b]:
				break
			case ws[b]:
				st = 15
			case b == '>':
				return s[:i]
			default:
				return ""
			}

		case 15: // </IMG SPACE
			switch {
			case ws[b]:
				break
			case b == '>':
				return s[:i]
			default:
				return ""
			}

		case 16: // <?...?
			switch b {
			case '>':
				return s[:i]
			case '?':
				break
			default:
				st = 4
			}

		case 17: // <!-
			switch b {
			case '-':
				st = 20
			default:
				return ""
			}

		case 18: // <!D
			switch {
			case byteutil.IsUppercaseLetter(b):
				break
			case ws[b]:
				st = 23
			default:
				return ""
			}

		case 19: // <![
			switch {
			case strings.HasPrefix(s[i-1:], "CDATA["):
				i += 5
				st = 24
			default:
				return ""
			}

		case 20: // <!--
			switch b {
			case '-':
				st = 21
			case '>':
				return ""
			}

		case 21: // <!-- -
			switch b {
			case '-':
				st = 22
			case '>':
				return ""
			default:
				st = 20
			}

		case 22: // <!-- --
			switch b {
			case '>':
				return s[:i]
			default:
				return ""
			}

		case 23: // <!DOCTYPE SPACE
			switch b {
			case '>':
				return s[:i]
			}

		case 24: // <![CDATA[
			switch b {
			case ']':
				st = 25
			}

		case 25: // <![CDATA[ ... ]
			switch b {
			case ']':
				st = 26
			default:
				st = 24
			}

		case 26: // <![CDATA[ ... ]]
			switch b {
			case '>':
				return s[:i]
			default:
				st = 24
			}
		}
	}

	return ""
}
