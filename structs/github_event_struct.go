package structs

type WatchEvent struct {
	CommonFields
	Action string
}

type PushEvent struct {
	CommonFields
	Commits []Commit
	NumberOfCommits float64 
	CreatedAt string
}

type Commit struct {
	Message string
	Author Author
	Url string
}

type Author struct {
	Email string
	Name string
}

type CommonFields struct {
	RepoName string
	RepoUrl string
	DisplayLogin string
	ProfileUrl string
}

