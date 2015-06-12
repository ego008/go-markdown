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
	"bufio"
	"io"
)

type writer interface {
	Write([]byte) (int, error)
	WriteByte(byte) error
	WriteString(string) (int, error)
	Flush() error
}

type monadicWriter struct {
	writer
	err error
}

func newMonadicWriter(w io.Writer) *monadicWriter {
	if w, ok := w.(writer); ok {
		return &monadicWriter{writer: w}
	}
	return &monadicWriter{writer: bufio.NewWriter(w)}
}

func (w *monadicWriter) Write(p []byte) (n int, err error) {
	if w.err != nil {
		return
	}

	n, err = w.writer.Write(p)
	w.err = err
	return
}

func (w *monadicWriter) WriteByte(b byte) (err error) {
	if w.err != nil {
		return
	}

	err = w.writer.WriteByte(b)
	w.err = err
	return
}

func (w *monadicWriter) WriteString(s string) (n int, err error) {
	if w.err != nil {
		return
	}

	n, err = w.writer.WriteString(s)
	w.err = err
	return
}

func (w *monadicWriter) Flush() (err error) {
	if w.err != nil {
		return
	}

	err = w.writer.Flush()
	w.err = err
	return
}
