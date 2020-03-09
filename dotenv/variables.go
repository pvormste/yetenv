package dotenv

type Variables map[string]string

func (v Variables) Count() int {
	return len(v)
}
