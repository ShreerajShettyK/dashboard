package models

import "time"

type GitMetric struct {
	ID           int
	RepoName     string
	UserName     string
	CommitID     string
	CommitDate   time.Time
	LinesAdded   int
	LinesDeleted int
	FilesAdded   int
	FilesDeleted int
}
