package credentials

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var (
	valueA = Value{APIKey: "EXOaaaaaaaaaaaaaaaaaaaaaaaa", APISecret: "secret-a"}
	valueB = Value{APIKey: "EXObbbbbbbbbbbbbbbbbbbbbbbb", APISecret: "secret-b"}
)

// writeSecretVolume writes key and secret in dir the way the kubelet updates a Secret volume:
// the files go in a new timestamped directory, the "..data" symlink is atomically swapped to it,
// and the visible files are symlinks through "..data".
func writeSecretVolume(t *testing.T, dir, key, secret string) {
	t.Helper()

	version, err := os.MkdirTemp(dir, "..version-")
	if err != nil {
		t.Fatal(err)
	}
	for name, content := range map[string]string{SecretDirAPIKeyFile: key, SecretDirAPISecretFile: secret} {
		if err := os.WriteFile(filepath.Join(version, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	tmp := filepath.Join(dir, "..data_tmp")
	if err := os.Symlink(filepath.Base(version), tmp); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(tmp, filepath.Join(dir, "..data")); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{SecretDirAPIKeyFile, SecretDirAPISecretFile} {
		link := filepath.Join(dir, name)
		if _, err := os.Lstat(link); err == nil {
			continue
		}
		if err := os.Symlink(filepath.Join("..data", name), link); err != nil {
			t.Fatal(err)
		}
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestSecretDirSource(t *testing.T) {
	t.Run("secret volume update", func(t *testing.T) {
		dir := t.TempDir()
		source := SecretDirSource(dir)

		writeSecretVolume(t, dir, valueA.APIKey, valueA.APISecret+"\n")
		if v, err := source(); err != nil || v != valueA {
			t.Fatalf("got %v, %v, want %v", v, err, valueA)
		}

		writeSecretVolume(t, dir, valueB.APIKey, valueB.APISecret)
		if v, err := source(); err != nil || v != valueB {
			t.Fatalf("after update, got %v, %v, want %v", v, err, valueB)
		}
	})

	t.Run("plain files rewritten in place", func(t *testing.T) {
		dir := t.TempDir()
		source := SecretDirSource(dir)

		writeFile(t, filepath.Join(dir, SecretDirAPIKeyFile), valueA.APIKey)
		writeFile(t, filepath.Join(dir, SecretDirAPISecretFile), valueA.APISecret)
		if v, err := source(); err != nil || v != valueA {
			t.Fatalf("got %v, %v, want %v", v, err, valueA)
		}

		writeFile(t, filepath.Join(dir, SecretDirAPIKeyFile), valueB.APIKey)
		writeFile(t, filepath.Join(dir, SecretDirAPISecretFile), valueB.APISecret)
		if v, err := source(); err != nil || v != valueB {
			t.Fatalf("after rewrite, got %v, %v, want %v", v, err, valueB)
		}
	})

	t.Run("empty secret", func(t *testing.T) {
		dir := t.TempDir()
		writeSecretVolume(t, dir, valueA.APIKey, " \n")

		if _, err := SecretDirSource(dir)(); !errors.Is(err, ErrMissingIncomplete) {
			t.Fatalf("got error %v, want %v", err, ErrMissingIncomplete)
		}
	})

	t.Run("missing file", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, filepath.Join(dir, SecretDirAPIKeyFile), valueA.APIKey)

		if _, err := SecretDirSource(dir)(); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("got error %v, want %v", err, os.ErrNotExist)
		}
	})
}

func TestJSONFileSource(t *testing.T) {
	path := filepath.Join(t.TempDir(), "credentials.json")
	source := JSONFileSource(path)

	if _, err := source(); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("missing file: got error %v, want %v", err, os.ErrNotExist)
	}

	writeFile(t, path, fmt.Sprintf(`{"name": "ccm", "api_key": %q, "api_secret": %q}`, valueA.APIKey, valueA.APISecret))
	if v, err := source(); err != nil || v != valueA {
		t.Fatalf("got %v, %v, want %v", v, err, valueA)
	}

	writeFile(t, path, fmt.Sprintf(`{"api_key": %q}`, valueA.APIKey))
	if _, err := source(); !errors.Is(err, ErrMissingIncomplete) {
		t.Fatalf("incomplete file: got error %v, want %v", err, ErrMissingIncomplete)
	}

	writeFile(t, path, `{"api_key": `)
	if _, err := source(); err == nil {
		t.Fatal("invalid JSON: got no error")
	}
}

// testSource is a Source returning a value set by the test, and counting its calls.
type testSource struct {
	mu    sync.Mutex
	value Value
	err   error
	calls atomic.Int32
}

func (s *testSource) set(v Value, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.value, s.err = v, err
}

func (s *testSource) read() (Value, error) {
	s.calls.Add(1)

	s.mu.Lock()
	defer s.mu.Unlock()

	return s.value, s.err
}

func newTestProvider(t *testing.T, source *testSource, opts ...ReloadableOpt) *ReloadableProvider {
	t.Helper()

	opts = append([]ReloadableOpt{ReloadableOptWithInterval(0)}, opts...)
	p, err := NewReloadableProvider(t.Context(), source.read, opts...)
	if err != nil {
		t.Fatal(err)
	}

	return p
}

// retrieve checks the expired state of p, then that Retrieve returns want.
func retrieve(t *testing.T, p *ReloadableProvider, wantExpired bool, want Value) {
	t.Helper()

	if got := p.IsExpired(); got != wantExpired {
		t.Fatalf("IsExpired() = %v, want %v", got, wantExpired)
	}
	v, err := p.Retrieve()
	if err != nil || v != want {
		t.Fatalf("Retrieve() = %v, %v, want %v", v, err, want)
	}
	if p.IsExpired() {
		t.Fatal("IsExpired() = true after Retrieve")
	}
}

func TestReloadableProviderReload(t *testing.T) {
	source := &testSource{value: valueA}
	p := newTestProvider(t, source)
	retrieve(t, p, true, valueA)

	// Unchanged credentials.
	if err := p.Reload(t.Context()); err != nil {
		t.Fatal(err)
	}
	retrieve(t, p, false, valueA)

	// New credentials.
	source.set(valueB, nil)
	p.lastReloadCall = time.Time{}
	if err := p.Reload(t.Context()); err != nil {
		t.Fatal(err)
	}
	retrieve(t, p, true, valueB)

	// The source can't be read: the last credentials are kept.
	source.set(Value{}, ErrMissingIncomplete)
	p.lastReloadCall = time.Time{}
	if err := p.Reload(t.Context()); !errors.Is(err, ErrMissingIncomplete) {
		t.Fatalf("got error %v, want %v", err, ErrMissingIncomplete)
	}
	retrieve(t, p, false, valueB)
}

func TestReloadableProviderReloadRateLimit(t *testing.T) {
	source := &testSource{value: valueA}
	p := newTestProvider(t, source)

	for range 3 {
		if err := p.Reload(t.Context()); err != nil {
			t.Fatal(err)
		}
	}
	// The first read is the one of NewReloadableProvider.
	if got := source.calls.Load(); got != 2 {
		t.Fatalf("source read %d times, want 2", got)
	}
}

func TestReloadableProviderValidator(t *testing.T) {
	errRejected := errors.New("invalid key")
	validator := func(_ context.Context, v Value) error {
		if v == valueB {
			return errRejected
		}
		return nil
	}

	source := &testSource{value: valueB}
	if _, err := NewReloadableProvider(t.Context(), source.read, ReloadableOptWithValidator(validator)); !errors.Is(err, errRejected) {
		t.Fatalf("initial credentials rejected: got error %v, want %v", err, errRejected)
	}

	source.set(valueA, nil)
	p := newTestProvider(t, source, ReloadableOptWithValidator(validator))
	retrieve(t, p, true, valueA)

	source.set(valueB, nil)
	if err := p.Reload(t.Context()); !errors.Is(err, errRejected) {
		t.Fatalf("got error %v, want %v", err, errRejected)
	}
	retrieve(t, p, false, valueA)
}

func TestNewReloadableProviderError(t *testing.T) {
	source := &testSource{err: ErrMissingIncomplete}
	if _, err := NewReloadableProvider(t.Context(), source.read); !errors.Is(err, ErrMissingIncomplete) {
		t.Fatalf("got error %v, want %v", err, ErrMissingIncomplete)
	}
}

func TestReloadableProviderPeriodicReload(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	source := &testSource{value: valueA}
	p, err := NewReloadableProvider(ctx, source.read, ReloadableOptWithInterval(5*time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}
	retrieve(t, p, true, valueA)

	source.set(valueB, nil)
	waitFor(t, p.IsExpired)
	retrieve(t, p, true, valueB)

	cancel()
	time.Sleep(50 * time.Millisecond)
	calls := source.calls.Load()
	time.Sleep(50 * time.Millisecond)
	if got := source.calls.Load(); got != calls {
		t.Fatalf("source read %d times after the context was cancelled", got-calls)
	}
}

func waitFor(t *testing.T, cond func() bool) {
	t.Helper()

	deadline := time.Now().Add(5 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatal("timed out")
		}
		time.Sleep(time.Millisecond)
	}
}

func TestReloadableCredentials(t *testing.T) {
	dir := t.TempDir()
	writeSecretVolume(t, dir, valueA.APIKey, valueA.APISecret)

	creds, err := NewReloadableCredentials(t.Context(), SecretDirSource(dir), ReloadableOptWithInterval(0))
	if err != nil {
		t.Fatal(err)
	}
	if v, err := creds.Get(); err != nil || v != valueA {
		t.Fatalf("Get() = %v, %v, want %v", v, err, valueA)
	}

	writeSecretVolume(t, dir, valueB.APIKey, valueB.APISecret)
	if v, err := creds.Get(); err != nil || v != valueA {
		t.Fatalf("before Reload, Get() = %v, %v, want %v", v, err, valueA)
	}
	if err := creds.Reload(t.Context()); err != nil {
		t.Fatal(err)
	}
	if v, err := creds.Get(); err != nil || v != valueB {
		t.Fatalf("after Reload, Get() = %v, %v, want %v", v, err, valueB)
	}
}

func TestCredentialsReloadChain(t *testing.T) {
	source := &testSource{value: valueA}
	creds := NewChainCredentials([]Provider{
		&StaticProvider{},
		newTestProvider(t, source),
	})
	if v, err := creds.Get(); err != nil || v != valueA {
		t.Fatalf("Get() = %v, %v, want %v", v, err, valueA)
	}

	source.set(valueB, nil)
	if err := creds.Reload(t.Context()); err != nil {
		t.Fatal(err)
	}
	if v, err := creds.Get(); err != nil || v != valueB {
		t.Fatalf("after Reload, Get() = %v, %v, want %v", v, err, valueB)
	}
}

func TestCredentialsReloadStatic(t *testing.T) {
	creds := NewStaticCredentials(valueA.APIKey, valueA.APISecret)
	if err := creds.Reload(t.Context()); err != nil {
		t.Fatal(err)
	}
	if v, err := creds.Get(); err != nil || v != valueA {
		t.Fatalf("Get() = %v, %v, want %v", v, err, valueA)
	}
}

func TestReloadableCredentialsConcurrentUse(t *testing.T) {
	source := &testSource{value: valueA}
	creds, err := NewReloadableCredentials(t.Context(), source.read, ReloadableOptWithInterval(time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	for i := range 8 {
		wg.Go(func() {
			for range 200 {
				if i == 0 {
					source.set(valueB, nil)
				}
				v, err := creds.Get()
				if err != nil || (v != valueA && v != valueB) {
					t.Errorf("Get() = %v, %v", v, err)
					return
				}
				_ = creds.Reload(t.Context())
			}
		})
	}
	wg.Wait()
}
