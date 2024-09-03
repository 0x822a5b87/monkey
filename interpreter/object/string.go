package object

type StringObj struct {
	Value string
	Env   *Environment
}

func (s *StringObj) Type() ObjType {
	return ObjString
}

func (s *StringObj) Inspect() string {
	return s.Value
}

func (s *StringObj) Add(object Object) Object {
	s.Value = s.Value + object.Inspect()
	return s
}

func (s *StringObj) Index(o Object) Object {
	other, ok := o.(*Integer)
	if !ok {
		return NativeNull
	}
	if other.Value >= int64(len(s.Value)) || other.Value < 0 {
		return NativeNull
	}
	ch := s.Value[other.Value]
	return &StringObj{Value: string(rune(ch))}
}

func (s *StringObj) First() Object {
	if s.Len().Value == 0 {
		return NativeNull
	}
	return s.Index(&Integer{Value: 0})
}

func (s *StringObj) Last() Object {
	if s.Len().Value == 0 {
		return NativeNull
	}
	return s.Index(&Integer{Value: s.Len().Value - 1})
}

func (s *StringObj) Len() Integer {
	return Integer{Value: int64(len(s.Value))}
}
