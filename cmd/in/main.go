package main

import (
	"encoding/json"
	"os"
	"strconv"
	"time"
)

type Version struct {
	Ref string `json:"ref"`
}

func main() {
	version := Version{
		Ref: strconv.FormatInt(time.Now().Unix(), 10),
	}
	json.NewEncoder(os.Stdout).Encode(version)
}
