package utils

// Bool2int converts true=> 1, false => 0
func Bool2int(b bool) int {
	var i int
	if b {
		i = 1
	} else {
		i = 0
	}
	return i
}

// Empty2int coverst an empty string to 0, 1 otherwise
func Empty2int(s string) int {
	var i int
	if s == "" {
		i = 0
	} else {
		i = 1
	}
	return i
}
