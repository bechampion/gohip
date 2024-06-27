gohip
=====

The HIP ( `Host Integrity Protection`) mechanism is a security scanner for the Palo Alto Networks GlobalProtect VPNs, in the same vein as Cisco's CSD and Juniper's Host Checker ([source](https://www.infradead.org/openconnect/hip.html)).

# Install

Download and install `gohip` from the [releases page](https://github.com/bechampion/gohip/releases).

# Usage

Create a file in your home directory called `.splitvpn` with the following content:

    MAIN_DEV="enp0s31f6" # Your main network interface
    GW="192.168.1.254"   # Your gateway

You can determine those values with

    ip -json r get 1.1.1.1 | jq '.[]| "DEV=\"\(.dev)\" \nGW=\"\(.gateway)\""' -r

Then start the vpn client with

    sudo -E gpclient connect --hip --csd-wrapper /usr/bin/gohip vpn.endpoint.com