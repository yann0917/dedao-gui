#!/bin/sh
set -eu

bundle=${1:?Usage: check-macos-bundle.sh /path/to/application.app}
plist="$bundle/Contents/Info.plist"

plutil -lint "$plist" >/dev/null
executable=$(plutil -extract CFBundleExecutable raw -o - "$plist")
case "$executable" in
  ""|.|..|*/*)
    echo "Invalid CFBundleExecutable: $executable" >&2
    exit 1
    ;;
esac

if [ ! -f "$bundle/Contents/MacOS/$executable" ] || [ ! -x "$bundle/Contents/MacOS/$executable" ]; then
  echo "Bundle executable missing or not executable: $bundle/Contents/MacOS/$executable" >&2
  exit 1
fi

echo "Bundle executable OK: $executable"
