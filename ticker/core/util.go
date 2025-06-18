package core

import "strconv"

func FormatUintWithCommas(n uint64) string {
	s := strconv.FormatUint(n, 10)
	length := len(s)
	if length <= 3 {
		return s
	}

	var result string
	remainder := length % 3

	// Add the first group (which may be less than 3 digits)
	if remainder > 0 {
		result = s[:remainder]
		if length > remainder {
			result += ","
		}
	}

	// Add the rest of the groups
	for i := remainder; i < length; i += 3 {
		result += s[i : i+3]
		if i+3 < length {
			result += ","
		}
	}

	return result
}
