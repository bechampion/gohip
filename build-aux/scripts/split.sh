#!/usr/bin/env bash

set -e

if [[ ! -f $HOME/.splitvpn ]]; then
    echo "$HOME/.splitvpn does exist. Please create it with the following content:"
cat << EOF

# beginning
MAIN_DEV="enp0s31f6" # Your main network interface
GW="192.168.1.254"   # Your gateway
# end

You can determine those values with
  ip -json r get 1.1.1.1 | jq '.[]| "DEV=\"\(.dev)\" \nGW=\"\(.gateway)\""' -r

EOF
    exit 1
fi

. ${HOME}/.splitvpn

DISNEY_NET="10/8"
VPN_DEV="tun0"

ip route del default
ip route add default via $GW dev $MAIN_DEV
ip route add $DISNEY_NET dev $VPN_DEV

exit 0
