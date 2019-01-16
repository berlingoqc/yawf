package utility

import "io/ioutil"

func ReadFileString(filepath string) (string, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
