# Backup Cobra

Backup Cobra is a powerful tool designed to clone repositories, generate detailed reports in Excel format, zip repositories and upload zip files to Amazon S3. It helps you keep track of your repositories, developers' activities, and various project metrics.

## Description

Backup Cobra automates the process of cloning repositories, generating detailed Excel reports, creating zip archives, and uploading files to Amazon S3. It supports Bitbucket repositories and provides insights into developer activities, commit histories, and project configurations.

## Features

- **Clone Repositories**: Clone all repositories from a Bitbucket workspace using `ghorg`.
- **Generate Excel Reports**: Create detailed Excel reports with branch information, including last commit dates, developer statistics, and project metrics.
- **Zip Repositories**: Generate zip files for each repository in a workspace.
- **Upload to AmazonS3**: Upload zip files and repositories to Amazon S3.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/backup-cobra.git
   cd backup-cobra
   ```

2. Install Go:
   Make sure Go is installed. You can download it from the [official Go website](https://go.dev/doc/install)

3. Install required dependencies:
   - [ghorg](https://github.com/gabrie30/ghorg) (for cloning repositories)
     ```bash
     go install github.com/gabrie30/ghorg@latest
     ```
   - [Excelize](https://github.com/xuri/excelize) (for generating Excel reports)
     ```bash
     go get -u github.com/xuri/excelize/v2
     ```
   - [cobra](https://github.com/spf13/cobra)
     ```bash
     go get -u github.com/spf13/cobra@latest
     ```

4. Install project dependencies:
   ```bash
   go mod tidy
   ```

## Configuration

### AWS and Bitbucket Credentials
Create a `.secrets` file to store your credentials:

1. Copy the example file:
   ```bash
   cp example.secrets .secrets
   ```

2. Fill in the following required credentials in `.secrets`:
   ```
   BITBUCKET_TOKEN=your_token
   BITBUCKET_USERNAME=your_username
   BITBUCKET_WORKSPACE=your_workspace
   AWS_ACCESS_KEY_ID=your_aws_key
   AWS_SECRET_ACCESS_KEY=your_aws_secret
   AWS_REGION=your_aws_region
   AWS_BUCKET_NAME=your_bucket_name
   UPLOAD_KEY=your_upload_prefix
   ```

For help getting AWS credentials, refer to the [AWS IAM documentation](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html).

### Project Configuration
Create a `.config` file to customize the tool's behavior:

1. Copy the example file:
   ```bash
   cp example.config .config
   ```

2. Configure the following settings:
   - `CPU`: Number of processors to use (default: 1)
   - `DEVELOPERS_MAP`: Allows standardizing developer names in reports. Format: "StandardName=NameVariant;". For example:
   - "DEVELOPERS_MAP=John Doe=John D;John Doe=J.Doe": all commits from "John D" and "J.Doe" will be attributed to "John Doe"
   - `DEFAULT_COLUMN`: Default columns for Excel reports
   - `TERMS_TO_SEARCH`: Terms to search in repository files (supports regex)
   - `FILES_TO_SEARCH`: Files to search for in repositories (supports regex)

## Quick Start Guide

Get started with Backup Cobra in minutes! Here are the basic commands you'll need:

```bash
# List all available commands
./backupcobra

# Clone all repositories from your Bitbucket workspace
./backupcobra clone

# Generate an Excel report for your repositories
./backupcobra report

# Create zip archives of your repositories
./backupcobra zip

# Upload files to S3
./backupcobra upload
```

## Command Reference and Examples

### Clone Command
Clone repositories from your Bitbucket workspace:

```bash
./backupcobra clone [flags]
```

Flags:
- `-d, --dir-path`: Directory path where repositories will be cloned (optional, defaults to ./repositories)
- `-m, --main-only`: Clone only the main/master branch
- `-s, --shallow`: Perform a shallow clone with only the latest commit

### Report Command
Generate an Excel report with repository analytics:

```bash
./backupcobra report [flags]
```

Flags:
- `-p, --dir-path`: Path to the cloned repositories directory

The report includes:
- Branch information
- Last commit dates
- Time since last commit
- Developer statistics
- Configuration file presence
- Search term matches
- Custom metrics defined in .config

### Zip Command
Create zip archives of repositories:

```bash
./backupcobra zip [flags]
```

Flags:
- `-p, --dir-path`: Path to the directory containing repositories to zip

Archives will be created in an `./archive` directory.

### Upload Command
Upload files to Amazon S3:

```bash
./backupcobra upload [flags]
```

Flags:
- `-p, --dir-path`: Path to the directory to upload (required)

Files will be uploaded to the configured S3 bucket using the prefix specified in `UPLOAD_KEY`.

## Examples of Common Operations

### Basic Workflow Example
```bash
# 1. Clone all repositories
./backupcobra clone -d ./my-repos

# 2. Generate a report
./backupcobra report -p ./my-repos

# 3. Create zip archives
./backupcobra zip -p ./my-repos

# 4. Upload the archives to S3
./backupcobra upload -p ./archive
```

### Advanced Clone Examples
```bash
# Clone only the main branches with shallow history (faster)
./backupcobra clone -d ./repos -m -s

# Clone to a specific directory with full history
./backupcobra clone -d /path/to/backup/repos

# Clone with default settings (uses ./repositories directory)
./backupcobra clone
```

### Report Generation Examples
```bash
# Generate report with default settings
./backupcobra report

# Generate report for a specific directory
./backupcobra report -p /path/to/repos

# The report will be saved as workspace_report.xlsx in the specified directory
```

### Zip Operation Examples
```bash
# Zip repositories with default settings
./backupcobra zip -p ./repositories

# Zip repositories from a specific directory
./backupcobra zip -p /path/to/repos

# The zip files will be created in the ./archive directory
```

### Upload Examples
```bash
# Upload a specific directory to S3
./backupcobra upload -p ./archive

# Upload from a custom path
./backupcobra upload -p /path/to/files
```

### Common Patterns
1. **Daily Backup Pattern**:
   ```bash
   ./backupcobra clone -d ./daily-backup
   ./backupcobra report -p ./daily-backup
   ./backupcobra zip -p ./daily-backup
   ./backupcobra upload -p ./archive
   ```

2. **Quick Status Check Pattern**:
   ```bash
   # Use shallow clone for faster operation
   ./backupcobra clone -d ./quick-check -s
   ./backupcobra report -p ./quick-check
   ```

3. **Main Branch Analysis Pattern**:
   ```bash
   # Clone only main branches
   ./backupcobra clone -d ./main-analysis -m
   ./backupcobra report -p ./main-analysis
   ```

## Project Structure

```
.
├── cmd/                    # Command implementations
├── config/                 # Configuration handling
├── processrepos/          # Core processing logic
│   ├── excel/            # Excel report generation
│   └── ...
├── utils/                 # Utility functions
│   ├── git/              # Git operations
│   ├── logger/           # Logging functionality
│   ├── structs/          # Data structures
│   └── styles/           # Excel styling
└── main.go               # Application entry point
```

## Logging

The application uses a custom logging system with the following levels:
- ERROR: Critical errors (red)
- WARN: Warnings (yellow)
- INFO: General information (blue)
- DEBUG: Debugging information (cyan)
- TRACE: Detailed debugging (black)
- SUCCESS: Successful operations (green)

Configure the log level in your `.config` file using the `LOG_LEVEL` setting.

1. **Repositories Cloning : Clones repositories form a specified workspace.**

[![](https://mermaid.ink/img/pako:eNqNk19vgjAUxb9Kc58ZQS0b8laxmWRqDeBmDC-NdBuJUIOwzKHffeXPnJu6rG_0nt8955a2hJWMBNggsmHMXzKehClSyxmzKUWHva7vS3Q_Yt49slGxFafV_f7mRpYohIEbDObOAw3QE_Me_BlxaAhKHycbmeUNck1VNVEWHp0x3w2Y51JfkSuZ5jxOtxftxswh4x9E7bblb228quQFX-mvAsdxWv3RgAzYPDjX85Svdx9XmMaDLhw6bmu_Ql2RtPOHMCQBQd_OiEyHaOCRqTOqG_08kaU7-_94lfiY02GTmUd9_1wfyQty4jkj95HWm2fEKhM8b03o4vTM_yRPorVUa1dRE7JkU-T3qp8-HzZnKNIINEhElvA4Une1rOAQ8leRiKpdCJF45sU6DyFMD0rKi1z6u3QFdp4VQoNMFi-vYD_z9VZ9FZtIxW6v-3F3w1OwS3gHG-u3tybuWriHsWl2MNZgB3bHNHQLd7vYwIbV7_Xv8EGDDylVB0PvN8vEd4bZtTqWBiKKc5lNmsdVv7HaYlkDteXhE3X_Aig?type=png)](https://mermaid.live/edit#pako:eNqNk19vgjAUxb9Kc58ZQS0b8laxmWRqDeBmDC-NdBuJUIOwzKHffeXPnJu6rG_0nt8955a2hJWMBNggsmHMXzKehClSyxmzKUWHva7vS3Q_Yt49slGxFafV_f7mRpYohIEbDObOAw3QE_Me_BlxaAhKHycbmeUNck1VNVEWHp0x3w2Y51JfkSuZ5jxOtxftxswh4x9E7bblb228quQFX-mvAsdxWv3RgAzYPDjX85Svdx9XmMaDLhw6bmu_Ql2RtPOHMCQBQd_OiEyHaOCRqTOqG_08kaU7-_94lfiY02GTmUd9_1wfyQty4jkj95HWm2fEKhM8b03o4vTM_yRPorVUa1dRE7JkU-T3qp8-HzZnKNIINEhElvA4Une1rOAQ8leRiKpdCJF45sU6DyFMD0rKi1z6u3QFdp4VQoNMFi-vYD_z9VZ9FZtIxW6v-3F3w1OwS3gHG-u3tybuWriHsWl2MNZgB3bHNHQLd7vYwIbV7_Xv8EGDDylVB0PvN8vEd4bZtTqWBiKKc5lNmsdVv7HaYlkDteXhE3X_Aig)

2. **Excel Report Generation : Creates an Excel file with detailed information about repositories, branches, commit, etc.**

[![](https://mermaid.ink/img/pako:eNp1kVFvgjAQx79Kc89IAItC35QRR3TDAMuSpS8NVCURakpZ5oDvvgpm2ZLtnq69__9-vV4HuSg4EKA1lw8lO0pW0RrpCHbxc4iG3jT7Dm0e42SDCGob_rPa97OZ6BCFdZStX4JtmKHXONmm-1UQUtD6sroIqSbLf6pbE41Iwn2cRlmcRGGqnbmoFSvr5k_cLg5Wu1-Okdawd_08MKDismJlocfqbn4K6sQrToHotOAH1p4V1RMPWspaJdJrnQNRsuUGSNEeT0AO7NzoU3spmOL3f_m-vbAaSAcfQLC5WLjY8fAcY9e1MTbgCsR2LdPDjoMtbHn-3F_iwYBPIXQHy_SncPHSch3P9gzgRamEfJr2MK5jRLyNhhE5fAGzeH3X?type=png)](https://mermaid.live/edit#pako:eNp1kVFvgjAQx79Kc89IAItC35QRR3TDAMuSpS8NVCURakpZ5oDvvgpm2ZLtnq69__9-vV4HuSg4EKA1lw8lO0pW0RrpCHbxc4iG3jT7Dm0e42SDCGob_rPa97OZ6BCFdZStX4JtmKHXONmm-1UQUtD6sroIqSbLf6pbE41Iwn2cRlmcRGGqnbmoFSvr5k_cLg5Wu1-Okdawd_08MKDismJlocfqbn4K6sQrToHotOAH1p4V1RMPWspaJdJrnQNRsuUGSNEeT0AO7NzoU3spmOL3f_m-vbAaSAcfQLC5WLjY8fAcY9e1MTbgCsR2LdPDjoMtbHn-3F_iwYBPIXQHy_SncPHSch3P9gzgRamEfJr2MK5jRLyNhhE5fAGzeH3X)

3. **Directory Compression : Creates zip files for specified directories.**

[![](https://mermaid.ink/img/pako:eNpVkF9rgzAUxb9KuM9W_JO0mrfSChPWKdqNIXkJTdoK1ZSYwDr1u8_qGOy83cv9nXM5PZyUkEBB6n3NL5o3rEWTks88K45oHFx36BFisC12L-lHgqo0R0WSZ2V6zIo0KRlQZDv5jxqG1Uot1GFbZW-oDNHuNXvfM0AUdbIV4EAjdcNrMUX3T5iBucpGPu0YCHnm9mYYsHacTrk1qny0J6BGW-mAVvZyBXrmt26a7F1wI3-__9veeQu0hy-g2F2vCQ4iHGJMiI-xAw-gPvHcCAcB9rAXxWG8waMD30pNDp4bLyJ445Eg8iMHpKiN0oelq7myOaKagTly_AEGf2IC?type=png)](https://mermaid.live/edit#pako:eNpVkF9rgzAUxb9KuM9W_JO0mrfSChPWKdqNIXkJTdoK1ZSYwDr1u8_qGOy83cv9nXM5PZyUkEBB6n3NL5o3rEWTks88K45oHFx36BFisC12L-lHgqo0R0WSZ2V6zIo0KRlQZDv5jxqG1Uot1GFbZW-oDNHuNXvfM0AUdbIV4EAjdcNrMUX3T5iBucpGPu0YCHnm9mYYsHacTrk1qny0J6BGW-mAVvZyBXrmt26a7F1wI3-__9veeQu0hy-g2F2vCQ4iHGJMiI-xAw-gPvHcCAcB9rAXxWG8waMD30pNDp4bLyJ445Eg8iMHpKiN0oelq7myOaKagTly_AEGf2IC)

4. **Amazon S3 upload : Uploads directories or zip files to a secure S3 bucket.**

[![](https://mermaid.ink/img/pako:eNptkMFugzAMhl8l8pkioEkLuVUMaUitQDDtUOUSkbRFKqRKE2kd8O5L6bQdWp9s6_P_2x6gUUICBanfWn7UvGM9crHPSzSNvj8OiMG2SDdbVGVlUecfRZVnNQNEkb3Kf3gcFwt1h9NiV1ZZXT_zQr3AN1X6nn9mc_NpotGSGwkedFJ3vBVuz-GuwcCcZCcZUJcKeeD2bBiwfnIot0bVt74BarSVHmhljyegB36-uspehFP8PfWve-E90AG-gGJ_tSI4ivESY0JCjD24AQ1J4Mc4inCAgzhZJms8efCtlFMI_OQRBK8DEsVh7IEUrVF693js_N_ZYj8PzJbTDzY0bzQ?type=png)](https://mermaid.live/edit#pako:eNptkMFugzAMhl8l8pkioEkLuVUMaUitQDDtUOUSkbRFKqRKE2kd8O5L6bQdWp9s6_P_2x6gUUICBanfWn7UvGM9crHPSzSNvj8OiMG2SDdbVGVlUecfRZVnNQNEkb3Kf3gcFwt1h9NiV1ZZXT_zQr3AN1X6nn9mc_NpotGSGwkedFJ3vBVuz-GuwcCcZCcZUJcKeeD2bBiwfnIot0bVt74BarSVHmhljyegB36-uspehFP8PfWve-E90AG-gGJ_tSI4ivESY0JCjD24AQ1J4Mc4inCAgzhZJms8efCtlFMI_OQRBK8DEsVh7IEUrVF693js_N_ZYj8PzJbTDzY0bzQ)

## Author

Main author: Louise Calvez

## License

This project is licensed under the MIT License - see the LICENSE file for details.