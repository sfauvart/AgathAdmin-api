#!/usr/bin/env bash

NONE='\033[00m'
RED='\033[01;31m'
GREEN='\033[01;32m'
YELLOW='\033[01;33m'
PURPLE='\033[01;35m'
CYAN='\033[01;36m'
WHITE='\033[01;37m'
BOLD='\033[1m'
UNDERLINE='\033[4m'

version_lt() {
  test "$(echo "$@" | tr " " "\n" | sort -rV | head -n 1)" != "$1";
}

_init() {
  ## Minimum required versions for build dependencies
  GIT_VERSION="1.0"
  GO_VERSION="1.6"
  OPENSSL_VERSION="1.0.2"
  UNAME=$(uname -sm)

  ## Check all dependencies are present
  MISSING=""
}

check_deps() {
  echo -n "Golang (>${GO_VERSION}): "
   if version_lt "$(env go version 2>/dev/null | sed 's/^.* go\([0-9.]*\).*$/\1/')" "${GO_VERSION}"; then
     MISSING="${MISSING} golang(>${GO_VERSION})"
     echo -e "${RED}KO${NONE}"
   else
     echo -e "${YELLOW}ok${NONE}"
   fi

  echo -n "Git (>${GIT_VERSION}): "
   if version_lt "$(env git --version 2>/dev/null | sed -e 's/^.* \([0-9.\].*\).*$/\1/' -e 's/^\([0-9.\]*\).*/\1/g')" "${GIT_VERSION}"; then
     MISSING="${MISSING} git(>${GIT_VERSION})"
     echo -e "${RED}KO${NONE}"
   else
     echo -e "${YELLOW}ok${NONE}"
   fi

   echo -n "OpenSSL (>${OPENSSL_VERSION}): "
    if version_lt "$(env openssl version 2>/dev/null | sed 's/^OpenSSL \([0-9.]*\).*$/\1/')" "${OPENSSL_VERSION}"; then
      MISSING="${MISSING} openssl(>${OPENSSL_VERSION})"
      echo -e "${RED}KO${NONE}"
    else
      echo -e "${YELLOW}ok${NONE}"
    fi
}

main() {
  echo -e "${GREEN}Checking build tools...${NONE}"
  check_deps

  ## If dependencies are missing, warn the user and abort
  if [ "x${MISSING}" != "x" ]; then
     echo -e "${RED}"
     echo "ERROR"
     echo
     echo "The following build tools are missing or not up to date:"
     echo
     echo "** ${MISSING} **"
     echo -e "${NONE}"
     exit 1
  fi
  echo "check_tools Done !"
}

_init && main "$@"
