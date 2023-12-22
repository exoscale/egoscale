package credentials

import "fmt"

type FileOpt func(*FileProvider)

// FileOptWithFilename returns a FileOpt overriding the default filename.
func FileOptWithFilename(filename string) FileOpt {
	return func(f *FileProvider) {
		f.filename = filename
	}
}

// FileOptWithAccount returns a FileOpt overriding the default account.
func FileOptWithAccount(account string) FileOpt {
	return func(f *FileProvider) {
		f.account = account
	}
}

type FileProvider struct {
	filename  string
	account   string
	retrieved bool
}

func NewFileCredentials(opts ...FileOpt) *Credentials {
	fp := &FileProvider{}
	for _, opt := range opts {
		opt(fp)
	}
	return NewCredentials(fp)
}

func (f *FileProvider) Retrieve() (Value, error) {
	return Value{}, fmt.Errorf("not implemented")
}

// IsExpired returns if the shared credentials have expired.
func (p *FileProvider) IsExpired() bool {
	return !p.retrieved
}
