#!/usr/bin/env bash

directory=$1

mkdir -p "$directory"
cp "cmd/functions/$directory/function.go" "$directory/function.go"
sed "s/replace\sgithub\.com\/prixfixeco\/api_server\s=>\s\.\.\/\.\.\/\.\.\//replace github\.com\/prixfixeco\/api_server => \.\.\//" cmd/functions/${directory}/go.mod > ${directory}/go.mod
(cd "$directory" && go mod tidy && go mod vendor && zip -r -D ${directory}_cloud_function.zip * && mv ${directory}_cloud_function.zip ../)
rm -rf "$directory"
