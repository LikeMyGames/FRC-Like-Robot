package subsystem

type (
	Subsystem interface {
		// SetState(string)
		Initialize()
		Periodic()
	}
)
