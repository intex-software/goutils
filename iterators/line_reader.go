package iterators

import (
	"bufio"
	"errors"
	"io"
	"iter"
)

func LineReader(r io.Reader) iter.Seq2[string, error] {
	br := bufio.NewReader(r)

	return func(yield func(string, error) bool) {
		for {
			line, err := br.ReadString('\n')
			eof := errors.Is(err, io.EOF)

			if err != nil && !eof {
				yield(line, err)
				break
			}

			if !yield(line, nil) {
				break
			}

			if eof {
				break
			}
		}
	}
}
