package api

import (
	"net/http"

	"github.com/Dataman-Cloud/swan/mesos"
	"github.com/Dataman-Cloud/swan/mole"
	"github.com/Dataman-Cloud/swan/types"
)

type Driver interface {
	KillTask(string, string) error
	LaunchTasks([]*mesos.Task) (map[string]error, error)

	ClusterName() string

	SubscribeEvent(http.ResponseWriter, string) error
	FullTaskEventsAndRecords() []*types.CombinedEvents

	ClusterAgents() map[string]*mole.ClusterAgent
	ClusterAgent(id string) *mole.ClusterAgent

	// for debug convenience
	Dump() interface{}
	Load() map[string]interface{}
}
