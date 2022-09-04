package ui

import "embed"

//go:generate npm install
//go:generate npm run build

//go:embed build
var FS embed.FS
