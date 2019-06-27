package vultr

import (
	"fmt"
	"log"
	"time"

	"github.com/JamesClonk/vultr/lib"
)

// waitForState simply blocks until the server is in a state we expect,
// while eventually timing out.
func waitForServerState(state string, power string, serverID string, client *lib.Client, timeout time.Duration) error {
	done := make(chan struct{})
	defer close(done)
	result := make(chan error, 1)
	go func() {
		attempts := 0
		for {
			attempts++
			log.Printf("Checking server status... (attempt: %d)", attempts)
			serverInfo, err := client.GetServer(serverID)
			if err != nil {
				result <- err
				return
			}
			if serverInfo.Status == state && (serverInfo.PowerStatus == power || power == "") {
				result <- nil
				return
			}

			time.Sleep(3 * time.Second)

			// Verify we shouldn't exit
			select {
			case <-done:
				// We finished, so just exit the goroutine
				return
			default:
				// Keep going
			}
		}
	}()
	log.Printf("Waiting for up to %d seconds for server", timeout/time.Second)
	select {
	case retval := <-result:
		return retval
	case <-time.After(timeout):
		err := fmt.Errorf("Timeout while waiting to for server")
		return err
	}
}

func waitForSnapshotState(state string, snapshotID string, client *lib.Client, timeout time.Duration) error {
	done := make(chan struct{})
	defer close(done)
	result := make(chan error, 1)
	go func() {
		attempts := 0
		for {
			attempts++
			log.Printf("Checking snapshot status... (attempt: %d)", attempts)

			// Crude workaround for the lack of singular GetSnapshot() method
			snapshots, err := client.GetSnapshots()
			if err != nil {
				result <- err
				return
			}

			found := false
			var snapshotInfo lib.Snapshot
			for _, s := range snapshots {
				if s.ID == snapshotID {
					snapshotInfo = s
					found = true
					break
				}
			}
			if !found {
				result <- fmt.Errorf("Snapshot is lost")
				return
			}

			if snapshotInfo.Status == state {
				result <- nil
				return
			}

			time.Sleep(3 * time.Second)

			// Verify we shouldn't exit
			select {
			case <-done:
				// We finished, so just exit the goroutine
				return
			default:
				// Keep going
			}
		}
	}()
	log.Printf("Waiting for up to %d seconds for snapshot", timeout/time.Second)
	select {
	case retval := <-result:
		return retval
	case <-time.After(timeout):
		err := fmt.Errorf("Timeout while waiting to for snapshot")
		return err
	}
}
