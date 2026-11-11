package credentials

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	// DefaultReloadInterval is how often a ReloadableProvider re-reads its Source by default.
	DefaultReloadInterval = 30 * time.Second

	// reloadMinInterval rate-limits the on-demand reloads (Reload), so that a burst of requests
	// rejected at the same time reads the Source, and calls the validator, only once.
	reloadMinInterval = 5 * time.Second

	// SecretDirAPIKeyFile and SecretDirAPISecretFile are the files SecretDirSource reads:
	// the keys of a Kubernetes Secret mounted as a volume.
	SecretDirAPIKeyFile    = "EXOSCALE_API_KEY"
	SecretDirAPISecretFile = "EXOSCALE_API_SECRET"
)

// A Reloader is a Provider whose credentials can change while the program runs.
// The v3 client calls Reload when the API rejects the key it signed a request with,
// then retries the request once if the credentials changed.
type Reloader interface {
	// Reload re-reads the credentials. Changes are reported by the Provider's IsExpired.
	Reload(ctx context.Context) error
}

// A Source reads the credentials from where they are stored, every time it is called.
type Source func() (Value, error)

// SecretDirSource returns a Source reading the API key and secret from the EXOSCALE_API_KEY and
// EXOSCALE_API_SECRET files of dir, typically a Kubernetes Secret mounted as a volume.
//
// The kubelet updates such a volume by swapping its "..data" symlink: the Source resolves it once,
// so the key and the secret always come from the same version of the Secret.
func SecretDirSource(dir string) Source {
	return func() (Value, error) {
		base := dir
		if resolved, err := filepath.EvalSymlinks(filepath.Join(dir, "..data")); err == nil {
			base = resolved
		}

		read := func(name string) (string, error) {
			b, err := os.ReadFile(filepath.Join(base, name))
			if err != nil {
				return "", err
			}

			return strings.TrimSpace(string(b)), nil
		}

		apiKey, err := read(SecretDirAPIKeyFile)
		if err != nil {
			return Value{}, fmt.Errorf("secret dir source: %w", err)
		}
		apiSecret, err := read(SecretDirAPISecretFile)
		if err != nil {
			return Value{}, fmt.Errorf("secret dir source: %w", err)
		}

		v := Value{APIKey: apiKey, APISecret: apiSecret}
		if !v.IsSet() {
			return Value{}, fmt.Errorf("secret dir source: %s: %w", dir, ErrMissingIncomplete)
		}

		return v, nil
	}
}

// JSONFileSource returns a Source reading the API key and secret from a JSON file of the form
// {"api_key": "...", "api_secret": "..."}, the format of the Exoscale cloud controller manager.
// Other fields are ignored.
func JSONFileSource(path string) Source {
	return func() (Value, error) {
		b, err := os.ReadFile(path)
		if err != nil {
			return Value{}, fmt.Errorf("json file source: %w", err)
		}

		var file struct {
			APIKey    string `json:"api_key"`
			APISecret string `json:"api_secret"`
		}
		if err := json.Unmarshal(b, &file); err != nil {
			return Value{}, fmt.Errorf("json file source: %s: %w", path, err)
		}

		v := Value{APIKey: strings.TrimSpace(file.APIKey), APISecret: strings.TrimSpace(file.APISecret)}
		if !v.IsSet() {
			return Value{}, fmt.Errorf("json file source: %s: %w", path, ErrMissingIncomplete)
		}

		return v, nil
	}
}

// ReloadableOpt represents a ReloadableProvider option.
type ReloadableOpt func(*ReloadableProvider)

// ReloadableOptWithInterval returns a ReloadableOpt setting how often the Source is re-read
// (DefaultReloadInterval by default). Zero or less disables the periodic re-read: the Source
// is then only re-read when the API rejects the key.
func ReloadableOptWithInterval(d time.Duration) ReloadableOpt {
	return func(p *ReloadableProvider) {
		p.interval = d
	}
}

// ReloadableOptWithValidator returns a ReloadableOpt checking new credentials before they are
// used: credentials it rejects are ignored, the previous ones are kept and the new ones are
// checked again on the next re-read. The validator should only reject credentials the API
// refuses, not fail on a transient error.
//
// The v3 package provides one: v3.CredentialsValidator.
func ReloadableOptWithValidator(validate func(ctx context.Context, v Value) error) ReloadableOpt {
	return func(p *ReloadableProvider) {
		p.validate = validate
	}
}

// ReloadableOptWithLogHandler returns a ReloadableOpt logging the reloads, and the errors when
// the Source can't be read or the credentials are rejected. Nothing is logged by default.
func ReloadableOptWithLogHandler(h slog.Handler) ReloadableOpt {
	return func(p *ReloadableProvider) {
		p.logger = slog.New(h)
	}
}

// A ReloadableProvider is a Provider whose credentials are re-read from a Source while the
// program runs, so that a rotated API key is used without restarting it, and without creating
// a new client: the v3 client gets its credentials for every request.
//
// The Source is re-read periodically, and right away when the API rejects the key (see Reloader).
// New complete credentials are used as soon as they are read, unless a validator rejects them.
// If the Source can't be read or is incomplete, the last credentials are kept and the error is
// logged on every re-read.
//
// Retrieve doesn't read the Source: it returns the last credentials read.
type ReloadableProvider struct {
	source   Source
	interval time.Duration
	validate func(context.Context, Value) error
	logger   *slog.Logger

	// reloadMu serializes the reads of the Source.
	reloadMu sync.Mutex
	// lastReloadCall is the time of the last read made by Reload, to rate-limit them.
	lastReloadCall time.Time

	mu      sync.RWMutex
	value   Value
	expired bool
}

// NewReloadableProvider returns a ReloadableProvider reading its credentials from source.
// The first read happens here and its error, or the validator's, is returned.
// The periodic re-read runs in the background until ctx is done.
func NewReloadableProvider(ctx context.Context, source Source, opts ...ReloadableOpt) (*ReloadableProvider, error) {
	p := &ReloadableProvider{
		source:   source,
		interval: DefaultReloadInterval,
		logger:   slog.New(slog.DiscardHandler),
	}
	for _, opt := range opts {
		opt(p)
	}

	v, err := p.source()
	if err != nil {
		return nil, fmt.Errorf("reloadable provider: %w", err)
	}
	if p.validate != nil {
		if err := p.validate(ctx, v); err != nil {
			return nil, fmt.Errorf("reloadable provider: credentials rejected: %w", err)
		}
	}
	p.value = v
	p.expired = true

	if p.interval > 0 {
		go p.run(ctx)
	}

	return p, nil
}

// NewReloadableCredentials returns Credentials backed by a ReloadableProvider: see
// NewReloadableProvider.
func NewReloadableCredentials(ctx context.Context, source Source, opts ...ReloadableOpt) (*Credentials, error) {
	p, err := NewReloadableProvider(ctx, source, opts...)
	if err != nil {
		return nil, err
	}

	return NewCredentials(p), nil
}

// Retrieve returns the last credentials read from the Source.
func (p *ReloadableProvider) Retrieve() (Value, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.expired = false

	return p.value, nil
}

// IsExpired returns true when new credentials were read since the last Retrieve.
func (p *ReloadableProvider) IsExpired() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.expired
}

// Reload re-reads the Source now, unless Reload already did less than a few seconds ago.
// It returns an error if the Source can't be read or the new credentials are rejected;
// the previous credentials are kept in that case.
func (p *ReloadableProvider) Reload(ctx context.Context) error {
	return p.reload(ctx, false)
}

func (p *ReloadableProvider) run(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = p.reload(ctx, true)
		}
	}
}

func (p *ReloadableProvider) reload(ctx context.Context, force bool) error {
	p.reloadMu.Lock()
	defer p.reloadMu.Unlock()

	if !force {
		if time.Since(p.lastReloadCall) < reloadMinInterval {
			return nil
		}
		p.lastReloadCall = time.Now()
	}

	v, err := p.source()
	if err != nil {
		p.logger.ErrorContext(ctx, "keeping current Exoscale API credentials, unable to read new ones", "error", err)
		return fmt.Errorf("reloadable provider: %w", err)
	}

	p.mu.RLock()
	current := p.value
	p.mu.RUnlock()
	if v == current {
		return nil
	}

	if p.validate != nil {
		if err := p.validate(ctx, v); err != nil {
			p.logger.ErrorContext(ctx, "keeping current Exoscale API credentials, new ones rejected",
				"api_key", v.APIKey, "error", err)
			return fmt.Errorf("reloadable provider: credentials rejected: %w", err)
		}
	}

	p.mu.Lock()
	p.value = v
	p.expired = true
	p.mu.Unlock()

	p.logger.InfoContext(ctx, "Exoscale API credentials reloaded", "api_key", v.APIKey)

	return nil
}
