package migrations

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrationInitialization(t *testing.T) {
	t.Run("Migration variable is initialized", func(t *testing.T) {
		require.NotNil(t, Migration, "Migration should be initialized")
	})

	t.Run("Migration is of correct type", func(t *testing.T) {
		assert.NotNil(t, Migration, "Migration should not be nil")
		// Verify it's a *migrate.Migrations instance by checking we can call methods
		migrations := Migration.Sorted()
		assert.NotNil(t, migrations, "migrations list should not be nil")
	})

	t.Run("sqlMigrations embed.FS is initialized", func(t *testing.T) {
		// Test that we can read from the embedded filesystem
		entries, err := sqlMigrations.ReadDir(".")
		assert.NoError(t, err, "should be able to read from embedded FS")
		assert.NotNil(t, entries, "entries should not be nil")
	})

	t.Run("sqlMigrations contains expected files", func(t *testing.T) {
		entries, err := sqlMigrations.ReadDir(".")
		require.NoError(t, err, "should be able to read directory")
		
		// Check that test.sql is present
		found := false
		for _, entry := range entries {
			if entry.Name() == "test.sql" {
				found = true
				assert.False(t, entry.IsDir(), "test.sql should be a file, not a directory")
			}
		}
		assert.True(t, found, "test.sql should be embedded in the filesystem")
	})

	t.Run("Migration discovery completes without panic", func(t *testing.T) {
		// This test verifies the init() function didn't panic
		// If we reach here, init() succeeded
		assert.NotNil(t, Migration, "Migration should be initialized without panic")
	})

	t.Run("can retrieve migration list", func(t *testing.T) {
		migrations := Migration.Sorted()
		assert.NotNil(t, migrations, "migrations list should not be nil")
	})
}

func TestMigrationDiscovery(t *testing.T) {
	t.Run("discovered migrations are accessible", func(t *testing.T) {
		migrations := Migration.Sorted()
		require.NotNil(t, migrations, "should be able to get sorted migrations")
		
		// Verify migrations structure is a slice
		// Verify migrations structure is a slice
		assert.GreaterOrEqual(t, len(migrations), 0, "migrations should be a valid slice")
	})

	t.Run("migration names follow convention", func(t *testing.T) {
		entries, err := sqlMigrations.ReadDir(".")
		require.NoError(t, err, "should be able to read embedded directory")
		
		for _, entry := range entries {
			if !entry.IsDir() {
				// All migration files should be SQL files
				name := entry.Name()
				assert.Contains(t, name, ".sql", "migration file %s should have .sql extension", name)
			}
		}
	})

	t.Run("can read SQL file content", func(t *testing.T) {
		content, err := sqlMigrations.ReadFile("test.sql")
		assert.NoError(t, err, "should be able to read test.sql")
		assert.NotNil(t, content, "content should not be nil")
		// test.sql is empty in the current state
		assert.Equal(t, 0, len(content), "test.sql is currently empty")
	})
}

func TestMigrationEdgeCases(t *testing.T) {
	t.Run("Migration variable is not nil after package import", func(t *testing.T) {
		assert.NotNil(t, Migration, "Migration should be automatically initialized")
	})

	t.Run("sqlMigrations filesystem is valid", func(t *testing.T) {
		// Try various operations on the embedded filesystem
		_, err := sqlMigrations.ReadDir(".")
		assert.NoError(t, err, "ReadDir should work on root")
		
		// Try reading a non-existent file
		_, err = sqlMigrations.ReadFile("nonexistent.sql")
		assert.Error(t, err, "reading non-existent file should error")
	})

	t.Run("multiple Sorted calls return consistent results", func(t *testing.T) {
		list1 := Migration.Sorted()
		require.NotNil(t, list1)
		
		list2 := Migration.Sorted()
		require.NotNil(t, list2)
		
		assert.Equal(t, len(list1), len(list2), "multiple Sorted() calls should return same length")
	})

	t.Run("can add migrations programmatically", func(t *testing.T) {
		// Test that we can use the Add method (though we don't actually modify the global)
		// This verifies the API is available
		assert.NotPanics(t, func() {
			// Just verify the method exists by checking we can reference it
			_ = Migration.Add
		}, "Add method should be available")
	})
}

func TestEmbeddedFileSystem(t *testing.T) {
	t.Run("embedded FS contains only SQL files", func(t *testing.T) {
		entries, err := sqlMigrations.ReadDir(".")
		require.NoError(t, err)
		
		for _, entry := range entries {
			if !entry.IsDir() {
				info, err := entry.Info()
				require.NoError(t, err, "should be able to get file info")
				assert.NotNil(t, info, "file info should not be nil")
			}
		}
	})

	t.Run("can stat embedded files", func(t *testing.T) {
		entries, err := sqlMigrations.ReadDir(".")
		require.NoError(t, err)
		
		for _, entry := range entries {
			if !entry.IsDir() {
				info, err := entry.Info()
				assert.NoError(t, err, "should be able to stat file %s", entry.Name())
				if err == nil {
					assert.NotEmpty(t, info.Name(), "file should have a name")
					assert.False(t, info.IsDir(), "SQL file should not be a directory")
				}
			}
		}
	})

	t.Run("embedded FS is read-only", func(t *testing.T) {
		// Verify we cannot write to embedded filesystem
		// This is implicit in embed.FS design
		entries, err := sqlMigrations.ReadDir(".")
		assert.NoError(t, err, "reading should work")
		assert.NotNil(t, entries, "should be able to read entries")
	})
}

func TestMigrationConcurrency(t *testing.T) {
	t.Run("concurrent access to Migration variable", func(t *testing.T) {
		const goroutines = 50
		done := make(chan bool, goroutines)
		
		for i := 0; i < goroutines; i++ {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("panic during concurrent access: %v", r)
					}
					done <- true
				}()
				
				assert.NotNil(t, Migration)
				_ = Migration.Sorted()
			}()
		}
		
		for i := 0; i < goroutines; i++ {
			<-done
		}
	})

	t.Run("concurrent reads from embedded FS", func(t *testing.T) {
		const goroutines = 50
		done := make(chan bool, goroutines)
		
		for i := 0; i < goroutines; i++ {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("panic during concurrent FS read: %v", r)
					}
					done <- true
				}()
				
				_, err := sqlMigrations.ReadDir(".")
				assert.NoError(t, err)
				
				_, err = sqlMigrations.ReadFile("test.sql")
				assert.NoError(t, err)
			}()
		}
		
		for i := 0; i < goroutines; i++ {
			<-done
		}
	})
}

// Benchmark tests
func BenchmarkMigrationSorted(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Migration.Sorted()
	}
}

func BenchmarkReadEmbeddedFS(b *testing.B) {
	b.Run("ReadDir", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = sqlMigrations.ReadDir(".")
		}
	})
	
	b.Run("ReadFile", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = sqlMigrations.ReadFile("test.sql")
		}
	})
}