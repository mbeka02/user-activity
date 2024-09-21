package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Repo struct {
	Name string `json:"name"`
}

type GithubEvent struct {
	Type      string `json:"type"`
	Repo      Repo   `json:"repo"`
	CreatedAt string `json:"created_at"`
	Payload   struct {
		Action  string `json:"action"`
		Ref     string `json:"ref"`
		RefType string `json:"ref_type"`
		Commits []struct {
			Message string `json:"message"`
		} `json:"commits"`
	}
}

func main() {

	argsWithoutProgram := os.Args[1:]

	//ensure user name and  activity have been passed
	if len(argsWithoutProgram) != 1 {
		log.Fatal("insufficient arguments , usage :app <username>")
	}
	username := argsWithoutProgram[0]
	events, err := fetchActivity(username)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("response=>%v", events)
	printActivity(username, events)
}

func fetchActivity(username string) ([]GithubEvent, error) {
	response, err := http.Get(fmt.Sprintf("https://api.github.com/users/%s/events/public", username))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexcpected status code:%v", response.StatusCode)
	}
	events := []GithubEvent{}
	//decode body and store it in events
	json.NewDecoder(response.Body).Decode(&events)
	return events, nil
}

func printActivity(username string, events []GithubEvent) {
	for _, event := range events {
		switch event.Type {
		case "CreateEvent":
			fmt.Printf("%s created a repository called : %s at : %s", username, event.Repo.Name, event.CreatedAt)
		case "DeleteEvent":
			fmt.Printf("%s was deleted by %s\n", event.Repo.Name, username)
		case "ForkEvent":
			fmt.Printf("%s forked:%s\n", event.Repo.Name, username)
		case "PushEvent":
			fmt.Printf(" %s pushed to this repository: %s\n", username, event.Repo.Name)
		case "IssuesEvent":
			fmt.Printf(" %s opened an issue in %s\n", username, event.Repo.Name)
		case "WatchEvent":
			fmt.Printf(" %s starred %s\n", username, event.Repo.Name)
		case "PullRequestEvent":
			fmt.Printf(" %s opened PR in %s\n", username, event.Repo.Name)
		case "IssueCommentEvent":
			fmt.Printf(" %s commented in %s\n", username, event.Repo.Name)
		case "PullRequestReviewCommentEvent":
			fmt.Printf(" %s commented on a PR in %s\n", username, event.Repo.Name)
		case "PullRequestReviewEvent":
			fmt.Printf("  %s reviewed a PR in %s\n", username, event.Repo.Name)
		default:
			fmt.Printf("Unknown event by %s: %s\n", event.Type, username)
		}
	}

}
