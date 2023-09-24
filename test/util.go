package test

func ContainsItem(item string, arr []string) bool {
	for _, i := range arr {
		if i == item {
			return true
		}
	}

	return false
}

func MapEqual(a map[string]interface{}, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}

	for key, valA := range a {
		valB, ok := b[key]
		if !ok || valA != valB {
			return false
		}
	}

	return true
}
