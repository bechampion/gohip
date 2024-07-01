#!/usr/bin/env bash

# THIS SCRIPT IS PART OF gohip. DO NOT MODIFY IT UNLESS YOU KNOW WHAT YOU ARE DOING.

set -e

CONFIG_FILE=/etc/vpnc/splitvpn

if [[ ! -f $CONFIG_FILE ]]; then
    echo "$CONFIG_FILE does not exist. Split tunneling will not be active. Please create it with the following content if you want to activate:"
cat << EOF

# beginning
MAIN_DEV="enp0s31f6" # Your main network interface
GW="192.168.1.254"   # Your gateway
# end

You can determine those values with
  ip -json r get 1.1.1.1 | jq '.[]| "MAIN_DEV=\"\(.dev)\" \nGW=\"\(.gateway)\""' -r

EOF
    exit 0
fi

VPN_NETS=( 10/8 )
VPN_DEV="tun0"

. $CONFIG_FILE

ip route del default
ip route add default via $GW dev $MAIN_DEV

for subnet in "${VPN_NETS[@]}"
do
  ip route add $subnet dev $VPN_DEV
done


exit 0
