package main

import (
  "log"
  "encoding/json"
  "io/ioutil"
)

type Config struct {
  Pin       string   `json:"pin"`
  Interface string   `json:"interface"`
  Buttons   []Button `json:"buttons"`
}

type Button struct {
  Name string `json:"name"`
  Mac  string `json:"mac"`
}

func GetConfig() *Config {
  data, err := ioutil.ReadFile("dashbridge.json")
  if err != nil {
    log.Fatalf("Config read failed: %v", err)
  }

  config := &Config{}
  err = json.Unmarshal(data, &config)

  if err != nil {
    log.Fatalf("Config load failed: %v", err)
  }

  return config
}
