package pi

import (
	"fmt"
	"io"
	"net"
	"os"
)

var (
	debugMode   = false
	debugOutput = io.Writer(os.Stdout)
)

// SetDebug sets the debugMode, logging every requests and pretty prints JSON and XML.
func SetDebug(debug bool) {
	debugMode = debug
}

// SetDebugOutput sets where we need to write the output of the debug.
// By default, it is the standard output.
func SetDebugOutput(writer io.Writer) {
	debugOutput = writer
}

// writeDebug writes debug string formatted as: [GET] to/from [IP address] debug_message.
// If you are working on localhost and your machine is using IPV6 addresses, you'll get ::1.
func writeDebug(method, remoteAddr, output string) {
	ip, _, _ := net.SplitHostPort(remoteAddr)
	fmt.Fprintf(debugOutput, "[%s] to/from [%s] %s\n", method, ip, output)
}
