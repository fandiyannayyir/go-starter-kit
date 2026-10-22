// Package migrations embeds the goose SQL migrations into the binary, so a
// deployed binary carries its own schema and `migrate` works from any directory
// without shipping the migrations folder alongside it.
package migrations

import "embed"

//go:embed *.sql
var Embed embed.FS
