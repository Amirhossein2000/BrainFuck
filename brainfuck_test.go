package brainfuck

import (
	"bytes"
	"io"
	"testing"
)

func TestStartAndAddOperation(t *testing.T) {
	reader := bytes.NewBuffer([]byte{})
	writer := bytes.NewBuffer([]byte{})
	interpreter := NewInterpreter(128, reader, writer)
	interpreter.AddOperation('*', func(codePointer, memPointer *int, mem, code []byte, reader io.Reader, writer io.Writer) {
		mem[*memPointer] *= 2
	})

	tests := []struct {
		input  string
		output string
	}{
		{
			input:  "++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++.\n",
			output: "Hello World!\n",
		},
		{
			input:  ">++++++++[<+++++++++>-]<.>++++[<+++++++>-]<+.+++++++..+++.>>++++++[<+++++++>-]<++.------------.>++++++[<+++++++++>-]<+.<.+++.------.--------.>>>++++[<++++++++>-]<+.\n",
			output: "Hello, World!",
		},
		{
			input:  "++++++++++[->++++++++++>+++++++++++>+++<<<]>+.>.++++++.<.>--.>++.<----.+++++++.--------.<---.+++.[-]>+++++.++[--<+>]<.[-]>>.[-]<<<\n",
			output: "enter number: ",
		},
		{
			input:  "++++[->++++++++<]>+.-[->++<]>-.\n",
			output: "!?",
		},
		{
			input:  "+******.\n",
			output: "@", // Code 64
		},
		{
			input:  "+******-.\n",
			output: "?", // Code 63
		},
		{
			input:  "+*****+.\n",
			output: "!", // Code 33
		},
		{
			input:  ",2++++++.\n",
			output: "8", // Code 15
		},
		{
			input:  ",a+++.,@.\n",
			output: "d@", // Code 15
		},
	}

	for _, tt := range tests {
		t.Run(tt.output, func(t *testing.T) {
			reader.Reset()
			writer.Reset()
			reader.WriteString(tt.input)
			err := interpreter.Run()
			if err != nil {
				t.Errorf("interpreter err: %s", err.Error())
			}
			if writer.String() != tt.output {
				t.Errorf("expected %q but received %q", tt.output, writer.String())
			}
		})
	}
}

func TestInterpreter_DeleteOperation(t *testing.T) {
	reader := bytes.NewBuffer([]byte{})
	writer := bytes.NewBuffer([]byte{})
	interpreter := NewInterpreter(128, reader, writer)
	interpreter.DeleteOperations('[', '-')

	tests := []struct {
		interpreter *Interpreter
		input       string
		error       string
	}{
		{
			interpreter: interpreter,
			input:       "[-]+.\n",
			error:       "operation [ is not allowed",
		},
		{
			interpreter: interpreter,
			input:       "----.\n",
			error:       "operation - is not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.error, func(t *testing.T) {
			reader.Reset()
			writer.Reset()
			reader.WriteString(tt.input)
			err := tt.interpreter.Run()
			if err == nil {
				t.Errorf("mustn't")
			} else if err.Error() != tt.error {
				t.Errorf("expected err message %q but received %q", tt.error, err.Error())
			}
		})
	}
}

func TestInterpreter_Errors(t *testing.T) {
	reader := bytes.NewBuffer([]byte{})
	writer := bytes.NewBuffer([]byte{})
	interpreter := NewInterpreter(2, reader, writer)

	tests := []struct {
		interpreter *Interpreter
		input       string
		error       string
	}{
		{
			interpreter: interpreter,
			input:       ">>>\n",
			error:       "memory pointer can't be greater than memory length",
		},
		{
			interpreter: interpreter,
			input:       "<\n",
			error:       "memory pointer can't be less than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.error, func(t *testing.T) {
			reader.Reset()
			writer.Reset()
			reader.WriteString(tt.input)
			err := tt.interpreter.Run()
			if err == nil {
				t.Errorf("error mustn't be nil")
			} else if err.Error() != tt.error {
				t.Errorf("expected err message %q but received %q", tt.error, err.Error())
			}
		})
	}
}
