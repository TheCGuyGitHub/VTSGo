# VTSGo
A Go library for communicating with the VTubeStudio API

## Installation
```
go get github.com/thecguygithub/vtsgo
```

## Usage
```
package main

import (
	"context"
	"log"
	"time"

	"github.com/thecguygithub/vtsgo"
)

func main() {
	vts := vtsgo.New(".token", vtsgo.PluginInfo{
		PluginName:      "TestPlugin",
		PluginDeveloper: "TheCGuy",
	})
	err := vts.Connect(context.Background(), "ws://localhost:8001")
	if err != nil {
		log.Fatalf("Failed to Connect to VTS: %v", err)
	}
	log.Println("Connected to VTubeStudio!")

	authenticateContext, _ := context.WithTimeout(context.Background(), 15*time.Second)
	err = vts.Authenticate(authenticateContext)
	if err != nil {
		log.Fatalf("Failed to Authenticate: %v", err)
	}

	log.Println("Authenticated with VTubeStudio!")
	resp, err := vts.CheckFaceFound(ctx)
	if err != nil {
		log.Fatalf("Failed to CheckFaceFound: %v", err)
	}
	log.Printf("FaceFound: %v", resp.Found)

	modelResp, err := vts.GetCurrentModel(ctx)
	if err != nil {
		log.Fatalf("Failed to GetCurrentModel: %v", err)
	}
	log.Printf("CurrentModel: %+v", modelResp)
}

```