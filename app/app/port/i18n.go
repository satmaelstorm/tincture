package port

// Localizer defines the interface for localization functionality
type Localizer interface {
	// Localize returns the localized string for the given key
	Localize(key string) string

	// LocalizeWithTemplate returns the localized string with template data
	LocalizeWithTemplate(key string, data interface{}) string
}
