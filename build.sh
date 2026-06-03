#! /bin/bash
set -e

GOPROXY="${GOPROXY:-https://goproxy.cn,direct}"
export GOPROXY

go build

cd qa-test-web
export npm_config_registry="${npm_config_registry:-https://registry.npmmirror.com}"
if [ ! -x node_modules/.bin/vue-tsc ] || [ ! -x node_modules/.bin/vite ]; then
  pnpm install --prod=false
fi
pnpm build

mkdir -p ../templates
rm -rf ../templates/assets

cp -rf dist/. ../templates/
