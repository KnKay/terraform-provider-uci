# Network

Networking is the obviously most important part of a network device.
The networking is having a WAN part, which is our upstream as well as internal things.

Up to now we have the WAN part.

## WAN
The WAN can be configured in static pr DHCP mode.
A default output of uci on a vanilla OpenWrt is:

```
root@OpenWrt:~# uci show network
network.loopback=interface
network.loopback.device='lo'
network.loopback.proto='static'
network.loopback.ipaddr='127.0.0.1'
network.loopback.netmask='255.0.0.0'
network.globals=globals
network.globals.ula_prefix='fde7:42dd:244b::/48'
network.@device[0]=device
network.@device[0].name='br-lan'
network.@device[0].type='bridge'
network.@device[0].ports='eth0'
network.wan=interface
network.wan.device='eth1'
network.wan.proto='dhcp'
network.mng=interface
network.mng.device='eth0'
network.mng.proto='dhcp'
root@OpenWrt:~#
```
