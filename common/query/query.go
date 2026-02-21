package query

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rs/zerolog"
)

type Config struct {
	Path string `yaml:"path"`
}

type QueryComponent struct {
	log     zerolog.Logger
	cfg     Config
	queries map[string]string
}

// NewQueryComponent creates a new component but does not load queries yet.
func NewQueryComponent(log zerolog.Logger, cfg Config) *QueryComponent {
	return &QueryComponent{
		log:     log,
		cfg:     cfg,
		queries: make(map[string]string),
	}
}

// Start loads all .sql files from the configured directory, parses them,
// and then blocks until the context is cancelled.
// It returns an error if no files are found or if any file cannot be read.
func (qc *QueryComponent) Start(ctx context.Context) error {
	// Find all .sql files in the directory
	files, err := filepath.Glob(filepath.Join(qc.cfg.Path, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to glob SQL files: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no SQL files found in path: %s", qc.cfg.Path)
	}

	// Load each file
	for _, file := range files {
		if err := qc.loadFile(file); err != nil {
			return fmt.Errorf("failed to load file %s: %w", file, err)
		}
	}

	qc.log.Debug().Msgf("Queries loaded successfully, total queries: %d", len(qc.queries))

	// Block until shutdown signal
	<-ctx.Done()

	qc.log.Debug().Msg("Query component context cancelled – stopping")
	return nil
}

// loadFile reads a single SQL file and adds its named queries to the map.
func (qc *QueryComponent) loadFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	content := string(data)
	// Split by "-- name:" markers (your existing format)
	sections := strings.SplitSeq(content, "-- name:")

	for section := range sections {
		if strings.TrimSpace(section) == "" {
			continue
		}

		lines := strings.Split(section, "\n")
		if len(lines) < 2 {
			continue
		}

		name := strings.TrimSpace(lines[0])
		query := strings.Join(lines[1:], "\n")
		query = strings.TrimSpace(query)
		query = strings.TrimSuffix(query, ";")

		qc.queries[name] = query
	}

	qc.log.Debug().Str("file", filepath.Base(filePath)).Msg("Loaded queries from file")
	return nil
}

// Stop performs any necessary cleanup. For this component, nothing is required.
func (qc *QueryComponent) Stop(ctx context.Context) error {
	qc.log.Debug().Msg("Query component stopped")
	return nil
}

// Get returns a query by its name and a boolean indicating if it exists.
func (qc *QueryComponent) Get(name string) (string, bool) {
	query, ok := qc.queries[name]
	return query, ok
}

// ExecuteTemplate processes a named query with the given data,
// converting named placeholders ($key) to positional parameters ($1, $2, …).
func (qc *QueryComponent) ExecuteTemplate(name string, data any) (string, []any, error) {
	queryTemplate, ok := qc.Get(name)
	if !ok {
		return "", nil, fmt.Errorf("query %s not found", name)
	}

	tmpl, err := template.New(name).Parse(queryTemplate)
	if err != nil {
		return "", nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", nil, err
	}

	query := buf.String()
	return convertNamedToPositional(query, data)
}

// convertNamedToPositional replaces $key placeholders with $1, $2, … and builds the args slice.
// It expects data to be a map[string]any.
func convertNamedToPositional(query string, data any) (string, []any, error) {
	paramMap, ok := data.(map[string]any)
	if !ok {
		return "", nil, fmt.Errorf("data must be map[string]any for named parameter conversion")
	}

	args := make([]any, 0)
	paramIndex := 1
	result := query

	// Replace each $key with $N and collect values in order
	for key, value := range paramMap {
		placeholder := "$" + key
		if strings.Contains(result, placeholder) {
			positional := fmt.Sprintf("$%d", paramIndex)
			result = strings.ReplaceAll(result, placeholder, positional)
			args = append(args, value)
			paramIndex++
		}
	}

	return result, args, nil
}
