#!/bin/sh
set -euo pipefail

mkdir -p build/man
for APP in holo holo-build; do
  find "../${APP}/doc" -name '*.pod' -exec cp -t build/man {} +
done
