#!/bin/bash
set -euxo pipefail

for T in dependencies runtime demo ; do
  docker build --target="$T" --tag="jfsmig/hegemonie-$T" .
  docker push "jfsmig/hegemonie-$T:latest"
done
