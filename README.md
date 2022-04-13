# BrainFuck

## Installation

```bash
got get github.com/Amirhossein2000/brainfuck
```

## Example

```go
package main

import (
	"fmt"
	"github.com/Amirhossein2000/brainfuck"
	"io"
	"os"
)

func main() {
	interpreter := brainfuck.NewInterpreter(128, os.Stdin, os.Stdout)
	interpreter.AddOperation('*', func(codePointer, memPointer *int, mem, code []byte, reader io.Reader, writer io.Writer) {
		mem[*memPointer] *= 2
	})

	err := interpreter.Run()
	if err != nil {
		fmt.Println("interpreter err:", err.Error())
	}
}
```
```go
    // if you want to have and command-line execute for more than once you can put the for
    for {
        err := interpreter.Run()
        if err != nil {
        fmt.Println("interpreter err:", err.Error())
    }
}
```
```bash
~ go run main.go
> ++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++.
> Hello World!
```

## Portability issues

### Cell size

In the classic distribution, the cells are of 8-bit size (cells are bytes), and this is still the most common size,
therefore I used the 8-bit in this implementation. Also, it's much easier to use the byte array.

### Array size

It's configurable in this implementation.

```go
    // in this example the memory size is 128
    interpreter := brainfuck.NewInterpreter(128, os.Stdin, os.Stdout)
```

### End-of-line code

\n is the most common one, therefore I used it in this implementation. have any improvement in your mind please feel free to open an issue.