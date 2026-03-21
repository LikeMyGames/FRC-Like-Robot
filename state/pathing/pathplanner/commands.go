package pathplanner

var (
	namedCommands map[string]*NamedCommand
)

func NewCommand(function func(any), argument any) *NamedCommand {
	return &NamedCommand{
		Function: function,
		Argument: argument,
	}
}

func RegisterCommand(name string, command *NamedCommand) {
	namedCommands[name] = command
}

func HasCommand(name string) bool {
	for s := range namedCommands {
		if s == name {
			return true
		}
	}
	return false
}

func RunCommand(name string) {
	if HasCommand(name) {
		namedCommands[name].Function(namedCommands[name].Argument)
	}
}
