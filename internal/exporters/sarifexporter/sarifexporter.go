// SARIF (Static Analysis Results Interchange Format)
// Based on https://docs.oasis-open.org/sarif/sarif/v2.1.0/sarif-v2.1.0.pdf

package sarifexporter

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mlw157/scout/internal/models"
)

type SARIFExporter struct {
	OutputFile string
}

func NewSARIFExporter(outputFile string) *SARIFExporter {
	return &SARIFExporter{OutputFile: outputFile}
}

// 3.13 - https://docs.oasis-open.org/sarif/sarif/v2.1.0/errata01/os/sarif-v2.1.0-errata01-os-complete.html#_Toc141790728
type SARIFLog struct {
	Schema  string `json:"$schema,omitempty"`
	Version string `json:"version"`
	Runs    []Run  `json:"runs"`
}

// Run object
// 3.14 - https://docs.oasis-open.org/sarif/sarif/v2.1.0/errata01/os/sarif-v2.1.0-errata01-os-complete.html#_Toc141790734
type Run struct {
	Tool    Tool     `json:"tool"`
	Results []Result `json:"results,omitempty"`
}

// Tool object
// 3.18 - https://docs.oasis-open.org/sarif/sarif/v2.1.0/errata01/os/sarif-v2.1.0-errata01-os-complete.html#_Toc141790779
type Tool struct {
	Driver ToolComponent `json:"driver"`
}

// ToolComponent object
// 3.19 - https://docs.oasis-open.org/sarif/sarif/v2.1.0/errata01/os/sarif-v2.1.0-errata01-os-complete.html#_Toc141790783
// used for "driver"
type ToolComponent struct {
	Name           string                `json:"name"`
	Version        string                `json:"version,omitempty"`
	InformationUri string                `json:"informationUri,omitempty"`
	Rules          []ReportingDescriptor `json:"rules,omitempty"`
}

// ReportingDescriptor object
// 3.49 - https://docs.oasis-open.org/sarif/sarif/v2.1.0/errata01/os/sarif-v2.1.0-errata01-os-complete.html#_Toc141791086
// defines a rule
type ReportingDescriptor struct {
	ID               string                    `json:"id"`
	Name             string                    `json:"name,omitempty"`
	ShortDescription *MultiformatMessageString `json:"shortDescription,omitempty"`
	FullDescription  *MultiformatMessageString `json:"fullDescription,omitempty"`
	HelpUri          string                    `json:"helpUri,omitempty"`
	DefaultConfig    *ReportingConfiguration   `json:"defaultConfiguration,omitempty"`
}

// ReportingConfiguration object
// 3.52 - https://docs.oasis-open.org/sarif/sarif/v2.1.0/errata01/os/sarif-v2.1.0-errata01-os-complete.html#_Toc141791112
type ReportingConfiguration struct {
	Level string `json:"level,omitempty"`
}

// MultiformatMessageString object
// 3.12 - https://docs.oasis-open.org/sarif/sarif/v2.1.0/errata01/os/sarif-v2.1.0-errata01-os-complete.html#_Toc141790723
type MultiformatMessageString struct {
	Text string `json:"text"`
}

// Result object
// 3.27 - https://docs.oasis-open.org/sarif/sarif/v2.1.0/errata01/os/sarif-v2.1.0-errata01-os-complete.html#_Toc141790888
type Result struct {
	RuleID    string     `json:"ruleId,omitempty"`
	Level     string     `json:"level,omitempty"`
	Message   Message    `json:"message"`
	Locations []Location `json:"locations,omitempty"`
}

// Message object
// 3.11 - https://docs.oasis-open.org/sarif/sarif/v2.1.0/errata01/os/sarif-v2.1.0-errata01-os-complete.html#_Toc141790709
type Message struct {
	Text string `json:"text,omitempty"`
}

// Location object
// 3.28 - https://docs.oasis-open.org/sarif/sarif/v2.1.0/errata01/os/sarif-v2.1.0-errata01-os-complete.html#_Toc141790920
type Location struct {
	PhysicalLocation *PhysicalLocation `json:"physicalLocation,omitempty"`
}

// PhysicalLocation object
// 3.29 - https://docs.oasis-open.org/sarif/sarif/v2.1.0/errata01/os/sarif-v2.1.0-errata01-os-complete.html#_Toc141790928
type PhysicalLocation struct {
	ArtifactLocation *ArtifactLocation `json:"artifactLocation,omitempty"`
}

// ArtifactLocation object
// 3.4 - https://docs.oasis-open.org/sarif/sarif/v2.1.0/errata01/os/sarif-v2.1.0-errata01-os-complete.html#_Toc141790677
type ArtifactLocation struct {
	Uri string `json:"uri,omitempty"`
}

func (s *SARIFExporter) Export(results []*models.ScanResult) error {
	file, err := os.Create(s.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", s.OutputFile, err)
	}
	defer file.Close()

	rulesMap := make(map[string]ReportingDescriptor)
	var sarifResults []Result

	for _, result := range results {
		for _, vuln := range result.Vulnerabilities {
			ruleID := vuln.CVE
			if ruleID == "" {
				ruleID = fmt.Sprintf("SCOUT-%s", vuln.Dependency.Name)
			}

			if _, exists := rulesMap[ruleID]; !exists {
				rulesMap[ruleID] = ReportingDescriptor{
					ID:               ruleID,
					Name:             vuln.Dependency.Name,
					ShortDescription: &MultiformatMessageString{Text: vuln.Summary},
					FullDescription:  &MultiformatMessageString{Text: vuln.Description},
					HelpUri:          vuln.URL,
					DefaultConfig:    &ReportingConfiguration{Level: mapSeverityToLevel(vuln.Severity)},
				}
			}

			level := mapSeverityToLevel(vuln.Severity)
			sarifResults = append(sarifResults, Result{
				RuleID: ruleID,
				Level:  level,
				Message: Message{
					Text: fmt.Sprintf("%s@%s is vulnerable: %s",
						vuln.Dependency.Name,
						vuln.Dependency.Version,
						vuln.Summary),
				},
				Locations: []Location{
					{
						PhysicalLocation: &PhysicalLocation{
							ArtifactLocation: &ArtifactLocation{
								Uri: result.SourceFile,
							},
						},
					},
				},
			})
		}
	}

	var rules []ReportingDescriptor
	for _, rule := range rulesMap {
		rules = append(rules, rule)
	}

	sarifReport := SARIFLog{
		Schema:  "https://docs.oasis-open.org/sarif/sarif/v2.1.0/errata01/os/schemas/sarif-schema-2.1.0.json",
		Version: "2.1.0",
		Runs: []Run{
			{
				Tool: Tool{
					Driver: ToolComponent{
						Name:           "Scout",
						Version:        "0.1.0",
						InformationUri: "https://github.com/mlw157/scout",
						Rules:          rules,
					},
				},
				Results: sarifResults,
			},
		},
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(sarifReport); err != nil {
		return fmt.Errorf("failed to encode results to SARIF format: %v", err)
	}

	log.Printf("Vulnerabilities exported to %s in SARIF format\n", s.OutputFile)
	return nil
}

func mapSeverityToLevel(severity string) string {
	switch strings.ToLower(severity) {
	case "critical", "high":
		return "error"
	case "medium", "moderate":
		return "warning"
	case "low":
		return "note"
	default:
		return "warning"
	}
}
