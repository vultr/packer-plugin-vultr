package vultr

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/vultr/govultr/v2"
)

func waitForISOState(state string, isoID string, client *govultr.Client, timeout time.Duration) error {
	done := make(chan struct{})
	defer close(done)
	result := make(chan error, 1)
	go func() {
		attempts := 0
		for {
			attempts++
			log.Printf("Checking ISO status... (attempt: %d)", attempts)

			iso, err := client.ISO.Get(context.Background(), isoID)
			if err != nil {
				result <- err
				return
			}

			if iso.Status == state {
				result <- nil
				return
			}

			time.Sleep(3 * time.Second)

			select {
			case <-done:
				return
			default:
			}
		}
	}()
	log.Printf("Waiting for up to %d seconds for ISO", timeout/time.Second)
	select {
	case err := <-result:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("timeout while waiting for ISO")
	}
}

// waitForState simply blocks until the server is in a state we expect,
// while eventually timing out.
func waitForServerState(state string, power string, serverID string, client *govultr.Client, timeout time.Duration) error {
	done := make(chan struct{})
	defer close(done)
	result := make(chan error, 1)
	go func() {
		attempts := 0
		for {
			attempts++
			log.Printf("Checking server status... (attempt: %d)", attempts)
			serverInfo, err := client.Instance.Get(context.Background(), serverID)
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
	case err := <-result:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("timeout while waiting for server")
	}
}

func waitForSnapshotState(state string, snapshotID string, client *govultr.Client, timeout time.Duration) error {
	done := make(chan struct{})
	defer close(done)
	result := make(chan error, 1)
	go func() {
		attempts := 0
		for {
			attempts++
			log.Printf("Checking snapshot status... (attempt: %d)", attempts)

			snapshot, err := client.Snapshot.Get(context.Background(), snapshotID)
			if err != nil {
				result <- err
				return
			}

			if snapshot.Status == state {
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
	case err := <-result:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("timeout while waiting for snapshot")
	}
}
