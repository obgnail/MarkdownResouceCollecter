package handler

// ChangeRelPathToAbsStrategy 将相对改成绝对
type ChangeRelPathToAbsStrategy struct{}

func (s *ChangeRelPathToAbsStrategy) BeforeRewrite(h *BaseHandler) error {
	return nil
}

func (s *ChangeRelPathToAbsStrategy) AfterRewrite(h *BaseHandler) error {
	return nil
}
