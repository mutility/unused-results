package csindirect

type struc struct {
	sink_, source_ int
}

func (s *struc) sink() *int  { return &s.sink_ }
func (s *struc) source() int { return s.source_ }

func use() {
	s := struc{3, 4}
	*s.sink() = s.source()
}
