#!/bin/bash
# Author: Ropon
# Email: luopeng@codoon.com

git clone "https://github.com/ropon/work_api.git" $1

pushd $(pwd)/$1
rm -rf .git
rm -f new_project.sh
find . -type f|grep -vE "air$"|xargs -L 1 sed -i "" "s/work_api/$1/g"
mv conf/work_api.conf conf/$1.conf
mkdir log
git init
git remote add origin https://github.com/ropon/$1.git
echo "init project success"
popd
