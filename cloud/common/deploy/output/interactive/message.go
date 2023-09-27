package interactive

import (
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
)

type ResourceState string

const (
	Created  ResourceState = "Created"
	Creating ResourceState = "Creating"
	Updating ResourceState = "Updating"
	Updated  ResourceState = "Updated"
	Deleting ResourceState = "Deleting"
	Deleted  ResourceState = "Deleted"
)

type ResourceUpdateMessage struct {
	Name  string
	Type  v1.ResourceType
	State ResourceState
}

type LogMessage struct {
	Message string
}

type LogMessageSubscriptionWriter struct {
	Sub chan LogMessage
}

func (l LogMessageSubscriptionWriter) Write(b []byte) (int, error) {
	l.Sub <- LogMessage{
		Message: string(b),
	}

	return len(b), nil
}
