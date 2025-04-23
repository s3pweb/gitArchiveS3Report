# Git Archive S3 Report

Tool to backup and analyze Bitbucket repositories. Features include cloning repositories, generating Excel reports, creating zip archives, and uploading to Amazon S3.

## Installation

### 1. Prerequisites
- Go 1.22 or higher (https://golang.org/dl/)
- A Bitbucket account with API access
- AWS credentials (if using S3 upload feature)

### 1.1 Generate Bitbucket Token
![alt text](resources/image.png)

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
https://github.com/gabrie30/ghorg?tab=readme-ov-file#installation

```bash
# Verify ghorg installation
which ghorg
```

### Configuration

#### Edit the `.env` file with your settings:
```bash
# Bitbucket Configuration
BITBUCKET_TOKEN=your_token
BITBUCKET_USERNAME=your_username
BITBUCKET_WORKSPACE=your_workspace

# AWS Configuration (Optional: For S3 upload feature)
AWS_ACCESS_KEY_ID=your_aws_key
AWS_SECRET_ACCESS_KEY=your_aws_secret
AWS_REGION=your_aws_region
AWS_BUCKET_NAME=your_bucket_name
AWS_PATH=repositories # this is the path created in the bucket where the zip file will be uploaded

# Logger Configuration
LOG_LEVEL=debug

# Application Configuration
# The number of CPU cores to use for cloning repositories
CPU=1

# Developer name mappings (optional)
# Format: alias1=Real Name 1;alias2=Real Name 2
DEVELOPERS_MAP=john=John Doe;jane=Jane Smith

# Default columns for the Excel report
DEFAULT_COLUMN=RepoName;BranchName;LastCommitDate;TimeSinceLastCommit;Commitnbr;HostLine;LastDeveloper;LastDeveloperPercentage;SelectiveCount;Count;ForbiddenCount

# Search terms and files for analysis
TERMS_TO_SEARCH=vault;swagger
FILES_TO_SEARCH=(?i)sonar-project.properties$;(?i)bitbucket-pipelines.yml$;(?i)Dockerfile$;(?i)docker-compose(-\w+)?\.yaml$

# Files that should NOT be present in the repository
FORBIDDEN_FILES_TO_SEARCH=(?i)\.npmrc$;(?i)\.env$;(?i)password\.txt$;(?i)credentials\.(json|txt|yaml)$;(?i)secret\.(key|txt|json)$;(?i)private\.key$;(?i)api_token\.txt$

# Terms and files to be counted separately (subset of the search terms and files)
TERMS_FILES_TO_COUNT=(?i)bitbucket-pipelines.yml$;(?i)sonar-project.properties$;vault

# Count thresholds (percentage values)
COUNT_THRESHOLD_LOW=30    # Below this percentage will be red
COUNT_THRESHOLD_MEDIUM=60 # Below this percentage will be orange, above will be green

# Default clone directory (where the repositories will be cloned)
DIR=../repositories
# Default zip directory (where the zip files will be stored)
DEST_DIR=../zipped

# JIRA Task Creation Configuration
JIRA_BASE_URL=https://example.atlassian.net
JIRA_TASK_ENABLED=true
JIRA_PARENT_TASK=S3DEVAGR-2180
JIRA_PROJECT_KEY=S3DEVAGR
JIRA_ISSUE_TYPE=Sous-tâche
JIRA_USERNAME=your_username
JIRA_API_TOKEN=your_token
JIRA_TITLE_TEMPLATE=Amélioration de la CI/CD pour le projet {{.RepoName}}
JIRA_DESC_TEMPLATE={{.MissingElements}}
# Documentation links format: Display Text|URL (separated by semicolons)
JIRA_DOC_LINKS=Title 1|https://example.atlassian.net/wiki/spaces/DEV/pages/123/CI+CD+Introduction;Title 2|https://example.atlassian.net/wiki/spaces/DEV/pages/456/External+Projects
```

The `JIRA_DOC_LINKS` field supports two formats:
- Simple URLs: `https://example.com/page1;https://example.com/page2`
- Custom display text: `Display Text|https://example.com;Another Text|https://example.com`

When using custom display text, the format is `Display Text|URL`, with multiple links separated by semicolons.

## Usage

### Display Available Commands
```bash
./git-archive-s3
```

### Clone Repositories
```bash
./git-archive-s3 clone [flags]
  -p, --dir-path string   Directory for cloned repositories (default: ./repositories) (optional)
  -m, --main-only         Clone only main/master/develop branches (optional)
  -s, --shallow           Perform shallow clone (latest commit only) (optional)
```

### Generate Report
```bash
./git-archive-s3 report [flags]
  -p, --dir-path string   Path to repositories directory (optional)
  -d, --dev-sheets        Generate developer-specific sheets (optional)
```

### Create ZIP Archive and Optionally Upload
```bash
./git-archive-s3 zip [flags]
  -p, --dir-path string       Path to source directory or file (required)
  -d, --dest-path string      Destination path to save the zip file (optional)
  -u, --upload                Upload the created zip file to S3 (optional)
  -r, --remove                Delete local zip file after successful upload (requires --upload)
```

#### Zip Examples:
```bash
# Create a zip archive from a directory
./git-archive-s3 zip -p /path/to/repositories

# Create a zip archive from a single file
./git-archive-s3 zip -p /path/to/specific/file.txt

# Specify custom destination for the zip file
./git-archive-s3 zip -p /path/to/source -d /custom/zip/location

# Create a zip archive and upload it to S3
./git-archive-s3 zip -p /path/to/repositories -u

# Create a zip archive, upload it to S3, and remove the local file
./git-archive-s3 zip -p /path/to/repositories -u -r
```

### Upload to S3
```bash
./git-archive-s3 upload [flags]
  -p, --dir-path string   Path to directory containing zip files or path to specific zip file (required)
  -a, --all               Upload all zip files in the specified directory (optional)
  -l, --last              Upload only the most recent zip file in the directory (optional)
```

#### Upload Examples:
```bash
# Upload a specific zip file
./git-archive-s3 upload -p /path/to/repository.zip

# Upload the most recent zip file from a directory 
./git-archive-s3 upload -p /path/to/zip/folder -l

# Upload all zip files from a directory
./git-archive-s3 upload -p /path/to/zip/folder -a

# When multiple zip files are found without specifying an option,
# the tool will list available files and options
```

### Forbidden Files Detection
The tool checks for files that should NOT be present in repositories (like .npmrc files containing tokens). Configure the list in the .env file:
```
FORBIDDEN_FILES_TO_SEARCH=(?i)\.npmrc$;(?i)\.env$;(?i)password.txt$
```

### JIRA Task Creation
The Excel report includes a "Create JIRA Task" button for each repository row that has missing elements. Clicking this button will create a JIRA task with customizable title and description.

The templates support the following variables:
- `{{.RepoName}}`: Repository name
- `{{.BranchName}}`: Branch name
- `{{.TopDeveloper}}`: Top contributor to the repository
- `{{.LastDeveloper}}`: Last developer who committed
- `{{.MissingElements}}`: List of missing elements (formatted as bullet points)
- `{{.ParentTask}}`: Parent JIRA task reference

#### Why a Local Server is Required
To use JIRA task creation from Excel files, you need to start a local server:
```bash
./git-archive-s3 serve
```

This local server is necessary because:
- Excel cannot make authenticated API calls directly to JIRA
- For security reasons, API credentials should not be stored in Excel files
- The local server acts as a secure bridge: it receives requests from Excel, retrieves credentials from your `.env` file, and makes authorized API calls to JIRA

The server runs locally on port 8081 by default (configurable with the `-p` flag).

## Notes
- Environment variables can be used to override any setting from the `.env` file
