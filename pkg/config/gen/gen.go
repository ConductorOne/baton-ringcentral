package main

import (
	cfg "github.com/conductorone/baton-ringcentral/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("ringcentral", cfg.Config)
}
