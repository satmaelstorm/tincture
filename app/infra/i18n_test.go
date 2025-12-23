package infra

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestI18nAdapter_Localize(t *testing.T) {
	// Create a new i18n adapter
	i18n, err := NewI18nAdapter()
	assert.NoError(t, err)

	// Test localization of a known key
	result := i18n.Localize("addTinctureButton")
	assert.Equal(t, "Добавить", result)

	// Test localization of an unknown key
	result = i18n.Localize("unknownKey")
	assert.Equal(t, "unknownKey", result)
}

func TestI18nAdapter_LocalizeWithTemplate(t *testing.T) {
	// Create a new i18n adapter
	i18n, err := NewI18nAdapter()
	assert.NoError(t, err)

	// Test localization with template data
	result := i18n.LocalizeWithTemplate("tinctureConfirmDelete", map[string]interface{}{
		"Name": "Test Tincture",
	})
	assert.Equal(t, "Вы хотите пометить Test Tincture как выпитое.", result)
}
