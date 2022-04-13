package brainfuck

import (
	"errors"
	"fmt"
	"io"
)

func scanInput(reader io.Reader) (byte, error) {
	input := make([]byte, 1)
	_, err := reader.Read(input)
	if err != nil && !errors.Is(err, io.EOF) {
		return 0, fmt.Errorf("input while scanning input: %w", err)
	}
	return input[0], nil
}
