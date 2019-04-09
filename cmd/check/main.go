package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"io/ioutil"
)

type Version struct {
	Ref string `json:"ref"`
}

type Versions []Version

func getVersionPathFromEnvironment() string {
	VERSIONPATH, bExists := os.LookupEnv("VERSIONPATH")
	if !bExists || VERSIONPATH == "" {
		panic("$VERSIONPATH is empty or unset. Cannot retrieve resource verison")
	}

	abspath, err := filepath.Abs(filepath.Join(filepath.Dir(VERSIONPATH), "version"))
	if err != nil {
		panic(err)
	}
	return abspath
}

func versionFromFile(vpath string) (Version) {
	version, err := ioutil.ReadFile(vpath)
	if err != nil { panic (err) }
	v := Version{}
	v.Ref = string(version)
	return v
}

func main() {
	versionlist := Versions{}
	versionpath := getVersionPathFromEnvironment()
	version := versionFromFile(versionpath)
	versionlist = append(versionlist, version)
	err := json.NewEncoder(os.Stdout).Encode(versionlist)
	if err != nil { panic(err) }
}
