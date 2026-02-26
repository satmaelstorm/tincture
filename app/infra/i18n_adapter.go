package infra

import (
	"embed"
	"os"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/satmaelstorm/tincture/app/app/port"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

//go:embed locales/*.yaml
var localeFiles embed.FS

type i18nAdapter struct {
	localizer *i18n.Localizer
}

func NewI18nAdapter() (port.Localizer, error) {
	// Create a new bundle
	bundle := i18n.NewBundle(language.Russian)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	// Load embedded translation files
	entries, err := localeFiles.ReadDir("locales")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		content, err := localeFiles.ReadFile("locales/" + entry.Name())
		if err != nil {
			return nil, err
		}
		_, err = bundle.ParseMessageFileBytes(content, entry.Name())
		if err != nil {
			return nil, err
		}
	}

	// Create a localizer
	lang := os.Getenv("LANGUAGE")
	if lang == "" {
		lang = "ru"
	}

	localizer := i18n.NewLocalizer(bundle, lang)

	return &i18nAdapter{
		localizer: localizer,
	}, nil
}

func (i *i18nAdapter) Localize(key string) string {
	localized, err := i.localizer.Localize(&i18n.LocalizeConfig{
		MessageID: key,
	})

	if err != nil {
		return key
	}

	return localized
}

func (i *i18nAdapter) LocalizeWithTemplate(key string, data interface{}) string {
	localized, err := i.localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: data,
	})

	if err != nil {
		return key
	}

	return localized
}
