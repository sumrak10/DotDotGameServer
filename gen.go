//go:generate protoc -I=proto --go_out=pkg/pb/go --go_opt=paths=source_relative proto/game/input.proto proto/game/output.proto proto/game/world/world.proto proto/game/world/nodes/nodes.proto

//js:generate protoc -I=proto --js_out=import_style=commonjs,binary:pkg/pb/js --grpc-web_out=import_style=typescript,mode=grpcwebtext:pkg/pb/js proto/game/input.proto proto/game/output.proto proto/game/world/world.proto proto/game/world/nodes/nodes.proto
//js:generate protoc -I=proto --js_out=import_style=closure,binary:pkg/pb/js proto/game/input.proto proto/game/output.proto proto/game/world/world.proto proto/game/world/nodes/nodes.proto
//js:generate protoc -I=proto --js_out=library=pkg/pb/js/game,binary:. proto/game/input.proto proto/game/output.proto proto/game/world/world.proto proto/game/world/nodes/nodes.proto

package main
