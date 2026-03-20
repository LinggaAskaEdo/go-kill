package query

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"text/template"

	"github.com/rs/zerolog"
)

type Config struct {
	Path string `yaml:"path"`
}

type QueryComponent struct {
	log          zerolog.Logger
	cfg          Config
	queries      map[string]string
	templates    map[string]*template.Template
	ready        chan struct{}
	mu           sync.RWMutex
	paramPattern *regexp.Regexp
}

// NewQueryComponent creates a new component but does not load queries yet.
func NewQueryComponent(log zerolog.Logger, cfg Config) *QueryComponent {
	return &QueryComponent{
		log:          log,
		cfg:          cfg,
		queries:      make(map[string]string),
		templates:    make(map[string]*template.Template),
		ready:        make(chan struct{}),
		paramPattern: regexp.MustCompile(`\$([a-zA-Z_][a-zA-Z0-9_]*)`),
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

	close(qc.ready) // signal readiness
	qc.log.Debug().Msgf("Queries loaded successfully, total queries: %d", len(qc.queries))
	<-ctx.Done() // Block until shutdown signal
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
	qc.mu.RLock()
	tmpl, exists := qc.templates[name]
	qc.mu.RUnlock()

	if !exists {
		queryTemplate, ok := qc.Get(name)
		if !ok {
			return "", nil, fmt.Errorf("query %s not found", name)
		}

		var err error
		tmpl, err = template.New(name).Parse(queryTemplate)
		if err != nil {
			return "", nil, err
		}

		qc.mu.Lock()
		qc.templates[name] = tmpl
		qc.mu.Unlock()
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", nil, err
	}

	query := buf.String()
	return qc.convertNamedToPositional(query, data)
}

// convertNamedToPositional replaces $key placeholders with $1, $2, … in the order they appear.
// It expects data to be a map[string]any.
func (qc *QueryComponent) convertNamedToPositional(query string, data any) (string, []any, error) {
	paramMap, ok := data.(map[string]any)
	if !ok {
		return "", nil, fmt.Errorf("data must be map[string]any for named parameter conversion")
	}

	matches := qc.paramPattern.FindAllStringSubmatchIndex(query, -1)

	if len(matches) == 0 {
		return query, nil, nil
	}

	args := make([]any, 0, len(matches))
	var result bytes.Buffer
	offset := 0

	for _, match := range matches {
		fullStart := match[0]
		fullEnd := match[1]
		keyStart := match[2]
		keyEnd := match[3]

		result.WriteString(query[offset:fullStart])

		key := query[keyStart:keyEnd]
		value, exists := paramMap[key]
		if !exists {
			return "", nil, fmt.Errorf("parameter $%s not found in data", key)
		}

		result.WriteString(fmt.Sprintf("$%d", len(args)+1))
		args = append(args, value)
		offset = fullEnd
	}

	result.WriteString(query[offset:])

	return result.String(), args, nil
}

// Ready returns a channel that is closed when the connection is established.
func (qc *QueryComponent) Ready() <-chan struct{} {
	return qc.ready
}
