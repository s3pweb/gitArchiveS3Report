package config

import (
	"os"
	"strings"

	"github.com/s3pweb/gitArchiveS3Report/utils"
	"github.com/s3pweb/gitArchiveS3Report/utils/logger"
	"github.com/spf13/viper"
)

var (
	// Global instance of the configuration
	cfg *Config
)

type Config struct {
	Bitbucket BitbucketConfig
	AWS       AWSConfig
	Logger    LoggerConfig
	App       AppConfig
}

type BitbucketConfig struct {
	Token     string
	Username  string
	Workspace string
}

type AWSConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	BucketName      string
	AWSUploadPath   string
}

type LoggerConfig struct {
	Level string
}

type AppConfig struct {
	CPU                  int
	DevelopersMap        string
	DefaultColumns       []string
	TermsToSearch        []string
	FilesToSearch        []string
	ForbiddenFiles       []string
	TermsFilesToCount    []string
	DefaultCloneDir      string
	DestDir              string
	MainBranchOnly       bool
	ShallowClone         bool
	DevSheets            bool
	CountThresholdLow    int
	CountThresholdMedium int
	JiraBaseURL          string
	JiraTaskEnabled      bool
	JiraParentTask       string
	JiraTitleTemplate    string
	JiraDescTemplate     string
	JiraDocLinks         []string
	JiraProjectKey       string
	JiraIssueType        string
	JiraUsername         string
	JiraAPIToken         string
}

// Init initializes the configuration
func Init() {
	log, _ := logger.NewLogger("Config", "trace")

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// Read .env file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Info("No .env file found, copying from example.env...")
			if err := copyFile("example.env", ".env"); err != nil {
				log.Error("Error copying example.env: %v", err)
				os.Exit(1)
			}
			if err := viper.ReadInConfig(); err != nil {
				log.Error("Error reading config file: %v", err)
				os.Exit(1)
			}
		} else {
			log.Error("Error reading config file: %v", err)
			os.Exit(1)
		}
	}

	// Configure viper mappings
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set default values
	viper.SetDefault("app.cpu", 1)
	viper.SetDefault("app.cloneDir", "./repositories")
	viper.SetDefault("app.countThresholdLow", 30)
	viper.SetDefault("app.countThresholdMedium", 60)

	// Load config into struct
	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		log.Error("Unable to decode config: %v", err)
		os.Exit(1)
	}

	// Application Configuration
	cfg.App.CPU = viper.GetInt("CPU")
	cfg.App.DefaultColumns = strings.Split(viper.GetString("DEFAULT_COLUMN"), ";")
	cfg.App.TermsToSearch = strings.Split(viper.GetString("TERMS_TO_SEARCH"), ";")
	cfg.App.FilesToSearch = strings.Split(viper.GetString("FILES_TO_SEARCH"), ";")
	cfg.App.TermsFilesToCount = strings.Split(viper.GetString("TERMS_FILES_TO_COUNT"), ";")
	cfg.App.DefaultCloneDir = viper.GetString("DIR")
	cfg.App.DestDir = viper.GetString("DEST_DIR")
	cfg.App.DevelopersMap = viper.GetString("DEVELOPERS_MAP")
	cfg.App.ForbiddenFiles = strings.Split(viper.GetString("FORBIDDEN_FILES_TO_SEARCH"), ";")

	// Count thresholds
	cfg.App.CountThresholdLow = viper.GetInt("COUNT_THRESHOLD_LOW")
	cfg.App.CountThresholdMedium = viper.GetInt("COUNT_THRESHOLD_MEDIUM")

	// Bitbucket Configuration
	cfg.Bitbucket.Token = viper.GetString("BITBUCKET_TOKEN")
	cfg.Bitbucket.Username = viper.GetString("BITBUCKET_USERNAME")
	cfg.Bitbucket.Workspace = viper.GetString("BITBUCKET_WORKSPACE")

	// AWS Configuration
	cfg.AWS.AccessKeyID = viper.GetString("AWS_ACCESS_KEY_ID")
	cfg.AWS.SecretAccessKey = viper.GetString("AWS_SECRET_ACCESS_KEY")
	cfg.AWS.Region = viper.GetString("AWS_REGION")
	cfg.AWS.BucketName = viper.GetString("AWS_BUCKET_NAME")
	cfg.AWS.AWSUploadPath = viper.GetString("AWS_PATH")

	// JIRA configuration
	cfg.App.JiraBaseURL = viper.GetString("JIRA_BASE_URL")
	cfg.App.JiraTaskEnabled = viper.GetBool("JIRA_TASK_ENABLED")
	cfg.App.JiraParentTask = viper.GetString("JIRA_PARENT_TASK")
	cfg.App.JiraTitleTemplate = viper.GetString("JIRA_TITLE_TEMPLATE")
	cfg.App.JiraDescTemplate = viper.GetString("JIRA_DESC_TEMPLATE")
	cfg.App.JiraDocLinks = strings.Split(viper.GetString("JIRA_DOC_LINKS"), ";")
	cfg.App.JiraDocLinks = utils.FilterEmpty(cfg.App.JiraDocLinks)
	cfg.App.JiraProjectKey = viper.GetString("JIRA_PROJECT_KEY")
	cfg.App.JiraIssueType = viper.GetString("JIRA_ISSUE_TYPE")
	cfg.App.JiraUsername = viper.GetString("JIRA_USERNAME")
	cfg.App.JiraAPIToken = viper.GetString("JIRA_API_TOKEN")
	// Set default values for JIRA templates if not provided
	if cfg.App.JiraTitleTemplate == "" {
		cfg.App.JiraTitleTemplate = "Amélioration de la CI/CD pour le projet {{.RepoName}}"
	}
	if cfg.App.JiraDescTemplate == "" {
		cfg.App.JiraDescTemplate = "Ajouter les éléments suivants dans le projet :\n{{.MissingElements}}\n\nJIRA parente: {{.ParentTask}}\nAssigné: {{.TopDeveloper}}"
	}

	// Filter empty values
	cfg.App.DefaultColumns = utils.FilterEmpty(cfg.App.DefaultColumns)
	cfg.App.TermsToSearch = utils.FilterEmpty(cfg.App.TermsToSearch)
	cfg.App.FilesToSearch = utils.FilterEmpty(cfg.App.FilesToSearch)
	cfg.App.TermsFilesToCount = utils.FilterEmpty(cfg.App.TermsFilesToCount)
	cfg.App.ForbiddenFiles = utils.FilterEmpty(cfg.App.ForbiddenFiles)
}

// Get returns the configuration instance
func Get() *Config {
	return cfg
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}
