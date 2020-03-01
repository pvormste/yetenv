package yetenv

type EnvVariables map[string]string

func (v EnvVariables) Count() int {
	return len(v)
}
