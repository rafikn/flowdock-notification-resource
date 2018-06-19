package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	os.Stdout.Write([]byte(fmt.Sprintf("{ \"version\" :{ \"ref\" :\"%d\"}}", time.Now().Unix())))
}
