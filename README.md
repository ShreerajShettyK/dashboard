# GIT and AWS Dashboard

This project is a dashboard application that provides metrics and insights for AWS and Git repositories. It uses Gorilla Mux for routing and MongoDB for database storage. The frontend is built using HTML, Bootstrap, CSS, and JavaScript, with templates.

## Getting Started

1. Ensure you have Go installed on your system.
2. Clone the repository.
3. Install dependencies:

## Steps

```
go mod tidy
```

4. Set up your MongoDB connection string in the configuration.
5. Configure AWS credentials and region.
6. Run the program:

```
go run main.go

```

## Endpoints

### AWS Endpoints

- `/aws_metrics/home` (GET): Displays the AWS metrics home page.
- `/aws_metrics/home/resources` (GET): Retrieves AWS resource metrics.
- `/aws_billing/services` (GET): Lists AWS services for billing information.
- `/aws_billing/service/{service}/instances` (GET): Lists instances for a specific AWS service.

### Git Endpoints

- `/git_metrics/home` (GET): Displays the Git metrics home page.
- `/git_metrics/home/commits` (GET): Retrieves Git commit metrics.
- `/git_metrics/home/repos` (GET): Lists Git repositories.
- `/git_metrics/home/authors` (GET): Lists Git authors.
- `/git_metrics/repoAuthors` (GET): Retrieves authors by repository.

## Components

- **Router**: Gorilla Mux
- **Database**: MongoDB
- **Frontend**: HTML, Bootstrap, CSS, JavaScript (templates in the `templates` folder)

## AWS Services Used

The `helpers` folder contains utilities for interacting with various AWS services:

- **CloudWatch**: Monitoring and observability service (cloudWatch.go)
- **Cost Explorer**: Visualize and manage AWS costs and usage (costExplorerHelper.go)
- **EC2**: Elastic Compute Cloud for scalable computing capacity (ec2Helper.go)
- **Secrets Manager**: Securely store and manage secrets (used in config)

## Additional Information

- The project uses Go for backend development.
- AWS SDK for Go is used to interact with AWS services.
- Ensure proper AWS credentials and permissions are set up for accessing the required services.
- MongoDB is used as the database. Make sure to configure the connection string correctly.
- The frontend is built using HTML, Bootstrap for styling, and JavaScript for interactivity.