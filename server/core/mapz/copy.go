package mapz

// CopyMap copies a map and returns the new map.
func CopyMap(in Map) Map {
	if in == nil {
		return nil
	}

	out := make(Map, len(in))

	for k, v := range in {
		switch v := v.(type) {
		case Map:
			out[k] = CopyMap(v)
		case List:
			out[k] = CopyList(v)
		default:
			out[k] = v
		}
	}

	return out
}

// CopyList copies a list and returns the new list.
func CopyList(in List) List {
	if in == nil {
		return nil
	}

	out := make(List, len(in))

	for i, v := range in {
		switch v := v.(type) {
		case Map:
			out[i] = CopyMap(v)
		case List:
			out[i] = CopyList(v)
		default:
			out[i] = v
		}
	}

	return out
}
