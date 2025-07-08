package node

import (
	"fmt"
	"os"
	"runtime"
	"time"

	pb "github.com/Friends-Of-Noso/NosoGo/protobuf"
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

// Generates block zero
func getBlockZero() *pb.Block {
	block := &pb.Block{
		Hash:         "COINBASE",
		Height:       0,
		PreviousHash: "BZERO",
		Timestamp:    time.Now().Unix(), // TODO: This should be dated to genesis
		MerkleRoot:   "MZERO",
	}
	// block.SetHash()
	return block
}
