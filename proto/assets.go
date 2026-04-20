package proto

import "embed"

//go:embed all:game
//go:embed protobuf.min.js
var StaticFiles embed.FS
