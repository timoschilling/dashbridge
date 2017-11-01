package main

import (
  "github.com/brutella/hc/accessory"
)

type Bridge struct {
  *accessory.Accessory
}

func NewBridge() *Bridge {
  acc := Bridge{}
  info := accessory.Info{
    Name:         "DashBridge",
    Manufacturer: "Amazon",
    Model:        "DashBridge",
  }
  acc.Accessory = accessory.New(info, accessory.TypeBridge)

  return &acc
}
