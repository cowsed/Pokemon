package scripts

type Stack struct {
	arr   []int
	count int
}

func NewStack() *Stack {
	return &Stack{
		arr:   []int{},
		count: 0,
	}
}
func (s *Stack) Push(value int) {
	//Extend array if its too short
	if s.count >= len(s.arr) {
		s.arr = append(s.arr, 0)
	}
	s.arr[s.count] = value
	s.count++

}
func (s *Stack) Pop() int {
	if s.count <= 0 {
		return 0
	}
	s.count--
	return s.arr[s.count]
}
