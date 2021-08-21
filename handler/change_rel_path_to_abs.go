package handler

// ChangeRelPathToAbsStrategy 将相对改成绝对
type ChangeRelPathToAbsStrategy struct{}

func (s *ChangeRelPathToAbsStrategy) Adjust(h *BaseHandler) error {
	return nil
}

func (s *ChangeRelPathToAbsStrategy) Extra(h *BaseHandler) error {
	return nil
}
