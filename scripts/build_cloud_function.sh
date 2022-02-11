#!/usr/bin/env bash

mkdir -p $1
cp cmd/functions/$1/function.go $1/function.go
sed "s/replace\sgithub\.com\/prixfixeco\/api_server\s=>\s\.\.\/\.\.\/\.\.\//replace github\.com\/prixfixeco\/api_server => \.\.\//" cmd/functions/$1/go.mod > $1/go.mod
(cd $1 && go mod tidy && go mod vendor)
(cd $1 && zip -r -D $1_cloud_function.zip * && mv $1_cloud_function.zip ../)
rm -rf $1
