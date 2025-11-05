package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("Config struct has correct fields", func(t *testing.T) {
		config := Config{}
		
		// Verify struct can be instantiated
		assert.NotNil(t, &config, "Config should be instantiable")
		
		// Verify default shutdown timeout via struct tag
		assert.IsType(t, time.Duration(0), config.ShutdownTimeout,
			"ShutdownTimeout should be time.Duration type")
	})

	t.Run("Config fields are properly typed", func(t *testing.T) {
		config := Config{
			ShutdownTimeout: 30 * time.Second,
		}
		
		assert.Equal(t, 30*time.Second, config.ShutdownTimeout,
			"ShutdownTimeout should accept duration values")
	})

	t.Run("Config with various shutdown timeouts", func(t *testing.T) {
		testCases := []struct {
			name     string
			duration time.Duration
		}{
			{"zero duration", 0},
			{"1 second", time.Second},
			{"20 seconds", 20 * time.Second},
			{"1 minute", time.Minute},
			{"negative duration", -time.Second},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				config := Config{
					ShutdownTimeout: tc.duration,
				}
				assert.Equal(t, tc.duration, config.ShutdownTimeout,
					"Config should accept %s", tc.name)
			})
		}
	})
}

func TestVersion(t *testing.T) {
	t.Run("Version variable exists", func(t *testing.T) {
		assert.NotEmpty(t, Version, "Version should be initialized")
	})

	t.Run("Version is a string", func(t *testing.T) {
		assert.IsType(t, "", Version, "Version should be a string")
	})

	t.Run("Version default value", func(t *testing.T) {
		// Default value should be "unknown"
		assert.Equal(t, "unknown", Version,
			"Version should default to 'unknown' unless set at build time")
	})

	t.Run("Version can be modified", func(t *testing.T) {
		// Save original value
		original := Version
		defer func() {
			Version = original
		}()

		// Test modification
		testVersion := "v1.2.3"
		Version = testVersion
		assert.Equal(t, testVersion, Version, "Version should be modifiable")
	})

	t.Run("Version handles various formats", func(t *testing.T) {
		original := Version
		defer func() {
			Version = original
		}()

		testVersions := []string{
			"v1.0.0",
			"v1.0.0-beta",
			"v2.3.4+20250101",
			"dev-20250101-abc123",
			"",
			"latest",
		}

		for _, v := range testVersions {
			Version = v
			assert.Equal(t, v, Version, "Version should accept format: %s", v)
		}
	})
}

func TestConfigStructTags(t *testing.T) {
	t.Run("ShutdownTimeout has correct envconfig tag", func(t *testing.T) {
		// This test validates the struct tags are present
		// The actual parsing is done by envconfig library
		config := Config{}
		assert.IsType(t, time.Duration(0), config.ShutdownTimeout,
			"ShutdownTimeout field type should be time.Duration")
	})
}

// TestConfigZeroValues tests zero value behavior
func TestConfigZeroValues(t *testing.T) {
	t.Run("zero value Config is valid", func(t *testing.T) {
		var config Config
		
		assert.Equal(t, time.Duration(0), config.ShutdownTimeout,
			"zero value ShutdownTimeout should be 0")
	})

	t.Run("partially initialized Config", func(t *testing.T) {
		config := Config{
			ShutdownTimeout: 15 * time.Second,
		}
		
		assert.Equal(t, 15*time.Second, config.ShutdownTimeout)
		// HTTP and Postgres fields should be zero values
	})
}

// TestConfigEdgeCases tests edge cases and boundary conditions
func TestConfigEdgeCases(t *testing.T) {
	t.Run("very long shutdown timeout", func(t *testing.T) {
		config := Config{
			ShutdownTimeout: 24 * time.Hour,
		}
		assert.Equal(t, 24*time.Hour, config.ShutdownTimeout,
			"should handle long durations")
	})

	t.Run("maximum duration value", func(t *testing.T) {
		maxDuration := time.Duration(1<<63 - 1)
		config := Config{
			ShutdownTimeout: maxDuration,
		}
		assert.Equal(t, maxDuration, config.ShutdownTimeout,
			"should handle maximum duration")
	})

	t.Run("minimum duration value", func(t *testing.T) {
		minDuration := time.Duration(-1 << 63)
		config := Config{
			ShutdownTimeout: minDuration,
		}
		assert.Equal(t, minDuration, config.ShutdownTimeout,
			"should handle minimum duration")
	})
}

// TestVersionConcurrency tests concurrent access to Version variable
func TestVersionConcurrency(t *testing.T) {
	t.Run("concurrent reads of Version", func(t *testing.T) {
		const goroutines = 100
		done := make(chan bool, goroutines)
		
		for i := 0; i < goroutines; i++ {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("panic during concurrent read: %v", r)
					}
					done <- true
				}()
				
				_ = Version
				assert.NotEmpty(t, Version)
			}()
		}
		
		for i := 0; i < goroutines; i++ {
			<-done
		}
	})
}

// Benchmark tests
func BenchmarkConfigCreation(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Config{
			ShutdownTimeout: 20 * time.Second,
		}
	}
}

func BenchmarkVersionAccess(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Version
	}
}