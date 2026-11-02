// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package config_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/tidwall/gjson"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/kratos/embedx"
)

// Every courier.templates.* key that the code reads must be declared in the
// config schema, otherwise Kratos refuses to start when it is configured.
func TestCourierTemplateKeysAreInSchema(t *testing.T) {
	src, err := os.ReadFile("config.go")
	require.NoError(t, err)

	keys := regexp.MustCompile(`"(courier\.templates\.[a-z_.]+)"`).FindAllStringSubmatch(string(src), -1)
	require.NotEmpty(t, keys)

	for _, m := range keys {
		key := m[1]
		t.Run(key, func(t *testing.T) {
			// courier.templates.<name>[...]: the template name must be declared.
			name := strings.Split(key, ".")[2]
			node := "properties.courier.properties.templates.properties." + name
			assert.True(t, gjson.GetBytes(embedx.ConfigSchema, node).Exists(), "%s missing from config.schema.json", key)
		})
	}
}
