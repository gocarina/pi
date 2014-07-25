package pi

import (
	"fmt"
	"io"
	"os"
)

var (
	debugMode   = false
	debugOutput = io.Writer(os.Stdout)
)

// SetDebug sets the debugMode, logging every requests and pretty prints JSON and XML.
func SetDebug(debug bool) { debugMode = debug }

// SetDebugOutput sets where we need to write the output of the debug.
// By default, it is the standard output.
func SetDebugOutput(writer io.Writer) { debugOutput = writer }

func writeDebug(method, remoteAddr, output string) { fmt.Fprintf(debugOutput, "[%s] to [%s] %s", method, remoteAddr, output) }
