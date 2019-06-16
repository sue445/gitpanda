package testutil

import (
	"io/ioutil"
	"path"
)

func ReadTestData(filename string) string {
	buf, err := ioutil.ReadFile(path.Join("test", filename))

	if err != nil {
		panic(err)
	}

	return string(buf)
}
