package project

import (
	"context"
	"errors"

	"github.com/google/go-github/v21/github"
)

func UpdateAccountInfo(user string) (*GitHubAccount, error) {

	client := github.NewClient(nil)

	u, _, err := client.Users.Get(context.Background(), user)
	if err != nil {
		return nil, err
	}
	ua := &GitHubAccount{
		Name:           user,
		URL:            *u.URL,
		ImgURL:         *u.AvatarURL,
		Location:       *u.Location,
		Email:          u.GetEmail(),
		Bio:            u.GetBio(),
		NbrPublicRepo:  u.GetPublicRepos(),
		NbrPublicGists: u.GetPublicGists(),
		NbrFollowing:   u.GetFollowing(),
		NbrFollorwers:  u.GetFollowers(),
	}
	return ua, nil
}

func UpdateRepositoryInfo(user string, name ...string) ([]*GitHubRepo, error) {
	client := github.NewClient(nil)
	opt := &github.RepositoryListOptions{
		Sort: "updated",
	}
	repos, _, err := client.Repositories.List(context.Background(), user, opt)
	if err != nil {
		return nil, err
	}
	if len(repos) == 0 {
		return nil, errors.New("No repository")
	}
	ll := make([]*GitHubRepo, 0)
	for _, repo := range repos {
		valid := false
		for _, n := range name {
			na := *repo.Name
			if na == n {
				valid = true
			}
		}
		if !valid {
			continue
		}
		l := &GitHubRepo{
			URl:          repo.GetURL(),
			Name:         repo.GetName(),
			Description:  repo.GetDescription(),
			StarCount:    repo.GetStargazersCount(),
			ForksCount:   repo.GetForksCount(),
			LastUpdateOn: repo.GetUpdatedAt().Time,
			CreatedOn:    repo.GetCreatedAt().Time,
		}
		ll = append(ll, l)

	}
	return ll, nil
}
