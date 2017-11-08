package model

type NameValue struct {
	Name string `json:"name"`
	Value string `json:"value"`
}

func NewNameValue(n, v string) NameValue {
	return NameValue{
		Name: n,
		Value: v,
	}
}

func NewNameValues(m map[string]string) []NameValue {
	n := make([]NameValue, len(m))

	i := 0
	for k, v := range m {
		n[i] = NewNameValue(k,v)
		i++
	}
	return n
}