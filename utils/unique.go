package utils

type UniString struct {
	store map[string]struct{}
}

func NewUniString() *UniString {
	return &UniString{
		store: map[string]struct{}{},
	}
}

func (u *UniString) Add(key string) {
	if _, ok := u.store[key]; !ok {
		u.store[key] = struct{}{}
	}
}

func (u *UniString) Slice() []string {
	var slice []string
	for key, _ := range u.store {
		slice = append(slice, key)
	}
	return slice
}
