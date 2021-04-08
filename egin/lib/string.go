package lib

import "strings"

type String struct {
	Str string
}

func (s *String) TrimLeft(cutset string) *String {
	s.Str = strings.TrimLeft(s.Str, cutset)
	return s
}

func (s *String) TrimRight(cutset string) *String {
	s.Str = strings.TrimRight(s.Str, cutset)
	return s
}

func (s *String) Done() string {
	return s.Str
}
