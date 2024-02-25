package ui

import (
	"embed"
)

// Specifying which folders to embed into executable using "go:embed <paths>" format
// "html" embeds ui/html folder
// "static" embeds ui/static folder

//go:embed "html" "static"
var Files embed.FS
