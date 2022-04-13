package brainfuck

type intStack []int

func (s *intStack) isEmpty() bool {
	return len(*s) == 0
}

func (s *intStack) push(item int) {
	*s = append(*s, item)
}

func (s *intStack) pop() (int, bool) {
	if s.isEmpty() {
		return 0, false
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return element, true
	}
}
