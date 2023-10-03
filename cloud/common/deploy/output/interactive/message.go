package interactive

type ResourceState string

const (
	Created  ResourceState = "Created"
	Creating ResourceState = "Creating"
	Updating ResourceState = "Updating"
	Updated  ResourceState = "Updated"
	Deleting ResourceState = "Deleting"
	Deleted  ResourceState = "Deleted"
)

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
