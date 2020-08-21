# vpngater

Fetches preferred !!FREE!! highspeed OpenVPN hosts file from VPNGate project and configures the same for your machine.

## Usage of vpngater:
    -latency int
        Minimum latency of VPN (ms) (default 40)
    -minSpeed int
        Minimum speed of VPN (Mbps) (default 40)
    -openvpnBin string
        Mention path of OpenVPN binary (default "/usr/sbin/openvpn")
     -proto string
        Preferred Protocol - tcp/udp/any (default "any")

## Downloads
Releases - 
https://github.com/prben/vpngatefetcher/releases


## Linux Priveleges
On Linux systems, you might need to run vpngater as root to allow it to run openvpn binary as root privileges.
