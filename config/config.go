package config

import (
	"fmt"
	"os"
	"strings"

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
}

type LoggerConfig struct {
	Level string
}

type AppConfig struct {
	CPU             int
	DevelopersMap   string
	DefaultColumns  []string
	TermsToSearch   []string
	FilesToSearch   []string
	DefaultCloneDir string
	MainBranchOnly  bool
	ShallowClone    bool
}

// Init initializes the configuration
func Init() error {
	// Default configuration
	setDefaults()

	// Configuration via file (optional)
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// The config file is optional
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Environment variables (override file)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Mapping environment variables
	mapEnvVariables()

	// Loading the configuration
	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("unable to decode config: %w", err)
	}

	// Validating the configuration
	return validateConfig(cfg)
}

func setDefaults() {
	// App defaults
	viper.SetDefault("app.defaultCloneDir", "./repositories")
	viper.SetDefault("app.cpu", 1)
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("app.defaultColumns", []string{
		"RepoName", "BranchName", "LastCommitDate", "TimeSinceLastCommit",
		"Commitnbr", "HostLine", "LastDeveloper", "LastDeveloperPercentage",
	})
	viper.SetDefault("app.termsToSearch", []string{"vaumt", "swagger"})
	viper.SetDefault("app.filesToSearch", []string{
		"(?i)sonar-project.properties$",
		"(?i)bitbucket-pipelines.yml$",
		"(?i)Dockerfile$",
		"(?i)docker-compose(-\\w+)?\\.yaml$",
	})
	viper.SetDefault("app.mainBranchOnly", false)
	viper.SetDefault("app.shallowClone", false)

}

func mapEnvVariables() {
	// Explicit mapping of environment variables
	envMappings := map[string]string{
		"BITBUCKET_TOKEN":       "bitbucket.token",
		"BITBUCKET_USERNAME":    "bitbucket.username",
		"BITBUCKET_WORKSPACE":   "bitbucket.workspace",
		"AWS_ACCESS_KEY_ID":     "aws.accessKeyID",
		"AWS_SECRET_ACCESS_KEY": "aws.secretAccessKey",
		"AWS_REGION":            "aws.region",
		"AWS_BUCKET_NAME":       "aws.bucketName",
		"LOG_LEVEL":             "logger.level",
		"APP_CPU":               "app.cpu",
		"APP_CLONE_DIR":         "app.defaultCloneDir",
		"APP_MAIN_BRANCH_ONLY":  "app.mainBranchOnly",
	}

	for env, path := range envMappings {
		if value := os.Getenv(env); value != "" {
			viper.Set(path, value)
		}
	}
}

func validateConfig(cfg *Config) error {
	if cfg.Bitbucket.Token == "" {
		return fmt.Errorf("bitbucket token is required (BITBUCKET_TOKEN)")
	}
	if cfg.Bitbucket.Username == "" {
		return fmt.Errorf("bitbucket username is required (BITBUCKET_USERNAME)")
	}
	if cfg.Bitbucket.Workspace == "" {
		return fmt.Errorf("bitbucket workspace is required (BITBUCKET_WORKSPACE)")
	}
	return nil
}

// Get returns the configuration instance
func Get() *Config {
	return cfg
}
