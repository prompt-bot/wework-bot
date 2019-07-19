package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(getEnv("payload"))
}


func getEnv(s string) string {
	return os.Getenv(s)
}
