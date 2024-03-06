package template

import _ "embed"

//go:embed postgres/files/create.tmpl
var PostgresCreate []byte

//go:embed postgres/files/read.tmpl
var PosrgresRead []byte

//go:embed postgres/files/update.tmpl
var PostgresUpdate []byte

//go:embed postgres/files/delete.tmpl
var PostgresDelete []byte
