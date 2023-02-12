terraform {
  required_providers {
    uci = {
      source = "hashicorp.com/edu/uci"
    }
  }
}

provider "uci" {}

data "uci_system" "example" {}