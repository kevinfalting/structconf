package stronf

// Set will attempt to set the value provided and is exposed for tests only.
func (f *Field) Set(val any) error {
	return f.set(val)
}
