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
	UploadKey       string
}

type LoggerConfig struct {
	Level string
}

type AppConfig struct {
	CPU            int
	DevelopersMap  string
	DefaultColumns []string
	TermsToSearch  []string
	FilesToSearch  []string
	Dir            string
	DestDir        string
	MainBranchOnly bool
	ShallowClone   bool
	DevSheets      bool
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
	cfg.App.Dir = viper.GetString("DIR")
	cfg.App.DevelopersMap = viper.GetString("DEVELOPERS_MAP")

	// Bitbucket Configuration
	cfg.Bitbucket.Token = viper.GetString("BITBUCKET_TOKEN")
	cfg.Bitbucket.Username = viper.GetString("BITBUCKET_USERNAME")
	cfg.Bitbucket.Workspace = viper.GetString("BITBUCKET_WORKSPACE")

	// AWS Configuration
	cfg.AWS.AccessKeyID = viper.GetString("AWS_ACCESS_KEY_ID")
	cfg.AWS.SecretAccessKey = viper.GetString("AWS_SECRET_ACCESS_KEY")
	cfg.AWS.Region = viper.GetString("AWS_REGION")
	cfg.AWS.BucketName = viper.GetString("AWS_BUCKET_NAME")
	cfg.AWS.UploadKey = viper.GetString("AWS_UPLOAD_KEY")

	// Filter empty values
	cfg.App.DefaultColumns = utils.FilterEmpty(cfg.App.DefaultColumns)
	cfg.App.TermsToSearch = utils.FilterEmpty(cfg.App.TermsToSearch)
	cfg.App.FilesToSearch = utils.FilterEmpty(cfg.App.FilesToSearch)

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
