package main

// could possibly use concurrency for sending mnulitple requests at once

import (
	"encoding/json"
	"fmt"
	"github.com/SPatel29/github-user-activity/structs"
	"io"
	"net/http"
	"os"
)

var myMap = make(map[string]interface{})
const (
	PUSH_EVENT = "PushEvent"
)

func main() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/users/samtiz/events", nil)
	if err != nil {
		fmt.Println("error creating request:", err.Error())
		os.Exit(1)
	}
	req.Header.Add("accept", "application/vnd.github+json")
	resp, err := client.Do(req) // sends the request. Does this send the request async or synchronously?
	if err != nil {
		fmt.Println("error getting response", err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Cannot read contents of response body")
	}
	var events []interface{}
	if err := json.Unmarshal(body, &events); err != nil {
		fmt.Println("Error decoding JSON:", err)
		os.Exit(1)
	}
	var test []structs.PushEvent
	myMap[PUSH_EVENT] = test
	fmt.Printf("Number of events: %d\n", len(events))
	var pushEventCounter = 0
	var deleteEventCounter = 0
	var unKnownEventCounter = 0
	pushEventInterface := myMap[PUSH_EVENT] 
	pushEventSlice := pushEventInterface.([]structs.PushEvent) 
	for _, event := range events {
		eventMap := event.(map[string]interface{})
		eventType := eventMap["type"].(string)     
		repo := eventMap["repo"].(map[string]interface{})
		repoName := repo["name"].(string)
		repoUrl := repo["url"].(string)
		actor := eventMap["actor"].(map[string]interface{})
		displayLogin := actor["display_login"].(string)
		profileUrl := actor["url"].(string)
		payload := eventMap["payload"].(map[string]interface{})
		switch eventType {
		case "PushEvent":
			commits := payload["commits"].([]interface{})
			numberOfCommits := payload["size"].(float64)
			createdAt := eventMap["created_at"].(string)
			commitSlice := []structs.Commit{}
			for _, commit := range commits {
				commitMap := commit.(map[string]interface{})
				commitMessage := commitMap["message"].(string)
				commitUrl := commitMap["url"].(string)
				author := commitMap["author"].(map[string]interface{})
				authorName := author["name"].(string)
				authorEmail := author["email"].(string)
				commitSlice = append(commitSlice, structs.Commit{
					Message: commitMessage,
					Author: structs.Author{
						Email: authorEmail,
						Name: authorName,
					},
					Url: commitUrl,
				})
			}
			pushEvent := structs.PushEvent{
				CommonFields: structs.CommonFields{
					RepoName: repoName, RepoUrl: repoUrl, DisplayLogin: displayLogin, ProfileUrl: profileUrl},
				Commits: commitSlice,
				NumberOfCommits: numberOfCommits,
				CreatedAt: createdAt,
			}
			pushEventSlice = append(pushEventSlice, pushEvent)
			pushEventCounter += 1
		case "DeleteEvent":
			deleteEventCounter += 1
		default:
			unKnownEventCounter += 1
		}

	}
	fmt.Printf("PushEvent Counter was: %d, DeleteEvent Counter was: %d, unKnownEventCounter was: %d\n", pushEventCounter, deleteEventCounter, unKnownEventCounter)
}
