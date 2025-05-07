package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/s3pweb/gitArchiveS3Report/config"
	"github.com/s3pweb/gitArchiveS3Report/utils/jira"
	"github.com/spf13/cobra"
)

var (
	port int
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a local server for JIRA integration",
	Long:  `Start a local HTTP server that provides an API for JIRA ticket creation.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()

		// Validate JIRA configuration
		if !cfg.App.JiraTaskEnabled {
			fmt.Println("JIRA task creation is disabled in configuration")
			os.Exit(1)
		}

		if cfg.App.JiraAPIToken == "" || cfg.App.JiraUsername == "" {
			fmt.Println("JIRA API credentials are not configured")
			os.Exit(1)
		}

		fmt.Printf("Starting local server on port %d...\n", port)
		fmt.Println("Server will handle JIRA ticket creation requests")
		fmt.Println("Use Ctrl+C to stop the server")

		http.HandleFunc("/create-jira-ticket", handleCreateTicket)
		http.HandleFunc("/", handleIndex)

		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			fmt.Printf("Error starting server: %v\n", err)
			os.Exit(1)
		}
	},
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Git Archive S3 Report - JIRA Integration</title>
			<style>
				body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
				h1 { color: #2c3e50; }
				.info { background-color: #f8f9fa; padding: 15px; border-radius: 5px; }
				pre { background-color: #f1f1f1; padding: 10px; border-radius: 5px; overflow-x: auto; }
			</style>
		</head>
		<body>
			<h1>Git Archive S3 Report - JIRA Integration</h1>
			<div class="info">
				<p>This service provides an API for creating JIRA tickets from the Git Archive S3 Report tool.</p>
				<p>The server is running properly and ready to handle requests.</p>
				<h2>API Usage</h2>
				<p>To create a JIRA ticket, send a POST request to:</p>
				<pre>/create-jira-ticket</pre>
				<p>With the following parameters:</p>
				<ul>
					<li><strong>title</strong>: Title of the JIRA ticket</li>
					<li><strong>description</strong>: Description of the JIRA ticket</li>
					<li><strong>assignee</strong>: (Optional) Username to assign the ticket to</li>
					<li><strong>parent</strong>: (Optional) Parent ticket key</li>
				</ul>
			</div>
		</body>
		</html>
	`))
}

func handleCreateTicket(w http.ResponseWriter, r *http.Request) {
	// Accept GET and POST
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Deserialize JSON body
	var requestData struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Assignee    string `json:"assignee"`
		Parent      string `json:"parent"`
	}

	var err error
	// Process POST or JSON parameters
	paramMap := map[string]*string{
		"title":       &requestData.Title,
		"description": &requestData.Description,
		"assignee":    &requestData.Assignee,
		"parent":      &requestData.Parent,
	}

	for key, value := range paramMap {
		rawValue := r.URL.Query().Get(key)
		decodedValue, err := url.QueryUnescape(rawValue)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error decoding parameter %s: %v", key, err), http.StatusBadRequest)
			return
		}
		*value = decodedValue
	}

	// Check required parameters
	if requestData.Title == "" {
		http.Error(w, "Missing required parameter: title", http.StatusBadRequest)
		return
	}

	// If parent is not specified, use the default parent
	if requestData.Parent == "" {
		cfg := config.Get()
		requestData.Parent = cfg.App.JiraParentTask
	}

	// Create the JIRA ticket
	ticketURL, err := jira.CreateJiraTicket(
		requestData.Title,
		requestData.Description,
		requestData.Assignee,
		requestData.Parent,
	)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating JIRA ticket: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the URL of the created ticket
	response := struct {
		Success bool   `json:"success"`
		URL     string `json:"url"`
		Message string `json:"message"`
	}{
		Success: true,
		URL:     ticketURL,
		Message: "JIRA ticket created successfully",
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", 8081, "Port to run the server on")
	rootCmd.AddCommand(serveCmd)
}
