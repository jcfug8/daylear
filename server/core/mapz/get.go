package mapz

import (
	"strconv"
	"strings"
)

// Get returns the value at a key in o. It follows the key chain
// to return the value of the last key in the chain.
func Get(o any, keys string) any {
	o, ok := o.(Map)
	if !ok {
		return nil
	}

	chain := strings.Split(keys, ".")
	var out = o

	for _, k := range chain {
		if o, ok := out.(List); ok {
			n, err := strconv.Atoi(k)
			if err != nil {
				return nil
			}

			out = o[n]
			continue
		}

		if o, ok := out.(Map); ok {
			out = o[k]
			continue
		}

		return nil
	}

	return out
}
