package public

import "embed"

//go:embed css/* js/*
var Assets embed.FS
