// package models

// import "time"

// type GitMetric struct {
// 	ID            int       `bson:"_id,omitempty"`
// 	RepoName      string    `bson:"repo_name"`
// 	UserName      string    `bson:"user_name"`
// 	CommitID      string    `bson:"commit_id"`
// 	CommitDate    time.Time `bson:"commit_date"`
// 	CommittedBy   string    `bson:"committed_by"`
// 	CommitMessage string    `bson:"commit_message"`
// 	FilesAdded    int       `bson:"files_added"`
// 	FilesDeleted  int       `bson:"files_deleted"`
// 	FilesUpdated  int       `bson:"files_updated"`
// 	LinesAdded    int       `bson:"lines_added"`
// 	LinesUpdated  int       `bson:"lines_updated"`
// 	LinesDeleted  int       `bson:"lines_deleted"`
// }

package models

import "time"

type GitMetric struct {
	ID            string    `bson:"_id,omitempty"`
	CommitMessage string    `bson:"commit_message"`
	FilesDeleted  int       `bson:"files_deleted"`
	FilesUpdated  int       `bson:"files_updated"`
	LinesDeleted  int       `bson:"lines_deleted"`
	CommitID      string    `bson:"commit_id"`
	CommittedBy   string    `bson:"commited_by"`
	FilesAdded    int       `bson:"files_added"`
	LinesAdded    int       `bson:"lines_added"`
	LinesUpdated  int       `bson:"lines_updated"`
	RepoName      string    `bson:"reponame"`
	CommitDate    time.Time `bson:"commit_date"`
}
