package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/open-telemetry/opamp-go/internal/examples/agent/agent"
)

func main() {
	var agentType string
	flag.StringVar(&agentType, "t", "io.opentelemetry.collector", "Agent Type String")

	var agentVersion string
	flag.StringVar(&agentVersion, "v", "1.0.0", "Agent Version String")

	var initialInsecureConnection bool
	flag.BoolVar(&initialInsecureConnection, "initial-insecure-connection", false, "Set SkipInsecureVerify for the initial connection to the OpAMP server.")

	var serverHost string
	flag.StringVar(&agentVersion, "server-host", "127.0.0.1", "Server Host String")

	flag.Parse()

	agent := agent.NewAgent(&agent.Logger{log.Default()}, agentType, agentVersion, initialInsecureConnection, serverHost)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
	agent.Shutdown()
}
