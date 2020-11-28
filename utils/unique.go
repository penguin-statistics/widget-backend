package utils

// UniString contains only unique strings which uses a underlying map which provides a constant time complexity
type UniString struct {
	// store uses its key to store the actual string, its value does not have any meaning. Notice that we used a struct{} in order to not let it allocate any more space
	store map[string]struct{}
}

// NewUniString creates a new UniString instance that contains only unique strings
func NewUniString() *UniString {
	return &UniString{
		store: map[string]struct{}{},
	}
}

// Add adds the string to the UniString instance only when there is no such string
func (u *UniString) Add(key string) {
	if _, ok := u.store[key]; !ok {
		u.store[key] = struct{}{}
	}
}

// Slice returns a string slice from UniString instance
func (u *UniString) Slice() []string {
	var slice []string
	for key := range u.store {
		slice = append(slice, key)
	}
	return slice
}
