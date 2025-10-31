package examples

import (
	"log"
	"testing"
	"time"

	"github.com/open-telemetry/opamp-go/internal/examples/agent/agent"
	"github.com/open-telemetry/opamp-go/internal/examples/server/data"
	"github.com/open-telemetry/opamp-go/internal/examples/server/opampsrv"
	"github.com/open-telemetry/opamp-go/protobufs"
)

func TestServerClient(t *testing.T) {

	serverLogger := log.Default()
	serverLogger.SetPrefix("SERVER ")

	opampSrv := opampsrv.NewServer(&opampsrv.Logger{Logger: serverLogger}, &data.AllAgents)
	opampSrv.Start()

	agentType := "io.opentelemetry.collector"
	agentVersion := "1.0.0"
	initialInsecureConnection := true

	agentLogger := log.Default()
	agentLogger.SetPrefix("AGENT ")

	agent := agent.NewAgent(&agent.Logger{Logger: agentLogger}, agentType, agentVersion, initialInsecureConnection)

	t.Logf("Agent InstanceId:%v", agent.InstanceId())

	config := &protobufs.AgentConfigMap{
		ConfigMap: map[string]*protobufs.AgentConfigFile{
			"": {Body: []byte("FooBar"), ContentType: ""},
		},
	}

	instanceId := data.InstanceId(agent.InstanceId())
	notifyNextStatusUpdate := make(chan struct{}, 1)
	data.AllAgents.SetCustomConfigForAgent(instanceId, config, notifyNextStatusUpdate)

	// Wait for up to 5 seconds for a Status update, which is expected
	// to be reported by the Agent after we set the remote config.
	timer := time.NewTicker(time.Second * 5)

	select {
	case <-notifyNextStatusUpdate:
	case <-timer.C:
	}

	time.Sleep(10 * time.Second)

	agent.Shutdown()
	opampSrv.Stop()

}
