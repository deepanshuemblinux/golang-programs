package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	file_name := os.Args[1]
	file, err := os.Open("something.txt")
	if err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf("The file %s does not exists\n", file_name)
			log.Fatal(err)
		} else if os.IsPermission(err) {
			err = fmt.Errorf("Permissions error while opening the file %s\n", file_name)
			log.Fatal(err)
		}

	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
