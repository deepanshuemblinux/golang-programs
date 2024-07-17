package main

import "fmt"

func LongestCommonPrefix(inp []string) string {

	var lcp []rune
lab:
	for idx, val := range inp[0] {
		for _, str := range inp[1:] {
			if idx >= len(str) {
				break lab
			}
			if rune(str[idx]) != val {
				break lab
			}
		}

		lcp = append(lcp, val)

	}
	return string(lcp)
}
func main() {
	input := []string{"family", "fam", "famttish"}
	fmt.Println(LongestCommonPrefix(input))

}
