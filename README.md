# terraform-provider-uci
A terraform provider to configure openwrt devices


## About
To configure my infrastructure 100% in terraform I wanted to be able to use UCI commands.
A quick search did not give any usefull provider. There is an ansible based thing but I tend to have my router in the infrastructure (for a big part at least).

As I wanted to learn how to write a terraform provider this seems to be a good goal.

## Actual goals

- [x] Build a provider
- [x] Set hostname
- [x] Configure WAN
- [ ] Handle backups (download upload)
- [ ] Configure strongswan
- [ ] Configure WLAN
- [ ] Set root password
- [ ] Handle package installs


## Links
[How to write a provider](https://developer.hashicorp.com/terraform/plugin/framework)
[Hashicups example](https://github.com/hashicorp/terraform-provider-hashicups-pf)
[fork of an UCI go lib](github.com/KnKay/go-uci)
