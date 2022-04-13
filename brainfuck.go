package brainfuck

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// If you want to add any custom operation you have to use this func
type OperationHandler func(codePointer, memPointer *int, mem, code []byte, reader io.Reader, writer io.Writer)

type Interpreter struct {
	memoryLen            int
	reader               io.Reader
	writer               io.Writer
	customOperations     map[byte]OperationHandler
	notAllowedOperations map[byte]bool
}

func NewInterpreter(memoryLen int, reader io.Reader, writer io.Writer) *Interpreter {
	if writer == nil {
		writer = os.Stdout
	}
	if reader == nil {
		reader = os.Stdin
	}
	if memoryLen <= 0 {
		memoryLen = 128
	}
	return &Interpreter{
		memoryLen:            memoryLen,
		reader:               reader,
		writer:               writer,
		customOperations:     make(map[byte]OperationHandler),
		notAllowedOperations: make(map[byte]bool),
	}
}

func (interpreter *Interpreter) Run() error {
	mem := make([]byte, interpreter.memoryLen)
	reader := interpreter.reader
	writer := interpreter.writer
	code := []byte{}
	memPointer := 0
	codePointer := 0
	loopTable := make(map[int]int) // map beginning to end of the loop
	stack := intStack{}
	unfinishedLoops := intStack{} // start index of the loops which they haven't scanned completely yet

	for {
		// Scan new input when the code pointer is out of any loop
		if codePointer == len(code) {
			input, err := scanInput(reader)
			if err != nil {
				return err
			}
			if _, ok := interpreter.notAllowedOperations[input]; ok {
				return fmt.Errorf("operation %s is not allowed", string(input))
			}
			code = append(code, input)
		}

		/* When we have to jump to the end of the loop but we haven't scanned the end of the loop yet, therefore
		   we need to ignore every input until we scan the end of the loop */
		if !unfinishedLoops.isEmpty() && code[codePointer] != '[' && code[codePointer] != ']' {
			codePointer++
			continue
		}

		switch code[codePointer] {
		case '\n':
			return nil
		case '>':
			memPointer++
			if memPointer > interpreter.memoryLen {
				return errors.New("memory pointer can't be greater than memory length")
			}
		case '<':
			memPointer--
			if memPointer < 0 {
				return errors.New("memory pointer can't be less than 0")
			}
		case '-':
			mem[memPointer]--
		case '+':
			mem[memPointer]++
		case '.':
			_, err := fmt.Fprintf(writer, string(mem[memPointer]))
			if err != nil {
				return fmt.Errorf("input while writing error: %s", err.Error())
			}
		case ',':
			input, err := scanInput(reader)
			if err != nil {
				return err
			}
			mem[memPointer] = input
		case '[':
			stack.push(codePointer)
			if mem[memPointer] == 0 {
				endLoopIndex, ok := loopTable[codePointer]
				if ok {
					codePointer = endLoopIndex - 1
				} else {
					unfinishedLoops.push(codePointer)
				}
			}
		case ']':
			beginningOfTheLoop, ok := stack.pop()
			if !ok {
				return errors.New("operation ] must be entered after a [")
			}
			loopTable[beginningOfTheLoop] = codePointer
			if mem[memPointer] != 0 {
				codePointer = beginningOfTheLoop - 1
				/* -1 exists because there is a codePointer++ in the end of the for-loop and if we break here the
				   next condition will be ignored */
			}
			if !unfinishedLoops.isEmpty() {
				unfinishedLoops.pop()
			}
		default:
			handler, ok := interpreter.customOperations[code[codePointer]]
			if ok {
				handler(&codePointer, &memPointer, mem, code, reader, writer)
			}
		}
		codePointer++
	}
}

func (interpreter *Interpreter) AddOperation(symbol byte, handler OperationHandler) {
	interpreter.customOperations[symbol] = handler
}

func (interpreter *Interpreter) DeleteOperations(symbols ...byte) {
	for _, symbol := range symbols {
		if symbol == '[' || symbol == ']' {
			// A loop with only one operation doesn't make any sense
			interpreter.notAllowedOperations['['] = true
			interpreter.notAllowedOperations[']'] = true
		} else {
			interpreter.notAllowedOperations[symbol] = true
		}
	}
}
