package main

import (
	"log"
	"os"

	"github.com/deepanshuemblinux/go-json-encoder/encoder"
)

type MyStruct struct {
	Name   string `myjson:"name"`
	Age    int    `myjson:"age"`
	Skills map[string]float64
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go-json-encoder <file_name.json>")
	}
	file, err := os.OpenFile(os.Args[1], os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	// m := make(map[string]int)
	// m["deepanshu"] = 90
	// m["bubu"] = 100
	// m["ram"] = 99

	// m := make(map[string]float64)
	// m["deepanshu"] = 90.5
	// m["bubu"] = 100
	// m["ram"] = 99.7
	// m := make(map[string]string)
	// m["deepanshu"] = "sharma"
	// m["bubu"] = "babbar-sher"
	// m["ram"] = "lalla"
	// m := make(map[int]string)
	// m[20000] = "Bubu"
	// m[40000] = "Ram"
	// m[50000] = "deepanshu"
	// m1 := make(map[string]map[int]string)
	// m1["1"] = m

	// m := make(map[string][]int)
	// m["deepanshu"] = []int{1, 2, 3, 4}
	// m["mukul"] = []int{5, 6, 7, 8}
	// m["Bubu"] = []int{9, 10, 11, 12}
	// m := make(map[string][]string)
	// m["deepanshu"] = []string{"a", "b", "c"}
	// m["mukul"] = []string{"d", "e", "f"}
	// m["Bubu"] = []string{"g", "h", "i"}
	s := MyStruct{
		Name:   "Deepanshu",
		Age:    33,
		Skills: make(map[string]float64),
	}
	s.Skills["Reading"] = 8.5
	s.Skills["Writing"] = 9.5
	s.Skills["Listening"] = 7.5
	//json.NewEncoder(file).Encode(s)
	encoder.NewEncoder(file).Encode(s)
}
