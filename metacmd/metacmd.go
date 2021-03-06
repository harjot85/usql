package metacmd

// Metacmd represents a command and associated meta information about it.
type Metacmd uint

// Decode converts a command name (or alias) into a Runner.
func Decode(name string, params []string) (Runner, error) {
	mc, ok := cmdMap[name]
	if !ok {
		return nil, ErrUnknownCommand
	}

	cmd := cmds[mc]
	if cmd.Min > len(params) {
		return nil, ErrMissingRequiredArgument
	}

	return RunnerFunc(func(h Handler) (Res, error) {
		return cmd.Process(h, name, params)
	}), nil
}

// Command types.
const (
	// None is an empty command.
	None Metacmd = iota

	// Question is question meta command (\?)
	Question

	// Quit is the quit meta command (\?).
	Quit

	// Copyright is the copyright meta command (\copyright).
	Copyright

	// Connect is the connect meta command (\c).
	Connect

	// Disconnect is the disconnect meta command (\Z).
	Disconnect

	// ConnInfo is the connection info meta command (\conninfo).
	ConnInfo

	// Drivers is the driver info meta command (\drivers).
	Drivers

	// Describe is the describe meta command (\d and variants).
	Describe

	// Exec is the execute meta command (\g and variants).
	Exec

	// Edit is the edit query buffer meta command (\e).
	Edit

	// Print is the print query buffer meta command (\p).
	Print

	// Reset is the reset query buffer meta command (\r).
	Reset

	// Echo is the echo meta command (\echo).
	Echo

	// Write is the write meta command (\w).
	Write

	// ChangeDir is the system change directory meta command (\cd).
	ChangeDir

	// SetEnv is the system set environment variable meta command (\setenv).
	SetEnv

	// Include is the system include file meta command (\i and variants).
	Include
)
