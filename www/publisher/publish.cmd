#!cmd alpine openssh bash hugo make git ca-certificates openssl
#!/bin/bash

REPO_NAME="gliderlabs/cmd"
BASE_URL="https://gliderlabs.github.io/cmd"
BUILD_CMD="make www-build"
BUILD_DIR="build/www/public"
DOCS_DIR="docs"
WWW_DIR="www"
COMMIT_MSG="published via www/publisher/publish"
TARGET_BRANCH="gh-pages"
EXPECT_BRANCH="master"

set -eo pipefail
mkdir -p /root/.ssh
cat > /root/.ssh/id_rsa << EOF
{{ $key }}
EOF
chmod 600 /root/.ssh/id_rsa

main() {
  echo -e "Host github.com\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config
  git config --global user.name "Publisher"
  git config --global user.email "team+robot@gliderlabs.com"
  chown -R $(whoami) /root/.ssh

  cat | tar -xpf -

  local checksum="$(checksum-src)"
  local url="${BASE_URL}/checksum.txt?$(date +%s)"
  local deployed="$(wget -qO- $url)"
  if [[ "$checksum" == "$deployed" ]]; then
    echo "No changes to deploy."
    exit
  fi
  rm -rf .git
  $BUILD_CMD
  cd $BUILD_DIR
  echo "$checksum" > checksum.txt
  git init
  git add .
  git commit -m "${COMMIT_MSG}"
  git remote add origin "git@github.com:${REPO_NAME}.git"
  git push -f origin "master:${TARGET_BRANCH}"
}

checksum-src() {
  echo "$(checksum-dir $DOCS_DIR)$(checksum-dir $WWW_DIR)" | md5sum | cut -d' ' -f1
}

checksum-dir() {
  find $1 -type f -exec md5sum {} \; | sort -k 2 | md5sum | cut -d' ' -f1
}

if [[ "$1" == "$EXPECT_BRANCH" ]]; then
  main "$@"
fi
