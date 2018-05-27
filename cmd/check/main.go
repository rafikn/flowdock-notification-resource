package main

import (
	"os"
	"fmt"
	"time"
)

func main() {
	os.Stdout.Write([]byte(fmt.Sprintf("{ \"version\" :{ \"ref\" :\"%d\"}}", time.Now().Unix())))
}
