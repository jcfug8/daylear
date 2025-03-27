package mapz

import "strings"

// Append appends a value to in. It uses the keys to determine where to
// append the value. For example, Append(mapOfMaps, "foo.bar.baz", mapz.Map{"two": 2})
// will append mapz.Map{"two": 2} to mapOfMaps inside the baz list. It creates the
// key and adds the value if the keys do not exist.
func Append(in any, keys string, value any) {
	value = Unchain(value)

	chain := strings.Split(keys, ".")
	tail := len(chain) - 1

	x := in

	for i, keys := range chain {
		switch v := x.(type) {
		case Map:
			if i == tail {
				target := v[keys]

				switch target := target.(type) {
				case List:
					v[keys] = append(target, value)
				case nil:
					v[keys] = value
				default:
					v[keys] = List{target, value}
				}

				return
			}

			if _, ok := v[keys]; !ok {
				v[keys] = Map{}
			}

			x = v[keys]
		default:
			return
		}
	}
}
