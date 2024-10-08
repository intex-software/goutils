package iterators

import (
	"encoding/json"
	"iter"
)

func JsonlReader[T any](input iter.Seq2[string, error]) iter.Seq2[T, error] {
	next, stop := iter.Pull2(input)

	return func(yield func(T, error) bool) {
		defer stop()

		for {
			var data T

			line, err, ok := next()
			if err != nil {
				yield(data, err)
				break
			} else if !ok {
				break
			}

			if len(line) == 0 {
				yield(data, ErrEmptyLine)
				continue
			}

			if err = json.Unmarshal([]byte(line), &data); err != nil {
				yield(data, err)
				break
			} else if ok = yield(data, err); !ok {
				break
			}
		}
	}
}
