package internal

func freeSnapshotHandles(snaps []snapshotHandle) {
	for _, s := range snaps {
		s.Free()
	}
}

func findSnapshotByName(snaps []snapshotHandle, snapName string) snapshotHandle {
	for _, s := range snaps {
		n, err := s.Name()
		if err != nil {
			continue
		}
		if n == snapName {
			return s
		}
	}
	return nil
}
