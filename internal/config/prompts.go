package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// LocalizedPrompts contains all prompts for a specific language
type LocalizedPrompts struct {
	// GenerateSummary is the prompt for generating document summaries
	GenerateSummary string `yaml:"generate_summary" json:"generate_summary"`
	// GenerateSessionTitle is the prompt for generating session titles
	GenerateSessionTitle string `yaml:"generate_session_title" json:"generate_session_title"`
	// SystemPrompt is the main RAG system prompt for chat responses
	SystemPrompt string `yaml:"system_prompt" json:"system_prompt"`
	// FallbackPrompt is used when no relevant context is found
	FallbackPrompt string `yaml:"fallback_prompt" json:"fallback_prompt"`
	// FallbackResponse is the default response when fallback is triggered
	FallbackResponse string `yaml:"fallback_response" json:"fallback_response"`
	// RewriteSystem is the system prompt for query rewriting
	RewriteSystem string `yaml:"rewrite_system" json:"rewrite_system"`
	// RewriteUser is the user prompt template for query rewriting
	RewriteUser string `yaml:"rewrite_user" json:"rewrite_user"`
	// KeywordsExtraction is the prompt for extracting keywords from queries
	KeywordsExtraction string `yaml:"keywords_extraction" json:"keywords_extraction"`
	// KeywordsExtractionUser is the user prompt template for keyword extraction
	KeywordsExtractionUser string `yaml:"keywords_extraction_user" json:"keywords_extraction_user"`
	// ExtractEntities is the prompt for entity extraction
	ExtractEntities string `yaml:"extract_entities" json:"extract_entities"`
	// ExtractRelationships is the prompt for relationship extraction
	ExtractRelationships string `yaml:"extract_relationships" json:"extract_relationships"`
	// GenerateQuestions is the prompt for generating questions from document chunks
	GenerateQuestions string `yaml:"generate_questions" json:"generate_questions"`
	// SimplifyQuery is the system prompt for query simplification
	SimplifyQuery string `yaml:"simplify_query" json:"simplify_query"`
	// SimplifyQueryUser is the user prompt template for query simplification
	SimplifyQueryUser string `yaml:"simplify_query_user" json:"simplify_query_user"`
}

// PromptsConfig contains localized prompts for all supported languages
type PromptsConfig struct {
	// DefaultLanguage is the default language code (e.g., "ko-KR", "en-US", "zh-CN")
	DefaultLanguage string `yaml:"default_language" json:"default_language"`
	// Prompts maps language codes to their localized prompts
	Prompts map[string]*LocalizedPrompts `yaml:"prompts" json:"prompts"`
}

// SupportedLanguages defines the supported language codes
var SupportedLanguages = []string{"ko-KR", "en-US", "zh-CN"}

// GetPrompts returns the prompts for the specified language, falling back to default
func (pc *PromptsConfig) GetPrompts(lang string) *LocalizedPrompts {
	if pc == nil || pc.Prompts == nil {
		return nil
	}

	// Try exact match
	if prompts, ok := pc.Prompts[lang]; ok {
		return prompts
	}

	// Try language prefix match (e.g., "ko" matches "ko-KR")
	langPrefix := strings.Split(lang, "-")[0]
	for code, prompts := range pc.Prompts {
		if strings.HasPrefix(code, langPrefix) {
			return prompts
		}
	}

	// Fall back to default language
	if pc.DefaultLanguage != "" {
		if prompts, ok := pc.Prompts[pc.DefaultLanguage]; ok {
			return prompts
		}
	}

	// Fall back to first available in deterministic order (SupportedLanguages order)
	for _, supportedLang := range SupportedLanguages {
		if prompts, ok := pc.Prompts[supportedLang]; ok {
			return prompts
		}
	}

	return nil
}

// GetGenerateSummaryPrompt returns the document summary prompt for the specified language
// Falls back to config.yaml prompt if localized version is not available
func (c *Config) GetGenerateSummaryPrompt(lang string) string {
	if c.LocalizedPrompts != nil {
		if prompts := c.LocalizedPrompts.GetPrompts(lang); prompts != nil && prompts.GenerateSummary != "" {
			return prompts.GenerateSummary
		}
	}
	// Fallback to config.yaml
	if c.Conversation != nil {
		return c.Conversation.GenerateSummaryPrompt
	}
	return ""
}

// GetSystemPrompt returns the RAG system prompt for the specified language
// Falls back to config.yaml prompt if localized version is not available
func (c *Config) GetSystemPrompt(lang string) string {
	if c.LocalizedPrompts != nil {
		if prompts := c.LocalizedPrompts.GetPrompts(lang); prompts != nil && prompts.SystemPrompt != "" {
			return prompts.SystemPrompt
		}
	}
	// Fallback to config.yaml
	if c.Conversation != nil && c.Conversation.Summary != nil {
		return c.Conversation.Summary.Prompt
	}
	return ""
}

// GetFallbackPrompt returns the fallback prompt for the specified language
func (c *Config) GetFallbackPrompt(lang string) string {
	if c.LocalizedPrompts != nil {
		if prompts := c.LocalizedPrompts.GetPrompts(lang); prompts != nil && prompts.FallbackPrompt != "" {
			return prompts.FallbackPrompt
		}
	}
	if c.Conversation != nil {
		return c.Conversation.FallbackPrompt
	}
	return ""
}

// GetFallbackResponse returns the fallback response for the specified language
func (c *Config) GetFallbackResponse(lang string) string {
	if c.LocalizedPrompts != nil {
		if prompts := c.LocalizedPrompts.GetPrompts(lang); prompts != nil && prompts.FallbackResponse != "" {
			return prompts.FallbackResponse
		}
	}
	if c.Conversation != nil {
		return c.Conversation.FallbackResponse
	}
	return ""
}

// GetRewriteSystemPrompt returns the query rewrite system prompt for the specified language
func (c *Config) GetRewriteSystemPrompt(lang string) string {
	if c.LocalizedPrompts != nil {
		if prompts := c.LocalizedPrompts.GetPrompts(lang); prompts != nil && prompts.RewriteSystem != "" {
			return prompts.RewriteSystem
		}
	}
	if c.Conversation != nil {
		return c.Conversation.RewritePromptSystem
	}
	return ""
}

// GetRewriteUserPrompt returns the query rewrite user prompt for the specified language
func (c *Config) GetRewriteUserPrompt(lang string) string {
	if c.LocalizedPrompts != nil {
		if prompts := c.LocalizedPrompts.GetPrompts(lang); prompts != nil && prompts.RewriteUser != "" {
			return prompts.RewriteUser
		}
	}
	if c.Conversation != nil {
		return c.Conversation.RewritePromptUser
	}
	return ""
}

// GetKeywordsExtractionPrompt returns the keywords extraction prompt for the specified language
func (c *Config) GetKeywordsExtractionPrompt(lang string) string {
	if c.LocalizedPrompts != nil {
		if prompts := c.LocalizedPrompts.GetPrompts(lang); prompts != nil && prompts.KeywordsExtraction != "" {
			return prompts.KeywordsExtraction
		}
	}
	// Note: config.yaml uses different field name pattern
	return ""
}

// GetGenerateSessionTitlePrompt returns the session title generation prompt for the specified language
func (c *Config) GetGenerateSessionTitlePrompt(lang string) string {
	if c.LocalizedPrompts != nil {
		if prompts := c.LocalizedPrompts.GetPrompts(lang); prompts != nil && prompts.GenerateSessionTitle != "" {
			return prompts.GenerateSessionTitle
		}
	}
	if c.Conversation != nil {
		return c.Conversation.GenerateSessionTitlePrompt
	}
	return ""
}

// GetExtractEntitiesPrompt returns the entity extraction prompt for the specified language
func (c *Config) GetExtractEntitiesPrompt(lang string) string {
	if c.LocalizedPrompts != nil {
		if prompts := c.LocalizedPrompts.GetPrompts(lang); prompts != nil && prompts.ExtractEntities != "" {
			return prompts.ExtractEntities
		}
	}
	if c.Conversation != nil {
		return c.Conversation.ExtractEntitiesPrompt
	}
	return ""
}

// GetDefaultLanguage returns the default language code
func (c *Config) GetDefaultLanguage() string {
	if c.LocalizedPrompts != nil && c.LocalizedPrompts.DefaultLanguage != "" {
		return c.LocalizedPrompts.DefaultLanguage
	}
	return "ko-KR" // Default to Korean
}

// loadLocalizedPrompts loads prompts from the prompts directory
func loadLocalizedPrompts(configDir string) (*PromptsConfig, error) {
	promptsDir := filepath.Join(configDir, "prompts")

	// Check if directory exists
	if _, err := os.Stat(promptsDir); os.IsNotExist(err) {
		return nil, nil // Directory doesn't exist, return nil
	}

	config := &PromptsConfig{
		DefaultLanguage: "ko-KR", // Default to Korean
		Prompts:         make(map[string]*LocalizedPrompts),
	}

	// Load each language file
	for _, lang := range SupportedLanguages {
		filename := fmt.Sprintf("%s.yaml", lang)
		filePath := filepath.Join(promptsDir, filename)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			continue // File doesn't exist, skip
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", filename, err)
		}

		var prompts LocalizedPrompts
		if err := yaml.Unmarshal(data, &prompts); err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", filename, err)
		}

		config.Prompts[lang] = &prompts
	}

	// Load default language setting if exists
	defaultFile := filepath.Join(promptsDir, "default.yaml")
	if _, err := os.Stat(defaultFile); err == nil {
		data, err := os.ReadFile(defaultFile)
		if err != nil {
			// Log error but continue with default language - file exists but couldn't be read
			fmt.Fprintf(os.Stderr, "Warning: failed to read %s: %v\n", defaultFile, err)
		} else {
			var defaultConfig struct {
				DefaultLanguage string `yaml:"default_language"`
			}
			if err := yaml.Unmarshal(data, &defaultConfig); err != nil {
				// Log parse error but continue with default language
				fmt.Fprintf(os.Stderr, "Warning: failed to parse %s: %v\n", defaultFile, err)
			} else if defaultConfig.DefaultLanguage != "" {
				config.DefaultLanguage = defaultConfig.DefaultLanguage
			}
		}
	}

	return config, nil
}
