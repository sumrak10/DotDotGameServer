//go:generate protoc -I=proto --go_out=pkg/pb/go --go_opt=paths=source_relative proto/game/input.proto proto/game/output.proto proto/game/world/world.proto proto/game/world/nodes/nodes.proto
package main

//mkdir certs
//openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
//-keyout certs/nginx-selfsigned.key \
//-out certs/nginx-selfsigned.crt \
//-subj "/C=RU/ST=Moscow/L=Moscow/O=Development/CN=144.91.67.239"
