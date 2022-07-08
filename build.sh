#! /bin/bash

cd qa-test-web
pnpm build

rm -rf ../templates/assets

cp -rf dist/* ../templates/