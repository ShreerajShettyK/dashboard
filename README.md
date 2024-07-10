endpoint to get a specific user and repo name:
http://localhost:8000/git_metrics?user_name=user_yjaCZ&repo_name=repo_iuyWI

Test entire projects code coverage

go test ./... -coverprofile cover.out
go tool cover -func cover.out