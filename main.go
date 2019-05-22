package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	authentication "k8s.io/api/authentication/v1beta1"
)

func checkGitUser(ctx context.Context, client *github.Client) string {
	// Check User
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return ""
	}
	return *user.Login
}

func checkGitGroups(ctx context.Context, client *github.Client, orgName string) []string {
	// Check User's Groups
	var groups []string
	opt := &github.ListOptions{PerPage: 10}
	for {
		teams, resp, err := client.Teams.ListUserTeams(ctx, opt)
		if err != nil {
			return groups
		}
		for _, team := range teams {
			if team.Organization.GetLogin() == orgName {
				groups = append(groups, team.GetName())
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return groups
}

func outHandler(w http.ResponseWriter, httpStatus int, trs authentication.TokenReviewStatus, msg string) {
	log.Println(msg)
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"apiVersion": "authentication.k8s.io/v1beta1",
		"kind":       "TokenReview",
		"status":     trs,
	})
	return
}

func main() {
	orgName := os.Getenv("GIT_ORG")

	http.HandleFunc("/authenticate", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var tr authentication.TokenReview
		err := decoder.Decode(&tr)
		if err != nil {
			trs := authentication.TokenReviewStatus{
				Authenticated: false,
			}
			outHandler(w, http.StatusBadRequest, trs, fmt.Sprintf("[Error] %s", err.Error()))
			return
		}

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: tr.Spec.Token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)

		user := checkGitUser(ctx, client)
		groups := checkGitGroups(ctx, client, orgName)

		if len(groups) == 0 || len(user) == 0 {
			trs := authentication.TokenReviewStatus{
				Authenticated: false,
			}
			// TODO: below function fail
			outHandler(w, http.StatusUnauthorized, trs, "[Error] Unauthorized user")
			return
		}

		trs := authentication.TokenReviewStatus{
			Authenticated: true,
			User: authentication.UserInfo{
				Username: user,
				UID:      user,
				Groups:   groups,
			},
		}
		outHandler(w, http.StatusOK, trs, fmt.Sprintf("[Success] login as %s in groups %s", user, groups))

	})
	log.Fatal(http.ListenAndServe(":3210", nil))
}
