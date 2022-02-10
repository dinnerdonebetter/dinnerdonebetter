#!/usr/bin/env bash

mkdir -p data_changes_cloud_func
cp cmd/functions/data_changes/function.go data_changes_cloud_func/function.go
sed "s/replace\sgithub\.com\/prixfixeco\/api_server\s=>\s\.\.\/\.\.\/\.\.\//replace github\.com\/prixfixeco\/api_server => \.\.\//" cmd/functions/data_changes/go.mod > data_changes_cloud_func/go.mod
(cd data_changes_cloud_func && go mod tidy && go mod vendor)
(cd data_changes_cloud_func && zip -r -D data_changes_cloud_function.zip * && mv data_changes_cloud_function.zip ../)
rm -rf data_changes_cloud_func
