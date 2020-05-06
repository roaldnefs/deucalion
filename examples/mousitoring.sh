#!/usr/bin/env bash

TIMESTAMP=$(date +%s)
STATEFILE=/tmp/mousitoring.state
FORCERUN=""

if [ ! -f "${STATEFILE}" ]; then
  touch "${STATEFILE}"
  FORCERUN="1"
else
  STATESTAMP=$(cut -d' ' -f1 < ${STATEFILE})
  if [[ $(( TIMESTAMP - STATESTAMP )) -ge 60 ]]; then
    FORCERUN="1"
  fi
fi

case $1 in
  firing)
    COLOUR=red
    LIGHT_EFFECT=4
    shift;;
  warning)
    COLOUR=yellow
    LIGHT_EFFECT=3
    shift;;
  silent)
    COLOUR=green
    LIGHT_EFFECT=2
    shift;;
  *)
    echo "[!] $(date +%F) :: Please enter a valid option: firing|warning|silent"
    exit 1
    shift;;
esac

if [ "${FORCERUN}" ]; then
  echo "[+] $(date +%F) :: Enforcing run because statefule didn't exist or time difference too big"
  rivalcfg -c "${COLOUR}" -e "${LIGHT_EFFECT}"
  echo "${TIMESTAMP} ${COLOUR}" > "${STATEFILE}"
  exit 0
fi

if [ "$(grep -c ${COLOUR} ${STATEFILE})" -ge 1 ]; then
  echo "[+] $(date +%F) :: Colour already set properly, no need to take action"
  echo "${TIMESTAMP} ${COLOUR}" > "${STATEFILE}"
  exit 0
else
  echo "[+] $(date +%F) :: Colour not set properly, changing colour"
  rivalcfg -c "${COLOUR}" -e ${LIGHT_EFFECT}
  echo "${TIMESTAMP} ${COLOUR}" > "${STATEFILE}"
  exit 0
fi
