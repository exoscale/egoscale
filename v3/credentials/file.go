package credentials

type FileProvider struct {
}

func NewFileCredentials() *Credentials {
	return NewCredentials(&FileProvider{})
}

func (f *FileProvider) Retrieve() (Value, error) {
	return Value{}, nil
}

func (f *FileProvider) IsExpired() bool {
	return false
}
