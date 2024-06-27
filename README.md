gohip
=====

The HIP ( `Host Integrity Protection`) mechanism is a security scanner for the Palo Alto Networks GlobalProtect VPNs, in the same vein as Cisco's CSD and Juniper's Host Checker ([source](https://www.infradead.org/openconnect/hip.html)).

# Install

Download and install `gohip` from the [releases page](https://github.com/bechampion/gohip/releases).

# Usage

    sudo -E gpclient connect --hip --csd-wrapper /usr/bin/gohip vpn.endpoint.com