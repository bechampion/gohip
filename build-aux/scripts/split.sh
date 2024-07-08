#!/usr/bin/env bash

# THIS SCRIPT IS PART OF gohip. DO NOT MODIFY IT UNLESS YOU KNOW WHAT YOU ARE DOING.

set -e

if [ -z "$BASH" ]; then echo "Please run this script $0 with bash"; exit; fi

CONFIG_FILE=/etc/vpnc/splitvpn
RUNNING_FROM_OPENCONNECT=false

# When the script is called from vpnc, CISCO_DEF_DOMAIN is set.
if [[ -n "${CISCO_DEF_DOMAIN}" ]]; then
  RUNNING_FROM_OPENCONNECT=true
fi

if [[ ! -f $CONFIG_FILE ]]; then
    echo "[INFO] File $CONFIG_FILE does not exist."
    echo "[INFO] Manual split tunneling will not be active (not needed if split tunneling is provided by concentrator)."
    if [[ $RUNNING_FROM_OPENCONNECT == false ]]; then
      echo "[INFO] Please create it with the following content if you want to activate:"
      cat << EOF

# beginning
ENABLED=true
MAIN_DEV="enp0s31f6" # Your main network interface
GW="192.168.1.254"   # Your gateway
# end

You can determine the values for MAIN_DEV and GW with:
  ip -json r get 1.1.1.1 | jq '.[]| "MAIN_DEV=\"\(.dev)\" \nGW=\"\(.gateway)\""' -r

EOF
    fi
    exit 0
fi

ENABLED=false # override in $CONFIG_FILE
VPN_NETS=( 10/8 )

. $CONFIG_FILE

if [[ $ENABLED == true && $RUNNING_FROM_OPENCONNECT == true ]]; then
  ip route del default
  ip route add default via $GW dev $MAIN_DEV

  for subnet in "${VPN_NETS[@]}"
  do
    ip route add $subnet dev $TUNDEV # TUNDEV is set via openconnect
  done
else
  echo "Manual split tunneling is enabled only when running from vpnc (RUNNING_FROM_OPENCONNECT=$RUNNING_FROM_OPENCONNECT), and enabled in config (ENABLED=$ENABLED)"
fi

exit 0
