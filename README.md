# Backup Cobra

Tool to backup and analyze Bitbucket repositories. Features include cloning repositories, generating Excel reports, creating zip archives, and uploading to Amazon S3.

## Installation

### 1. Prerequisites
- Go 1.22 or higher
- A Bitbucket account with API access
- AWS credentials (if using S3 upload feature)

### 2. Get the code
```bash
git clone https://github.com/s3pweb/gitArchiveS3Report.git
cd gitArchiveS3Report
```

### 3. Install dependencies

#### Linux Installation
```bash
# Install ghorg (required for cloning repositories)
go install github.com/gabrie30/ghorg@v1.8.4

# Install project libraries
go mod tidy
```

#### macOS Installation
```bash
# Option 1: Install via Go (recommended)
go install github.com/gabrie30/ghorg@v1.8.4
# Add Go bin to your PATH if not already done
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc  # For zsh
# OR
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bash_profile  # For bash
source ~/.zshrc  # Or source ~/.bash_profile for bash

# Option 2: Install via Homebrew (alternative)
brew install ghorg

# Install project libraries
go mod tidy
```

Note: Why two different installation commands?
- `go install` is used for ghorg because it's an external executable tool that our application calls as a separate process.
  You can see this in `clone_repos.go` where we use `exec.Command("ghorg", ...)` to run it.
- `go mod tidy` installs library dependencies (like cobra, excelize) that are imported and used directly in our code,
  as defined in `go.mod`.

### 4. Configure the tool

Create configuration files:
```bash
# Create secrets file
cat > .secrets << EOF
# For backupcobra
BITBUCKET_TOKEN=your_token
BITBUCKET_USERNAME=your_username
BITBUCKET_WORKSPACE=your_workspace
AWS_ACCESS_KEY_ID=your_aws_key
AWS_SECRET_ACCESS_KEY=your_aws_secret
AWS_REGION=your_aws_region
AWS_BUCKET_NAME=your_bucket_name
UPLOAD_KEY=your_upload_prefix

# For ghorg directly (required if using brew install on macOS)
GHORG_BITBUCKET_USERNAME=your_username
GHORG_BITBUCKET_TOKEN=your_token
EOF

# Source the secrets file to set environment variables
source .secrets

# Create config file
cat > .config << EOF
CPU=1
DEFAULT_COLUMN=RepoName;BranchName;LastCommitDate;TimeSinceLastCommit;Commitnbr;HostLine;LastDeveloper;LastDeveloperPercentage
TERMS_TO_SEARCH=vaumt;swagger
FILES_TO_SEARCH=(?i)sonar-project.properties$;(?i)bitbucket-pipelines.yml$;(?i)Dockerfile$;(?i)docker-compose(-\\w+)?\\.yaml$
LOG_LEVEL=info
EOF
```

### 5. Build the tool
```bash
go build -o backupcobra
```

## Basic Usage

### Display available commands
```bash
./backupcobra
```

## Command Details

### Clone Options
```bash
./backupcobra clone [flags]
  -d, --dir-path string   Directory for cloned repositories (default: ./repositories)
  -m, --main-only        Clone only main/master branches
  -s, --shallow          Perform shallow clone (latest commit only)
```

### Report Options
```bash
./backupcobra report [flags]
  -p, --dir-path string   Path to repositories directory
```

### Zip Options
```bash
./backupcobra zip [flags]
  -p, --dir-path string   Path to repositories directory
```

### Upload Options
```bash
./backupcobra upload [flags]
  -p, --dir-path string   Path to directory to upload (required)
```