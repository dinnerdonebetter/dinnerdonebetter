#!/usr/bin/env bash

mkdir -p data_changes_cloud_func/vendor
rsync -r vendor data_changes_cloud_func
cp cmd/functions/data_changes/function.go data_changes_cloud_func/function.go
cp go.mod data_changes_cloud_func/go.mod
cd data_changes_cloud_func/ || exit
zip data_changes_cloud_function.zip function.go vendor/**/* vendor/modules.txt go.mod
mv data_changes_cloud_function.zip ../
cd ../
rm -rf data_changes_cloud_function