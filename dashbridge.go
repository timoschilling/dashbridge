package main

import (
  "log"
  "fmt"
  "github.com/brutella/hc"
  "github.com/brutella/hc/accessory"
  "github.com/google/gopacket"
  "github.com/google/gopacket/layers"
  "github.com/google/gopacket/pcap"
)

func FindButton(identifier string, switches []*accessory.Switch) (*accessory.Switch, bool) {
  for _, switch_accessory := range switches {
    if switch_accessory.Info.SerialNumber.GetValue() == identifier {
      return switch_accessory, true
    }
  }
  return nil, false
}

func main() {
  config := GetConfig()
  accessories := []*accessory.Accessory{}
  switches := []*accessory.Switch{}

  for _, button := range config.Buttons {
    switch_accessory := accessory.NewSwitch(accessory.Info{
      Name:         button.Name,
      Manufacturer: "Amazon",
      SerialNumber: button.Mac,
      Model:        fmt.Sprintf("Dash Button %s", button.Name),
    })
    switches = append(switches, switch_accessory)
    accessories = append(accessories, switch_accessory.Accessory)
  }

  t, err := hc.NewIPTransport(hc.Config{Pin: config.Pin, StoragePath: "database"}, NewBridge().Accessory, accessories...)
  if err != nil {
    log.Fatal(err)
  }

  go func() {
    hc.OnTermination(func() {
      <-t.Stop()
    })

    t.Start()
  }()

  log.Printf("Starting up on interface[%v]...", config.Interface)

  h, err := pcap.OpenLive(config.Interface, 65536, true, pcap.BlockForever)

  if err != nil || h == nil {
    log.Fatalf("Error opening interface: %s\nPerhaps you need to run as root?\n", err)
  }
  defer h.Close()

  err = h.SetBPFFilter("arp")
  if err != nil {
    log.Fatalf("Unable to set filter! %s\n", err)
  }
  log.Println("Listening for Dash buttons...")

  packetSource := gopacket.NewPacketSource(h, h.LinkType())

  for packet := range packetSource.Packets() {
    ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
    ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)

    button, err := FindButton(ethernetPacket.SrcMAC.String(), switches)
    if err {
      log.Printf("Found %s\n", button.Info.Name.GetValue())
      button.Switch.On.SetValue(!button.Switch.On.GetValue())
    } else {
      log.Printf("Unable to find %s\n", ethernetPacket.SrcMAC.String())
    }
  }
}