package main

import (
	"fmt"
)

func isVowel(r rune) bool {
	switch r {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	}
	return false
}
func reverseVowels(str string) string {
	runes := []rune(str)
	b := 0
	e := len(runes) - 1
	for {
		for {
			if isVowel(runes[b]) {
				break
			}
			b++
		}
		for {
			if isVowel(runes[e]) {
				break
			}
			e--
		}
		if b >= e {
			break
		}
		temp := runes[b]
		runes[b] = runes[e]
		runes[e] = temp
		b++
		e--
	}
	return string(runes)
}
func main() {
	str := "simple"
	fmt.Println("Given string is ", str)
	str = reverseVowels(str)
	fmt.Println("Output is ", str)
	str = "complex"
	fmt.Println("Given string is ", str)
	str = reverseVowels(str)
	fmt.Println("Output is ", str)

}
