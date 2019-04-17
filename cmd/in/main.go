package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"io/ioutil"
)

/*
	RESOURCE IN COMMAND

	Reads contents of version file from $VERSIONPATH and returns that version.

	Note: this is not quite coded to spec per the Concourse CI custom resource guide.
	(Spec: https://concourse-ci.org/implementing-resources.html)

	Much like the CHECK command (see ./cmd/check/main.go source) we are simply returning the
	latest version from the version file. This is a functional solution for now, but will ignore any
	input version params to the resource.

	FUTURE WORK: per spec, return the ref for the latest. If an input version param was provided,
	return that version if it exists instead of latest.
*/

type Version struct {
	Ref string `json:"ref"`
}

type Resource struct {
	Version *Version `json:"version"`
}

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
	versionpath := getVersionPathFromEnvironment()
	version := versionFromFile(versionpath)
	res := Resource{}
	res.Version = &version
	err := json.NewEncoder(os.Stdout).Encode(res)
	if err != nil { panic(err) }
}
