package yetenv

type dotenvVariables map[string]string

func (v dotenvVariables) count() int {
	return len(v)
}
