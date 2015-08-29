#!/bin/bash

set -e

function str_to_array {
  eval "local input=\"\$$1\""
  input="$(echo "$input" | awk '
  {
    split($0, chars, "")
    for (i = 1; i <= length($0); i++) {
      if (i > 1) {
        printf(", ")
      }
      printf("'\''%s'\''", chars[i])
    }
  }
  ')"
  eval "$1=\"$input\""
}

function update_access_key {
  str_to_array LOGIN_ALIYUNID_TICKET
  str_to_array FLOWDOCK_TOKEN
  awk "
  /LOGIN_ALIYUNID_TICKET/ {
    print \"var LOGIN_ALIYUNID_TICKET = []byte{${LOGIN_ALIYUNID_TICKET}}\"
    next
  }
  /FLOWDOCK_TOKEN/ {
    print \"var FLOWDOCK_TOKEN = []byte{${FLOWDOCK_TOKEN}}\"
    next
  }
  {
    print
  }
  " access.go > _access.go

  mv _access.go access.go
}

while test -z "$LOGIN_ALIYUNID_TICKET"; do
  echo -n "Please paste cookie 'login_aliyunid_ticket': "
  read -s LOGIN_ALIYUNID_TICKET
  echo
done
while test -z "$FLOWDOCK_TOKEN"; do
  echo -n "Please paste flowdock token: "
  read -s FLOWDOCK_TOKEN
  echo
done
update_access_key

if test -n "$BUILD_DOCKER"; then
  docker-compose up build
  docker-compose rm --force -v
  docker-compose build app
else
  go build
fi

LOGIN_ALIYUNID_TICKET=""
FLOWDOCK_TOKEN=""
update_access_key
