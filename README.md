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
DEFAULT_COLUMN=RepoName;BranchName;LastCommitDate;TimeSinceLastCommit;Commitnbr;HostLine;LastDeveloper;LastDeveloperPercentage;SelectiveCount;Count

# Search terms and files for analysis
TERMS_TO_SEARCH=vault;swagger
FILES_TO_SEARCH=(?i)sonar-project.properties$;(?i)bitbucket-pipelines.yml$;(?i)Dockerfile$;(?i)docker-compose(-\w+)?\.yaml$

# Terms and files to be counted separately (subset of the search terms and files)
TERMS_FILES_TO_COUNT=(?i)bitbucket-pipelines.yml$;(?i)sonar-project.properties$;vault

# Terms and files to be counted separately (subset of the search terms and files)
# These items will be used for the SelectiveCount calculation
TERMS_FILES_TO_COUNT=(?i)bitbucket-pipelines.yml$;(?i)sonar-project.properties$;vault

# Count thresholds (percentage values)
# These thresholds determine the color-coding in the report:
COUNT_THRESHOLD_LOW=30    # Below this percentage will be red
COUNT_THRESHOLD_MEDIUM=60 # Below this percentage will be orange, above will be green

# Count thresholds (percentage values)
COUNT_THRESHOLD_LOW=30    # Below this percentage will be red
COUNT_THRESHOLD_MEDIUM=60 # Below this percentage will be orange, above will be green

# Default clone directory (where the repositories will be cloned)
DIR=../repositories
# Default zip directory (where the zip files will be stored)
DEST_DIR=../zipped
```

## Explanation of FILES_TO_SEARCH regex patterns:
- `(?i)`: Case-insensitive matching
- `$`: End of the string

### Examples
- `(?i)sonar-project.properties$`
  - Matches: `sonar-project.properties`, `Sonar-Project.Properties`
  - Does not match: `sonar-project.properties.txt`, `my-sonar-project.properties`
- `(?i)docker-compose(-\w+)?\.yaml$`
  - Matches: `docker-compose.yaml`, `docker-compose-test.yaml`
  - Does not match: `docker-compose.yaml.backup`

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

## Notes
- Environment variables can be used to override any setting from the `.env` file