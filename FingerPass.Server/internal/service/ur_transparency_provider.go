package service

import (
  // "log"

  // "github.com/matwate/corner/internal/config"
)

type UrTransparencyProvider struct{}

func (mt UrTransparencyProvider) Accepts(key string) bool {
	return true
}

func (mt UrTransparencyProvider) Accept(key, mime string, data []byte) error {
	// log.Printf("Accepting %s\n", key)
	return nil
}
