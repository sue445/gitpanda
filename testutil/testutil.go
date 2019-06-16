package testutil

import (
	"io/ioutil"
)

func ReadTestData(filename string) string {
	buf, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	return string(buf)
}
