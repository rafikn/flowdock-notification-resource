package main

import (
	"fmt"
	"os"
)

func main() {
	os.Stdout.Write([]byte(fmt.Sprintf("{ \"version\" :{ \"ref\" :\"none\"}}")))
}
