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

type Versions []Version

func main() {
	versions := Versions{}
	version := Version{
		Ref: strconv.FormatInt(time.Now().Unix(), 10),
	}
	versions = append(versions, version)
	json.NewEncoder(os.Stdout).Encode(versions)
}
