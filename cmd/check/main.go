package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"io/ioutil"
)

/*
	RESOURCE CHECK COMMAND

	Reads contents of version file from $VERSIONPATH and returns that version.

	Note: this is not quite coded to spec per the Concourse CI custom resource guide.
	The check command should return an array of versions after retrieving them from
	the target repository. (e.g. git hashes, ecr hashes, ... ).
	If an input version param is given, that version should be returned in the array if it exists.
	(Spec: https://concourse-ci.org/implementing-resources.html)

	INSTEAD: We are returning a single semver from the version file as the only element
	in the versions array. The assumption being that the version file always indicates the latest version.
	We are ignoring any version params to the resource definition beyond the repository and tag params.

	One possible advantage is that we can refer to the version file using a semver resource in a pipeline.
	That will let us build this image and automatically increment the version in the pipeline rather than
	populating the version file manually.

	FUTURE WORK: this could use the aws sdk to retrive a list of all digests from ECR and return those instead.
	That could then respect the version param to check if a given hash is in the list.
*/

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
