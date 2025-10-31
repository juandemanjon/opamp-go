package uisrv

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/open-telemetry/opamp-go/internal/examples/server/data"
	"github.com/open-telemetry/opamp-go/protobufs"
)

type CustomConfigForInstance struct {
	InstanceId  string `json:"instance_id"`
	Body        []byte `json:"body"`
	ContentType string `json:"content_type,omitempty"`
}

func apiSaveCustomConfigForInstance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var cs CustomConfigForInstance
	err := json.NewDecoder(r.Body).Decode(&cs)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(cs.InstanceId) == 0 {
		http.Error(w, "Empty instance_id", http.StatusBadRequest)
		return
	}
	if len(cs.Body) == 0 {
		http.Error(w, "Empty body", http.StatusBadRequest)
		return
	}

	uid, err := uuid.Parse(cs.InstanceId)
	if err != nil {
		http.Error(w, "Error parsing instance_id", http.StatusBadRequest)
		return
	}

	instanceId := data.InstanceId(uid)
	agent := data.AllAgents.GetAgentReadonlyClone(instanceId)
	if agent == nil {
		http.Error(w, "Cannot find agent by instance_id", http.StatusNotFound)
		return
	}

	config := &protobufs.AgentConfigMap{
		ConfigMap: map[string]*protobufs.AgentConfigFile{
			"": {Body: cs.Body, ContentType: cs.ContentType},
		},
	}

	notifyNextStatusUpdate := make(chan struct{}, 1)
	data.AllAgents.SetCustomConfigForAgent(instanceId, config, notifyNextStatusUpdate)

	// Wait for up to 5 seconds for a Status update, which is expected
	// to be reported by the Agent after we set the remote config.
	timer := time.NewTicker(time.Second * 5)

	select {
	case <-notifyNextStatusUpdate:
	case <-timer.C:
	}

	w.WriteHeader(http.StatusCreated)
}
