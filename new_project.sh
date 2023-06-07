#!/bin/bash
# Author: Ropon
# Email: luopeng@codoon.com

git clone "https://github.com/ropon/work_api.git" $1

pushd $(pwd)/$1
rm -rf .git
rm -f new_project.sh
find . -type f|grep -vE "air$"|xargs -L 1 sed -i "" "s/work_api/$1/g"
mv conf/work_api.json conf/$1.json && cp conf/$1.json conf/${1}_dev.json
find . -type f|grep -vE "air$"|xargs -L 1 sed -i "" "s/2345/$2/g"
mkdir log
git config --global init.defaultBranch master
git init
go mod tidy
git remote add origin https://github.com/ropon/$1.git
echo "init project success"
popd
