package mapz

// Unchain returns a receiver if the given value is a chain
// otherwise it returns the given value.
func Unchain(value any) any {
	switch v := value.(type) {
	case *ChainStruct:
		return v.receiver
	default:
		return value
	}
}

// Chain returns a chain with the given receiver.
func Chain(receiver any) *ChainStruct {
	return &ChainStruct{receiver: receiver}
}

// ChainStruct represents a chain of objects.
type ChainStruct struct {
	receiver any
}

// Set sets a value at the key point.
func (chain *ChainStruct) Set(keys string, value any) *ChainStruct {
	Set(chain.receiver, keys, value)
	return chain
}

// Append appends a value to a chain at the key point.
func (chain *ChainStruct) Append(keys string, value any) *ChainStruct {
	Append(chain.receiver, keys, value)
	return chain
}
