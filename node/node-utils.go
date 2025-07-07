package node

import (
	"fmt"
	"os"
	"runtime"
)

// Checks for privileges of the user running the application.
// If ran under Linux/Darwin and port < 1024, user has to be root
func checkPort(port int, flag string, defaultPort int) error {
	goos := runtime.GOOS

	if goos == "linux" || goos == "darwin" {
		if port < 1024 && os.Geteuid() != 0 {
			return fmt.Errorf("port %d requires root privileges on %s; try a port >= 1024 (e.g. --%s %d)", port, goos, flag, defaultPort)
		}
	}

	return nil
}
