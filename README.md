# Backup Cobra

Backup Cobra is a powerful tool designed to clone repositories, generate detailed reports in Excel format, zip repositories and upload zip files to Amazon S3. It helps you keep track of your repositories, developers' activities, and various project metrics.

## Description

Backup Cobra automates the process of cloning repositories, generate detailed reports in Excel format, zip repositories and upload zip files to Amazon S3. It supports Bitbucket repositories and provides insights into developer activities, commit histories, and project configurations.

## Features

- **Clone Repositories**: Clone all repositories from a Bitbucket workspace using `ghorg`.
- **Generate Excel Reports**: Create Excel reports with detailed branch information, including last commit dates, top developers, and various project metrics.
- **Zip Repositories**: Generate zip files for each repositories of a workspace.
- **Upload to AmazonS3**: Upload a zip files, repositories of a project into amazon S3.

## Installation

1. Clone the repository :
   ```bash
   git clone https://github.com/yourusername/backup-cobra.git
   cd backup-cobra

2. Install Go : Make sure Go is installed. You can dowload it from the official Go website 
   https://go.dev/doc/install

3. Install dependencies :
- [ghorg](https://github.com/gabrie30/ghorg) (for cloning repositories)
      `go install github.com/gabrie30/ghorg@latest`
- [Excelize](https://github.com/xuri/excelize) (for generating Excel reports)
      `go get -u github.com/xuri/excelize/v2`
- [cobra](https://github.com/spf13/cobra)
      `go get -u github.com/spf13/cobra@latest`

4. Install project dependencies : From the project root, run : 
   `go mod tidy`

5. Configure AWS and Bitbucket credentials :
To store sensitive information like API keys or credentials, you can use a `.secrets` file. An example file named `example.secrets` is provided in the project to help you get started.
- Copy the `example.secrets` file and rename it to `.secrets`:
   ```bash
   cp example.secrets .secrets

- Fill it with the secrets required

- For help to get the keys for aws amazon S3 you can use this documentation : `https://docs.aws.amazon.com/fr_fr/IAM/latest/UserGuide/id_credentials_access-keys.html`

6. Configure your .config file :
The `.config` file allows you to customize certain behaviors of the project. An exemple file name `example.config` is provided in the project to help you get started.
- Copy the `example.config` file and rename it to `.config`:
   ```bash
   cp example.config .config

In this `.config` file you have some variables you can configurate :

- **CPU** : you can choose the number of processors you want to use (by default 1).

- **DEVELOPERS_MAP** : For developpers it will group some username, exemple : "DEVELOPERS_MAP=Jean D=Jean Dupont;"

- **DEFAULT_COLUMN** : The list of the default column that the program do

- **TERMS_TO_SEARCH** : List of terms you want to search for in the files of the repository (handle regex)

- **FILES_TO_SEARCH** : List of file you want to search for in a repository (handle regex)

## Usage 

### Clone Repositories

`./backupcobra clone --dir-path <path where to clone the repositories>`

### Generate an Excel Report

`./backupcobra report --dir-path <path where the repositories were cloned>`

### Compress Directories

`./backupcobra zip --dir-path <path where the repositories were cloned>`

### Upload to Amazon S3

`./backupcobra upload --dir-path <path where the repositories were cloned>`

## Features

[![](https://mermaid.ink/img/pako:eNqNk19vgjAUxb9Kc58ZQS0b8laxmWRqDeBmDC-NdBuJUIOwzKHffeXPnJu6rG_0nt8955a2hJWMBNggsmHMXzKehClSyxmzKUWHva7vS3Q_Yt49slGxFafV_f7mRpYohIEbDObOAw3QE_Me_BlxaAhKHycbmeUNck1VNVEWHp0x3w2Y51JfkSuZ5jxOtxftxswh4x9E7bblb228quQFX-mvAsdxWv3RgAzYPDjX85Svdx9XmMaDLhw6bmu_Ql2RtPOHMCQBQd_OiEyHaOCRqTOqG_08kaU7-_94lfiY02GTmUd9_1wfyQty4jkj95HWm2fEKhM8b03o4vTM_yRPorVUa1dRE7JkU-T3qp8-HzZnKNIINEhElvA4Une1rOAQ8leRiKpdCJF45sU6DyFMD0rKi1z6u3QFdp4VQoNMFi-vYD_z9VZ9FZtIxW6v-3F3w1OwS3gHG-u3tybuWriHsWl2MNZgB3bHNHQLd7vYwIbV7_Xv8EGDDylVB0PvN8vEd4bZtTqWBiKKc5lNmsdVv7HaYlkDteXhE3X_Aig?type=png)](https://mermaid.live/edit#pako:eNqNk19vgjAUxb9Kc58ZQS0b8laxmWRqDeBmDC-NdBuJUIOwzKHffeXPnJu6rG_0nt8955a2hJWMBNggsmHMXzKehClSyxmzKUWHva7vS3Q_Yt49slGxFafV_f7mRpYohIEbDObOAw3QE_Me_BlxaAhKHycbmeUNck1VNVEWHp0x3w2Y51JfkSuZ5jxOtxftxswh4x9E7bblb228quQFX-mvAsdxWv3RgAzYPDjX85Svdx9XmMaDLhw6bmu_Ql2RtPOHMCQBQd_OiEyHaOCRqTOqG_08kaU7-_94lfiY02GTmUd9_1wfyQty4jkj95HWm2fEKhM8b03o4vTM_yRPorVUa1dRE7JkU-T3qp8-HzZnKNIINEhElvA4Une1rOAQ8leRiKpdCJF45sU6DyFMD0rKi1z6u3QFdp4VQoNMFi-vYD_z9VZ9FZtIxW6v-3F3w1OwS3gHG-u3tybuWriHsWl2MNZgB3bHNHQLd7vYwIbV7_Xv8EGDDylVB0PvN8vEd4bZtTqWBiKKc5lNmsdVv7HaYlkDteXhE3X_Aig)

1. **Repositories Cloning : Clones repositories form a specified workspace.**

[![](https://mermaid.ink/img/pako:eNp1kVFvgjAQx79Kc89IAItC35QRR3TDAMuSpS8NVCURakpZ5oDvvgpm2ZLtnq69__9-vV4HuSg4EKA1lw8lO0pW0RrpCHbxc4iG3jT7Dm0e42SDCGob_rPa97OZ6BCFdZStX4JtmKHXONmm-1UQUtD6sroIqSbLf6pbE41Iwn2cRlmcRGGqnbmoFSvr5k_cLg5Wu1-Okdawd_08MKDismJlocfqbn4K6sQrToHotOAH1p4V1RMPWspaJdJrnQNRsuUGSNEeT0AO7NzoU3spmOL3f_m-vbAaSAcfQLC5WLjY8fAcY9e1MTbgCsR2LdPDjoMtbHn-3F_iwYBPIXQHy_SncPHSch3P9gzgRamEfJr2MK5jRLyNhhE5fAGzeH3X?type=png)](https://mermaid.live/edit#pako:eNp1kVFvgjAQx79Kc89IAItC35QRR3TDAMuSpS8NVCURakpZ5oDvvgpm2ZLtnq69__9-vV4HuSg4EKA1lw8lO0pW0RrpCHbxc4iG3jT7Dm0e42SDCGob_rPa97OZ6BCFdZStX4JtmKHXONmm-1UQUtD6sroIqSbLf6pbE41Iwn2cRlmcRGGqnbmoFSvr5k_cLg5Wu1-Okdawd_08MKDismJlocfqbn4K6sQrToHotOAH1p4V1RMPWspaJdJrnQNRsuUGSNEeT0AO7NzoU3spmOL3f_m-vbAaSAcfQLC5WLjY8fAcY9e1MTbgCsR2LdPDjoMtbHn-3F_iwYBPIXQHy_SncPHSch3P9gzgRamEfJr2MK5jRLyNhhE5fAGzeH3X)

2. **Excel Report Generation : Creates an Excel file with detailed information about repositories, branches, commit, etc.**

[![](https://mermaid.ink/img/pako:eNpVkF9rgzAUxb9KuM9W_JO0mrfSChPWKdqNIXkJTdoK1ZSYwDr1u8_qGOy83cv9nXM5PZyUkEBB6n3NL5o3rEWTks88K45oHFx36BFisC12L-lHgqo0R0WSZ2V6zIo0KRlQZDv5jxqG1Uot1GFbZW-oDNHuNXvfM0AUdbIV4EAjdcNrMUX3T5iBucpGPu0YCHnm9mYYsHacTrk1qny0J6BGW-mAVvZyBXrmt26a7F1wI3-__9veeQu0hy-g2F2vCQ4iHGJMiI-xAw-gPvHcCAcB9rAXxWG8waMD30pNDp4bLyJ445Eg8iMHpKiN0oelq7myOaKagTly_AEGf2IC?type=png)](https://mermaid.live/edit#pako:eNpVkF9rgzAUxb9KuM9W_JO0mrfSChPWKdqNIXkJTdoK1ZSYwDr1u8_qGOy83cv9nXM5PZyUkEBB6n3NL5o3rEWTks88K45oHFx36BFisC12L-lHgqo0R0WSZ2V6zIo0KRlQZDv5jxqG1Uot1GFbZW-oDNHuNXvfM0AUdbIV4EAjdcNrMUX3T5iBucpGPu0YCHnm9mYYsHacTrk1qny0J6BGW-mAVvZyBXrmt26a7F1wI3-__9veeQu0hy-g2F2vCQ4iHGJMiI-xAw-gPvHcCAcB9rAXxWG8waMD30pNDp4bLyJ445Eg8iMHpKiN0oelq7myOaKagTly_AEGf2IC)

3. **Directory Compression : Creates zip files for specified directories.**

[![](https://mermaid.ink/img/pako:eNptkMFugzAMhl8l8pkioEkLuVUMaUitQDDtUOUSkbRFKqRKE2kd8O5L6bQdWp9s6_P_2x6gUUICBanfWn7UvGM9crHPSzSNvj8OiMG2SDdbVGVlUecfRZVnNQNEkb3Kf3gcFwt1h9NiV1ZZXT_zQr3AN1X6nn9mc_NpotGSGwkedFJ3vBVuz-GuwcCcZCcZUJcKeeD2bBiwfnIot0bVt74BarSVHmhljyegB36-uspehFP8PfWve-E90AG-gGJ_tSI4ivESY0JCjD24AQ1J4Mc4inCAgzhZJms8efCtlFMI_OQRBK8DEsVh7IEUrVF693js_N_ZYj8PzJbTDzY0bzQ?type=png)](https://mermaid.live/edit#pako:eNptkMFugzAMhl8l8pkioEkLuVUMaUitQDDtUOUSkbRFKqRKE2kd8O5L6bQdWp9s6_P_2x6gUUICBanfWn7UvGM9crHPSzSNvj8OiMG2SDdbVGVlUecfRZVnNQNEkb3Kf3gcFwt1h9NiV1ZZXT_zQr3AN1X6nn9mc_NpotGSGwkedFJ3vBVuz-GuwcCcZCcZUJcKeeD2bBiwfnIot0bVt74BarSVHmhljyegB36-uspehFP8PfWve-E90AG-gGJ_tSI4ivESY0JCjD24AQ1J4Mc4inCAgzhZJms8efCtlFMI_OQRBK8DEsVh7IEUrVF693js_N_ZYj8PzJbTDzY0bzQ)

4. **Amazon S3 upload : Uploads directories or zip files to a secure S3 bucket.**

[![](https://mermaid.ink/img/pako:eNp1kV1rgzAUhv9KONdW1MZWc9c66aTdLOoYjNwETVuhmhLjWKf-96VaxgbbucrH-5yH5HSQi4IDAS4fSnaUrKI10hXs4ucQDb1p9h3aPMbJBhHUNvznbd_PZqJDFNZRtn4JtmGGXuNkm-5XQUhB58vqIqSakP9StyZakYT7OI2yOInCVJO5qBUr6-ZP3S4OVrtfxGhr2DsHAyouK1YW-kndjaagTrziFIheFvzA2rOiQOtBR1mrRHqtcyBKttwAKdrjCciBnRu9ay8FU_z-K9-nF1YD6eADCDYXCxc7Hp5j7Lo2xgZcgdiuZXrYcbCFLc-f-0s8GPAphO5gmf5ULl5aruPZngG8KJWQT9MMxlGMircRGJXDF0WyfEM?type=png)](https://mermaid.live/edit#pako:eNp1kV1rgzAUhv9KONdW1MZWc9c66aTdLOoYjNwETVuhmhLjWKf-96VaxgbbucrH-5yH5HSQi4IDAS4fSnaUrKI10hXs4ucQDb1p9h3aPMbJBhHUNvznbd_PZqJDFNZRtn4JtmGGXuNkm-5XQUhB58vqIqSakP9StyZakYT7OI2yOInCVJO5qBUr6-ZP3S4OVrtfxGhr2DsHAyouK1YW-kndjaagTrziFIheFvzA2rOiQOtBR1mrRHqtcyBKttwAKdrjCciBnRu9ay8FU_z-K9-nF1YD6eADCDYXCxc7Hp5j7Lo2xgZcgdiuZXrYcbCFLc-f-0s8GPAphO5gmf5ULl5aruPZngG8KJWQT9MMxlGMircRGJXDF0WyfEM)

## Project Structure

- cmd/ : Contains the project's CLI commands.

- processrepos/ : Main application logic.

- utils/ : some useful functions

## Authors

Main author : Louise Calvez
