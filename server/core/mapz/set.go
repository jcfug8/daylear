package mapz

import "strings"

// Set sets the value of a key in a map and returns the map.
func Set(in any, keys string, value any) any {
	if in == nil {
		in = Map{}
	}

	value = Unchain(value)

	chain := strings.Split(keys, ".")
	tail := len(chain) - 1

	x := in

	for i, keys := range chain {
		if v, ok := x.(Map); ok {
			if i == tail {
				v[keys] = value
				break
			}

			if _, ok := v[keys]; !ok {
				v[keys] = Map{}
			}

			x = v[keys]
		}
	}

	return in
}
