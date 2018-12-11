package controllers

type loginSerializer struct {
	LoginCode string `valid:"MinSize(1); MaxSize(128)"`
}

func (s *loginSerializer) Validate() (map[string]string, error) {
	return validateStruct(s)
}
