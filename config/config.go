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
func Init() {
	// Environment variables (override file)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Mapping environment variables
	mapEnvVariables()

	// Loading the configuration
	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		fmt.Printf("unable to decode config: %v\n", err)
	}
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

// Get returns the configuration instance
func Get() *Config {
	return cfg
}
