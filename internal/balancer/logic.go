package balancer

import (
	"time"
)

type State int

const (
	Set State = iota
	Scaling
	Idle
)

var idleEnd time.Time
var maxCluster string
var maxReplicas int32
var state State = Set

func logic(replicas map[string]int32, idle int) (State, []string) {
	if state == Set {
		for cluster, reps := range replicas {
			if reps > maxReplicas {
				maxCluster = cluster
				maxReplicas = reps
			}
		}
		state = Scaling
	}
	if state == Scaling {
		scaleClusters := []string{}
		for cluster, reps := range replicas {
			if reps < maxReplicas {
				scaleClusters = append(scaleClusters, cluster)
			}
		}
		if len(scaleClusters) > 0 {
			return Scaling, scaleClusters
		}
		idleEnd = time.Now().Add(time.Duration(idle) * time.Minute)
		state = Idle
		return Idle, []string{}
	}
	// CASE IDLE
	now := time.Now()
	if now.After(idleEnd) {
		state = Set
		return Set, []string{}
	}
	return Idle, []string{}
}
