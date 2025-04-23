// utils/excel/jira_integration.go

package excel

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"text/template"

	"github.com/s3pweb/gitArchiveS3Report/config"
	"github.com/s3pweb/gitArchiveS3Report/utils/structs"
	"github.com/xuri/excelize/v2"
)

// JiraTaskData holds the data for JIRA task template rendering
type JiraTaskData struct {
	RepoName        string
	BranchName      string
	TopDeveloper    string
	LastDeveloper   string
	MissingElements string
	ParentTask      string
}

// AddJiraButtons adds a JIRA task creation button to each row in the Excel sheet
func AddJiraButtons(f *excelize.File, sheet string, branchesInfo []structs.BranchInfo) error {
	cfg := config.Get()

	// Skip if JIRA task creation is disabled or URL is not set
	if !cfg.App.JiraTaskEnabled || cfg.App.JiraBaseURL == "" {
		return nil
	}

	// Add JIRA column header
	headerStyle, err := f.NewStyle(`{
		"fill": {
			"type": "pattern",
			"color": ["#2c3e50"],
			"pattern": 1
		},
		"font": {
			"bold": true,
			"color": "#ffffff"
		},
		"alignment": {
			"horizontal": "center",
			"vertical": "center"
		},
		"border": [
			{
				"type": "left",
				"color": "#000000",
				"style": 1
			},
			{
				"type": "top",
				"color": "#000000",
				"style": 1
			},
			{
				"type": "bottom",
				"color": "#000000",
				"style": 1
			},
			{
				"type": "right",
				"color": "#000000",
				"style": 1
			}
		]
	}`)
	if err != nil {
		return err
	}

	// Determine the last column
	columns := cfg.App.DefaultColumns
	columns = append(columns, cfg.App.FilesToSearch...)
	columns = append(columns, cfg.App.TermsToSearch...)
	columns = append(columns, cfg.App.ForbiddenFiles...)
	lastColumn := rune('A' + len(columns))

	// Add JIRA column header
	jiraCell := fmt.Sprintf("%c1", lastColumn)
	f.SetCellValue(sheet, jiraCell, "CREATE JIRA TASK")
	f.SetCellStyle(sheet, jiraCell, jiraCell, headerStyle)
	f.SetColWidth(sheet, string(lastColumn), string(lastColumn), 25)

	// Pre-compile templates
	titleTmpl, err := template.New("title").Parse(cfg.App.JiraTitleTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse JIRA title template: %v", err)
	}

	descTmpl, err := template.New("desc").Parse(cfg.App.JiraDescTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse JIRA description template: %v", err)
	}

	// Add JIRA button to each row
	for i, branchInfo := range branchesInfo {
		row := i + 2 // +2 because row 1 is the header and we start at 0 index

		// Determine missing elements
		var elementsToAdd []string
		var filesToRemove []string

		// Check missing required files
		for file, exists := range branchInfo.FilesToSearch {
			if !exists {
				cleanFile := cleanRegexPattern(file)
				elementsToAdd = append(elementsToAdd, cleanFile)
			}
		}

		// Check missing required terms
		for term, exists := range branchInfo.TermsToSearch {
			if !exists {
				cleanTerm := cleanRegexPattern(term)
				elementsToAdd = append(elementsToAdd, cleanTerm)
			}
		}

		// Check forbidden files that exist
		for file, exists := range branchInfo.ForbiddenFiles {
			if exists {
				cleanFile := cleanRegexPattern(file)
				filesToRemove = append(filesToRemove, cleanFile)
			}
		}

		// If there are missing elements or forbidden files, create a JIRA task link
		if len(elementsToAdd) > 0 || len(filesToRemove) > 0 {
			// Build a well-formatted description with clear sections
			var descriptionBuilder strings.Builder

			// Add section for elements to add
			if len(elementsToAdd) > 0 {
				descriptionBuilder.WriteString("Éléments à ajouter :\n")
				for _, element := range elementsToAdd {
					descriptionBuilder.WriteString("- " + element + "\n")
				}
				descriptionBuilder.WriteString("\n")
			}

			// Add section for files to remove
			if len(filesToRemove) > 0 {
				descriptionBuilder.WriteString("Fichiers à supprimer :\n")
				for _, file := range filesToRemove {
					descriptionBuilder.WriteString("- " + file + "\n")
				}
				descriptionBuilder.WriteString("\n")
			}

			// Add additional documentation links
			if len(cfg.App.JiraDocLinks) > 0 {
				descriptionBuilder.WriteString("\n\nDocumentation : \n")
				for _, link := range cfg.App.JiraDocLinks {
					// Split the link into text and URL if it contains a "|"
					parts := strings.SplitN(link, "|", 2)
					if len(parts) == 2 {
						// Use the first part as the text and the second part as the URL
						descriptionBuilder.WriteString("- [" + parts[0] + "|" + parts[1] + "]\n")
					} else {
						// If no "|" is found, treat the whole string as a URL
						descriptionBuilder.WriteString("- " + link + "\n")
					}
				}
			}

			// Add parent JIRA reference
			descriptionBuilder.WriteString("JIRA parente: " + cfg.App.JiraParentTask)

			// Use this formatted description in the template data
			data := JiraTaskData{
				RepoName:        branchInfo.RepoName,
				BranchName:      branchInfo.BranchName,
				TopDeveloper:    branchInfo.TopDeveloper,
				LastDeveloper:   branchInfo.LastDeveloper,
				MissingElements: descriptionBuilder.String(),
				ParentTask:      cfg.App.JiraParentTask,
			}

			// Render title template
			var titleBuf bytes.Buffer
			if err := titleTmpl.Execute(&titleBuf, data); err != nil {
				return fmt.Errorf("failed to render JIRA title template: %v", err)
			}
			title := titleBuf.String()

			// Render description template
			var descBuf bytes.Buffer
			if err := descTmpl.Execute(&descBuf, data); err != nil {
				return fmt.Errorf("failed to render JIRA description template: %v", err)
			}
			description := descBuf.String()

			// Add button to the cell
			cell := fmt.Sprintf("%c%d", lastColumn, row)

			// If the JIRA API token and username are set, create a link to our local server
			if cfg.App.JiraAPIToken != "" && cfg.App.JiraUsername != "" {
				// Create a link to the local server
				jiraLink := fmt.Sprintf("http://localhost:8081/create-jira-ticket?title=%s&description=%s&assignee=%s&parent=%s",
					url.QueryEscape(title),
					url.QueryEscape(description),
					url.QueryEscape(branchInfo.TopDeveloper),
					url.QueryEscape(cfg.App.JiraParentTask))

				// Create a hyperlink in the cell
				f.SetCellHyperLink(sheet, cell, jiraLink, "External")
				f.SetCellValue(sheet, cell, "Create JIRA Task")
			} else {
				// Fallback to the JIRA base URL (if not using local server)
				jiraLink := fmt.Sprintf("%s/secure/CreateIssue.jspa?summary=%s&description=%s",
					cfg.App.JiraBaseURL,
					url.QueryEscape(title),
					url.QueryEscape(description))

				f.SetCellHyperLink(sheet, cell, jiraLink, "External")
				f.SetCellValue(sheet, cell, "Create JIRA Task")
			}

			// Style the button
			buttonStyle, err := f.NewStyle(`{
				"font": {
					"color": "#0563C1",
					"underline": "single"
				},
				"alignment": {
					"horizontal": "center",
					"vertical": "center"
				},
				"border": [
					{
						"type": "left",
						"color": "#000000",
						"style": 1
					},
					{
						"type": "top",
						"color": "#000000",
						"style": 1
					},
					{
						"type": "bottom",
						"color": "#000000",
						"style": 1
					},
					{
						"type": "right",
						"color": "#000000",
						"style": 1
					}
				]
			}`)
			if err == nil {
				f.SetCellStyle(sheet, cell, cell, buttonStyle)
			}
		}
	}

	return nil
}

// cleanRegexPattern cleans up the regex pattern for better readability
func cleanRegexPattern(pattern string) string {
	pattern = strings.Replace(pattern, "(?i)", "", -1)
	pattern = strings.Replace(pattern, "$", "", -1)
	pattern = strings.Replace(pattern, "(-\\w+)?\\.", ".", -1)

	return pattern
}
