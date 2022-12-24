package util

import "strings"

func ParseWithoutCompression(number string) string {
	area := strings.Split(number, ".")
	var area1 []string
	chars := []rune{}
	if len(area) > 1 {
		chars = []rune(area[1])
	}
	for i := 0; i < len(chars); i++ {
		if i < 4 {
			char := string(chars[i])
			area1 = append(area1, char)
		}
	}
	str1 := strings.Join(area1, "")
	var result string
	if len(area) > 1 {
		result = area[0] + "." + str1
	} else {
		result = area[0]
	}
	return result
}
