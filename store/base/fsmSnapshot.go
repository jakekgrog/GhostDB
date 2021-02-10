package base

import (
	"log"
	"io"
	
	"github.com/hashicorp/raft"
)
type fsmSnapshot struct{}

// Snapshot returns a snapshot of the key-value store.
func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	// TODO: implement
	log.Printf("snapshot")
	return &fsmSnapshot{}, nil
}

// Restore stores the key-value store to a previous state.
func (f *fsm) Restore(rc io.ReadCloser) error {
	// TODO: implement
	log.Printf("restore [%v]", rc)
	return nil
}

// Persist is called on an fsmSnapshot and is used to write to file
// parameters: (a raft.SnapshotSink used to write snapshots)
// returns: error
func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	// TODO: implement
	err := func() error {
		b := []byte("hello from persist")
		if _, err := sink.Write(b); err != nil {
			return err
		}
		return sink.Close()
	}()
	if err != nil {
		sink.Cancel()
	}
	return err
}

func (f *fsmSnapshot) Release() {}