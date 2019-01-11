package project

import "testing"

func TestGithubRequest(t *testing.T) {
	_, err := UpdateAccountInfo("berlingoqc")
	if err != nil {
		t.Fatal(err)
	}
	_, err = UpdateRepositoryInfo("berlingoqc")
	if err != nil {
		t.Fatal(err)
	}
}
