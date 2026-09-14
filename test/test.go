package main

import (
	"context"
	"log"
	"time"

	"github.com/thecguygithub/vtsgo"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	vts := vtsgo.New(".token", vtsgo.PluginInfo{
		PluginName:      "TestPlugin",
		PluginDeveloper: "TheCGuy",
	})
	err := vts.Connect(ctx, "ws://localhost:8001")
	if err != nil {
		log.Fatalf("Failed to Connect to VTS: %v", err)
	}
	log.Println("Connected to VTubeStudio!")
	status, err := vts.GetSessionStatus(ctx)
	if err != nil {
		log.Fatalf("Failed to get session status: %v", err)
	}
	log.Println("Response:")
	log.Printf("%+v", status)

	authenticateContext, _ := context.WithTimeout(context.Background(), 15*time.Second)
	err = vts.Authenticate(authenticateContext)
	if err != nil {
		log.Fatalf("Failed to Authenticate: %v", err)
	}

	log.Println("Connected to VTubeStudio!")
	status, err = vts.GetSessionStatus(ctx)
	if err != nil {
		log.Fatalf("Failed to get session status: %v", err)
	}
	log.Println("Response:")
	log.Printf("%+v", status)
}
