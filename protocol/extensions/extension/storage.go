package extension

type newSnapshotStorage func(string) (SnapshotStorage, error)

type SnapshotStorage interface {
	StoreDirectory()
}

var registeredSnapshotStorages map[string]newSnapshotStorage

func RegisterSnapshotStorage(name string, ps newSnapshotStorage) {
	if registeredSnapshotStorages == nil {
		registeredSnapshotStorages = make(map[string]newSnapshotStorage)
	}
	registeredSnapshotStorages[name] = ps
}

func NewSnapshotStorage(name string, path string) (SnapshotStorage, error) {
	storage, ok := registeredSnapshotStorages[name]
	if !ok {
		return nil, ErrExtensionNotFound
	}
	storageInst, err := storage(path)
	if err != nil {
		return nil, err
	}
	return storageInst, nil
}
