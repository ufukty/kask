package functions

import "fmt"

func Dict(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("there should be even number of arguments")
	}
	dict := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("argument %d: key is not string %T", i, key)
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}
