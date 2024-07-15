package models

import "time"

type GitMetric struct {
	ID                  string    `bson:"_id,omitempty"`
	CommitMessage       string    `bson:"commit_message"`
	FilesDeleted        int       `bson:"files_deleted"`
	FilesUpdated        int       `bson:"files_updated"`
	LinesDeleted        int       `bson:"lines_deleted"`
	CommitID            string    `bson:"commit_id"`
	CommittedBy         string    `bson:"commited_by"`
	FilesAdded          int       `bson:"files_added"`
	LinesAdded          int       `bson:"lines_added"`
	LinesUpdated        int       `bson:"lines_updated"`
	RepoName            string    `bson:"reponame"`
	CommitDate          time.Time `bson:"commit_date"`
	FormattedCommitDate string    // New field for formatted date
}

// GitMetricsViewData holds the data for rendering the Git metrics view
type GitMetricsViewData struct {
	Metrics           []GitMetric
	Repos             []string
	Authors           []string
	RepoNameParameter string
	UserNameParameter string
	CurrentPage       int64
	PreviousPage      int64
	NextPage          int64
}
