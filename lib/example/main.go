package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(getenv("payload"))
}


func getenv(s string) string {
	return os.Getenv(s)
}
