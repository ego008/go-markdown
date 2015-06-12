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

func isHeaderLine(s string) bool {
	if s == "" {
		return false
	}

	st := 0
	n := 0
	for i := 0; i < len(s); i++ {
		b := s[i]
		switch st {
		case 0: // initial state
			switch b {
			case '|':
				st = 1
			case ':':
				st = 2
			case '-':
				st = 3
				n++
			case ' ':
				break
			default:
				return false
			}

		case 1: // |
			switch b {
			case ' ':
				break
			case ':':
				st = 2
			case '-':
				st = 3
				n++
			default:
				return false
			}

		case 2: // |:
			switch b {
			case ' ':
				break
			case '-':
				st = 3
				n++
			default:
				return false
			}

		case 3: // |:-
			switch b {
			case '-':
				break
			case ':':
				st = 4
			case '|':
				st = 5
			case ' ':
				st = 6
			default:
				return false
			}

		case 4: // |:---:
			switch b {
			case ' ':
				break
			case '|':
				st = 5
			default:
				return false
			}

		case 5: // |:---:|
			switch b {
			case ' ':
				break
			case ':':
				st = 2
			case '-':
				st = 3
				n++
			default:
				return false
			}

		case 6: // |:--- SPACE
			switch b {
			case ' ':
				break
			case ':':
				st = 4
			case '|':
				st = 5
			default:
				return false
			}
		}
	}

	return n >= 1
}
