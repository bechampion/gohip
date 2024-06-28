gohip
=====

The HIP ( `Host Integrity Protection`) mechanism is a security scanner for the Palo Alto Networks GlobalProtect VPNs, in the same vein as Cisco's CSD and Juniper's Host Checker ([source](https://www.infradead.org/openconnect/hip.html)).

# Install

Download and install `gohip` from the [releases page](https://github.com/bechampion/gohip/releases).

# Usage

Create file `/etc/vpnc/splitvpn` with the following content:

    MAIN_DEV="enp0s31f6" # Your main network interface
    GW="192.168.1.254"   # Your gateway

You can determine those values with

    ip -json r get 1.1.1.1 | jq '.[]| "MAIN_DEV=\"\(.dev)\" \nGW=\"\(.gateway)\""' -r

Then start the vpn client with

    sudo -E gpclient connect --hip --csd-wrapper /usr/bin/gohip vpn.endpoint.com
	
# Using your default browser (so creds are remembered , hsm etc)

You can pass `--default-browser` to gpclient is your in version `2.3.2` at least ,  and that should open on what xdg-config says
```
sudo -E gpclient connect --hip --csd-wrapper /home/jgarcia/Projects/disney/xgohip iad1gwavo.gp.disneystreaming.com --default-browser
```

To set your default browser you can do
```
xdg-settings set default-web-browser firefox.desktop
```


## Notes
If you use firefox/chrome from snaps this will not work ,  gpclient drops a file on `/tmp` that needs to be reachable by the browser and that doesn't work with snaps
