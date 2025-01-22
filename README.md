# Backup Cobra

Tool to backup and analyze Bitbucket repositories. Features include cloning repositories, generating Excel reports, creating zip archives, and uploading to Amazon S3.

## Installation

### 1. Prerequisites
- Go 1.22 or higher (https://golang.org/dl/)
- A Bitbucket account with API access
- AWS credentials (if using S3 upload feature)

### 2. Install dependencies

### Installation Steps

1. Clone and build the project:
```bash
# Get the code
git clone https://github.com/s3pweb/gitArchiveS3Report.git
cd gitArchiveS3Report

# Install dependencies and build
go mod tidy
go build -o git-archive-s3
```

2. Install ghorg:
```bash
go install github.com/gabrie30/ghorg@v1.8.4

# Verify ghorg installation
which ghorg
```

### Configuration

1. Create configuration files:
```bash
# Copy example files
cp example.secrets .secrets
cp example.config .config
```

2. Edit `.secrets` file:
```bash
export BITBUCKET_TOKEN=your_token
export BITBUCKET_USERNAME=your_username
export BITBUCKET_WORKSPACE=your_workspace


# Optional: For S3 upload feature
export AWS_ACCESS_KEY_ID=your_aws_key
export AWS_SECRET_ACCESS_KEY=your_aws_secret
export AWS_REGION=your_aws_region
export AWS_BUCKET_NAME=your_bucket_name
export UPLOAD_KEY=your_upload_prefix
```

3. Edit `.config` file:
```bash
CPU=1
DEFAULT_COLUMN=RepoName;BranchName;LastCommitDate;TimeSinceLastCommit;Commitnbr;HostLine;LastDeveloper;LastDeveloperPercentage
TERMS_TO_SEARCH=vaumt;swagger
FILES_TO_SEARCH=(?i)sonar-project.properties$;(?i)bitbucket-pipelines.yml$;(?i)Dockerfile$;(?i)docker-compose(-\\w+)?\\.yaml$
LOG_LEVEL=info
```

4. Load the secrets:
```bash
source .secrets
```

Modify these values according to your needs.

## Usage

### Display Available Commands
```bash
./git-archive-s3
```

### Clone Repositories
```bash
./git-archive-s3 clone [flags]
  -d, --dir-path string   Directory for cloned repositories (default: ./repositories)
  -m, --main-only        Clone only main/master branches
  -s, --shallow          Perform shallow clone (latest commit only)
```

### Generate Report
```bash
./git-archive-s3 report [flags]
  -p, --dir-path string   Path to repositories directory
```

### Create ZIP Archives
```bash
./git-archive-s3 zip [flags]
  -p, --dir-path string   Path to repositories directory
```

### Upload to S3
```bash
./git-archive-s3 upload [flags]
  -p, --dir-path string   Path to directory to upload
```