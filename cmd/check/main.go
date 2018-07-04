package main

import (
	"encoding/json"
	"os"
)

type Version struct {
	Ref string `json:"ref"`
}

type Versions []Version

func main() {
	versions := Versions{}
	json.NewEncoder(os.Stdout).Encode(versions)
}
