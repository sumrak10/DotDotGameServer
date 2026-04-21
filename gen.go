//go:generate protoc -I=proto --go_out=pkg/pb/go --go_opt=paths=source_relative proto/game/input.proto proto/game/output.proto proto/game/world/world.proto proto/game/world/nodes/nodes.proto
package main
