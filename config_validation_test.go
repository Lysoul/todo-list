package main

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvExampleFile(t *testing.T) {
	t.Run("env.example exists", func(t *testing.T) {
		_, err := os.Stat(".env.example")
		assert.NoError(t, err, ".env.example file should exist")
	})

	t.Run("env.example is readable", func(t *testing.T) {
		file, err := os.Open(".env.example")
		require.NoError(t, err, "should be able to open .env.example")
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineCount := 0
		for scanner.Scan() {
			lineCount++
		}
		assert.NoError(t, scanner.Err(), "should be able to read file")
		assert.Greater(t, lineCount, 0, "file should not be empty")
	})

	t.Run("env.example contains required keys", func(t *testing.T) {
		requiredKeys := []string{
			"HTTP_PORT",
			"HTTP_PATH_PREFIX",
			"HTTP_ENABLE_CORS",
			"SERVICE_NAME",
			"POSTGRES_URL",
		}

		content, err := os.ReadFile(".env.example")
		require.NoError(t, err, "should be able to read .env.example")

		fileContent := string(content)
		for _, key := range requiredKeys {
			assert.Contains(t, fileContent, key,
				"config should contain required key: %s", key)
		}
	})

	t.Run("env.example has valid key-value format", func(t *testing.T) {
		file, err := os.Open(".env.example")
		require.NoError(t, err)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := strings.TrimSpace(scanner.Text())
			
			// Skip empty lines and comments
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			// Check for key=value format
			parts := strings.SplitN(line, "=", 2)
			assert.Len(t, parts, 2,
				"line %d should be in KEY=VALUE format: %s", lineNum, line)
			
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				assert.NotEmpty(t, key,
					"line %d should have non-empty key", lineNum)
			}
		}
		assert.NoError(t, scanner.Err())
	})

	t.Run("HTTP_PORT has valid value", func(t *testing.T) {
		content, err := os.ReadFile(".env.example")
		require.NoError(t, err)

		found := false
		scanner := bufio.NewScanner(strings.NewReader(string(content)))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "HTTP_PORT=") {
				found = true
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					port := strings.TrimSpace(parts[1])
					assert.NotEmpty(t, port, "HTTP_PORT should have a value")
					// Verify it's a number
					assert.Regexp(t, `^\d+$`, port,
						"HTTP_PORT should be a valid port number")
				}
			}
		}
		assert.True(t, found, "HTTP_PORT should be defined")
	})

	t.Run("POSTGRES_URL has valid format", func(t *testing.T) {
		content, err := os.ReadFile(".env.example")
		require.NoError(t, err)

		found := false
		scanner := bufio.NewScanner(strings.NewReader(string(content)))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "POSTGRES_URL=") {
				found = true
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					url := strings.TrimSpace(parts[1])
					assert.NotEmpty(t, url, "POSTGRES_URL should have a value")
					assert.Contains(t, url, "postgres://",
						"POSTGRES_URL should start with postgres://")
				}
			}
		}
		assert.True(t, found, "POSTGRES_URL should be defined")
	})

	t.Run("boolean flags have valid values", func(t *testing.T) {
		content, err := os.ReadFile(".env.example")
		require.NoError(t, err)

		booleanKeys := []string{
			"HTTP_ENABLE_CORS",
			"TRACE_INSECURE",
			"POSTGRES_DEBUG",
			"POSTGRES_MIGRATE",
		}

		fileContent := string(content)
		for _, key := range booleanKeys {
			scanner := bufio.NewScanner(strings.NewReader(fileContent))
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if strings.HasPrefix(line, key+"=") {
					parts := strings.SplitN(line, "=", 2)
					if len(parts) == 2 {
						value := strings.ToLower(strings.TrimSpace(parts[1]))
						assert.Contains(t, []string{"true", "false"}, value,
							"%s should be true or false", key)
					}
				}
			}
		}
	})
}

func TestEnvLocalFile(t *testing.T) {
	t.Run("env.local exists", func(t *testing.T) {
		_, err := os.Stat(".env.local")
		assert.NoError(t, err, ".env.local file should exist")
	})

	t.Run("env.local and env.example have same structure", func(t *testing.T) {
		exampleContent, err := os.ReadFile(".env.example")
		require.NoError(t, err)

		localContent, err := os.ReadFile(".env.local")
		require.NoError(t, err)

		// Extract keys from both files
		exampleKeys := extractKeys(string(exampleContent))
		localKeys := extractKeys(string(localContent))

		assert.ElementsMatch(t, exampleKeys, localKeys,
			".env.local should have same keys as .env.example")
	})
}

func TestVSCodeSettings(t *testing.T) {
	t.Run("vscode settings file exists", func(t *testing.T) {
		_, err := os.Stat(".vscode/settings.json")
		assert.NoError(t, err, ".vscode/settings.json should exist")
	})

	t.Run("vscode settings is valid JSON", func(t *testing.T) {
		content, err := os.ReadFile(".vscode/settings.json")
		require.NoError(t, err)

		// Basic JSON validation
		contentStr := string(content)
		assert.True(t, strings.HasPrefix(contentStr, "{"),
			"JSON should start with {")
		assert.True(t, strings.HasSuffix(strings.TrimSpace(contentStr), "}"),
			"JSON should end with }")
	})

	t.Run("vscode settings contains mermaid configuration", func(t *testing.T) {
		content, err := os.ReadFile(".vscode/settings.json")
		require.NoError(t, err)

		contentStr := string(content)
		assert.Contains(t, contentStr, "markdown-mermaid",
			"settings should contain mermaid configuration")
	})
}

func TestJustfile(t *testing.T) {
	t.Run("justfile exists", func(t *testing.T) {
		_, err := os.Stat("justfile")
		assert.NoError(t, err, "justfile should exist")
	})

	t.Run("justfile is readable", func(t *testing.T) {
		content, err := os.ReadFile("justfile")
		assert.NoError(t, err, "justfile should be readable")
		assert.NotEmpty(t, content, "justfile should not be empty")
	})

	t.Run("justfile contains required recipes", func(t *testing.T) {
		content, err := os.ReadFile("justfile")
		require.NoError(t, err)

		contentStr := string(content)
		requiredRecipes := []string{"dbinit", "cleango"}

		for _, recipe := range requiredRecipes {
			assert.Contains(t, contentStr, recipe,
				"justfile should contain recipe: %s", recipe)
		}

		// Also check for "env" recipe which has a parameter
		assert.Contains(t, contentStr, "env ",
			"justfile should contain env recipe")
	})

	t.Run("justfile has dotenv-load setting", func(t *testing.T) {
		content, err := os.ReadFile("justfile")
		require.NoError(t, err)

		assert.Contains(t, string(content), "set dotenv-load",
			"justfile should enable dotenv loading")
	})
}
// Helper function to extract keys from env file content
func extractKeys(content string) []string {
	var keys []string
	scanner := bufio.NewScanner(strings.NewReader(content))
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			if key != "" {
				keys = append(keys, key)
			}
		}
	}
	
	return keys
}

func TestConfigurationConsistency(t *testing.T) {
	t.Run("all environment files are consistent", func(t *testing.T) {
		exampleContent, err := os.ReadFile(".env.example")
		if err != nil {
			t.Skip(".env.example not found")
		}

		localContent, err := os.ReadFile(".env.local")
		if err != nil {
			t.Skip(".env.local not found")
		}

		exampleKeys := extractKeys(string(exampleContent))
		localKeys := extractKeys(string(localContent))

		// Check for missing keys
		for _, key := range exampleKeys {
			assert.Contains(t, localKeys, key,
				"key %s from .env.example should exist in .env.local", key)
		}

		for _, key := range localKeys {
			assert.Contains(t, exampleKeys, key,
				"key %s from .env.local should exist in .env.example", key)
		}
	})

	t.Run("required configuration groups are present", func(t *testing.T) {
		content, err := os.ReadFile(".env.example")
		require.NoError(t, err)

		contentStr := string(content)
		
		// Check for HTTP configuration group
		httpKeys := []string{"HTTP_PORT", "HTTP_PATH_PREFIX", "HTTP_ENABLE_CORS"}
		for _, key := range httpKeys {
			assert.Contains(t, contentStr, key,
				"HTTP configuration should include %s", key)
		}

		// Check for Postgres configuration group
		postgresKeys := []string{"POSTGRES_URL", "POSTGRES_DEBUG", "POSTGRES_MIGRATE"}
		for _, key := range postgresKeys {
			assert.Contains(t, contentStr, key,
				"Postgres configuration should include %s", key)
		}

		// Check for service configuration group
		serviceKeys := []string{"SERVICE_NAME", "TRACE_COLLECTOR_URL"}
		for _, key := range serviceKeys {
			assert.Contains(t, contentStr, key,
				"Service configuration should include %s", key)
		}
	})
}

// Benchmark tests
func BenchmarkEnvFileRead(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = os.ReadFile(".env.example")
	}
}

func BenchmarkExtractKeys(b *testing.B) {
	content, _ := os.ReadFile(".env.example")
	contentStr := string(content)
	
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = extractKeys(contentStr)
	}
}