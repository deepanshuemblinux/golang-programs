package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	file_name := "something.txt"
	file, err := os.Open(file_name)
	defer file.Close()
	if err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf("The file %s does not exist\n", file_name)
		} else if os.IsPermission(err) {
			err = fmt.Errorf("Permssions error while opening the file %s\n", file_name)
		}
		log.Fatal(err)
	}
	//Read file line by line using bufferred io scanner
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
