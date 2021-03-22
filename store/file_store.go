package store

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/rlaskowski/go-email/config"
	"github.com/rlaskowski/go-email/model"
)

type FileStore struct {
	storePath      string
	poolSerializer sync.Pool
	fileMutex      sync.RWMutex
}

func NewFileStore(storePath string) *FileStore {
	f := &FileStore{
		storePath: storePath,
	}

	return f
}

/* func NewFileStore(file *datastore.FileSystem) *FileStore {
	FileStore := &FileStore{
		file: file,
	}
	FileStore.poolSerializer.New = func() interface{} {
		return serialization.NewSerializer()
	}
	return FileStore
}

func (f *FileStore) Base(hash string) string {
	fileInfo, err := f.readMetadata(hash)
	if err != nil {
		return ""
	}

	return fileInfo.OriginalName
}*/

func (f *FileStore) Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

/*
func (f *FileStore) EncodeName(name string) (string, error) {
	encode := base64.RawURLEncoding.EncodeToString([]byte(name))
	if !(len(encode) > 0) {
		return "", errors.New(("Bad encode file id"))
	}

	return encode, nil
}

 func (f *FileStore) DecodeName(hash string) string {
	decode, err := base64.RawURLEncoding.DecodeString(hash)
	if err != nil {
		return ""
	}

	return string(decode)
} */

func (f *FileStore) ControllDir(uuid string) string {
	var directoryName, path string

	if len(uuid) > 4 {
		nextToLastLen := 2
		offset := len(uuid) - nextToLastLen - 1
		directoryName = uuid[offset : offset+nextToLastLen]
		path = filepath.Join(f.storePath, directoryName)
	}

	return path
}

func (f *FileStore) FindByUUID(uuid string) (*model.File, error) {
	path := filepath.Join(f.ControllDir(uuid), uuid)
	stat, err := os.Stat(path)

	if err != nil {
		return nil, err
	}

	return &model.File{
		Name:    uuid,
		Size:    stat.Size(),
		ModTime: stat.ModTime(),
	}, nil
}

/*
func (f *FileStore) FindByPath(path string) (*model.File, error) {
	file, err := f.FindByUUID(uuid)

	if err != nil {
		return nil, err
	}

	file.Name = filepath.Join(f.Dir(uuid), uuid)

	return file, nil
} */

func (f *FileStore) Store(reader io.Reader) (string, error) {
	id := f.fileId()
	path := filepath.Join(f.ControllDir(id), id)
	fmt.Println(f.ControllDir(id))
	if !f.Exists(f.ControllDir(id)) {
		err := os.Mkdir(f.ControllDir(id), config.FilePermissions)

		if err != nil {
			return "", fmt.Errorf("Error before create control folder %s", path)
		}
	}

	return id, f.store(path, reader)
}

func (f *FileStore) store(path string, reader io.Reader) error {
	destinationFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, config.FilePermissions)
	defer destinationFile.Close()

	if err != nil {
		return err
	}

	buff := make([]byte, config.FileCopyBuff)
	_, err = io.CopyBuffer(destinationFile, reader, buff)

	if err != nil {
		if err := f.Remove(path); err != nil {
			return err
		}
		return err
	}

	return nil
}

func (f *FileStore) Remove(file string) error {
	err := os.Remove(file)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileStore) fileId() string {
	uuid := uuid.New()
	return uuid.String()
}

/* func (f *FileStore) metadataName(hash string) string {
	name := fmt.Sprintf(".%s.data", hash)
	return name
}

func (f *FileStore) metadataPath(hash string) string {
	name := f.metadataName(hash)
	path := filepath.Join(f.Dir(hash), name)
	return path
}

func (f *FileStore) createMetadata(hash string, fileInfo *model.FileInfo) error {
	path := f.metadataPath(hash)

	data, err := os.Create(path)
	defer data.Close()

	if err != nil {
		return err
	}

	if err := f.acquireGobSerializer().Serialize(data, fileInfo); err != nil {
		return err
	}

	if err := config.HideFile(path); err != nil {
		return err
	}

	return nil
}

func (f *FileStore) readMetadata(hash string) (*model.FileInfo, error) {
	f.fileMutex.Lock()
	defer f.fileMutex.Unlock()

	path := f.metadataPath(hash)

	data, err := os.Open(path)
	defer data.Close()

	if err != nil {
		return nil, err
	}

	fileInfo := new(model.FileInfo)

	if err := f.acquireGobSerializer().Deserialize(data, fileInfo); err != nil {
		return nil, err
	}
	return fileInfo, nil
}

func (f *FileStore) remove(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}

	metadataPath := f.metadataPath(filepath.Base(path))
	err = os.Remove(metadataPath)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileStore) acquireGobSerializer() *serialization.GobSerializer {
	s := f.poolSerializer.Get().(*serialization.GobSerializer)
	defer f.poolSerializer.Put(s)

	return s
}*/
