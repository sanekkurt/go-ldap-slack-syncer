package utils

func splitStringIntoBoundedStrings(in string, maxLengthOneString uint) []string {
	var (
		l      = uint(len(in))
		result []string
		from   = uint(0)
		to     = maxLengthOneString
	)

	if maxLengthOneString == 0 {
		return nil
	}

	for {
		if l > to {
			result = append(result, in[from:to])
			from = to
			to += maxLengthOneString

			continue
		}
		result = append(result, in[from:l])

		break
	}

	return result
}

func SplitStr(in []string, maxLengthOneString uint) []string {
	var (
		result []string

		oneStr string

		length = len(in)
	)

	for index, str := range in {
		if uint(len(str)) > maxLengthOneString {
			for _, partStr := range splitStringIntoBoundedStrings(str, maxLengthOneString) {
				result = append(result, partStr)
			}

		}

		if uint(len(oneStr)+len(str)) < maxLengthOneString {
			oneStr += str

			if index == length-1 {
				result = append(result, oneStr)
			}
		} else {
			result = append(result, oneStr)

			oneStr = str
		}
	}

	return result
}
