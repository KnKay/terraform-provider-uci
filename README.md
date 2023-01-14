# terraform-provider-uci
A terraform provider to configure openwrt devices


## About
To configure my infrastructure 100% in terraform I wanted to be able to use UCI commands. 
A quick search did not give any usefull provider. There is an ansible based thing but I tend to have my router in the infrastructure (for a big part at least). 

As I wanted to learn how to write a terraform provider this seems to be a good goal. 

## Actual goals

- [ ] Build a provider
- [ ] Set root password
- [ ] Handle backups (download upload)
- [ ] Handle package installs
- [ ] Configure WAN
- [ ] Configure WLAN
- [ ] Configure Wireguard
