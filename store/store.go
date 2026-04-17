package store

import (
	"bytes"
	"os"
	"strings"

	"github.com/evolvedevlab/weaveset/data"
)

type Storer interface {
	Save(*data.List) error
}

type FileSystemStore struct {
	file         *os.File
	storeDirPath string
}

func NewFileSystemStore(dirPath string) Storer {
	return &FileSystemStore{
		storeDirPath: dirPath,
	}
}

func (s *FileSystemStore) Save(list *data.List) error {
	if s.file == nil {
		var err error

		filename := strings.Join(strings.Split(strings.ToLower(list.Name), " "), "_")
		s.file, err = os.Create(s.storeDirPath + "/" + filename + ".md")
		if err != nil {
			return err
		}
	}

	buf := new(bytes.Buffer)

	buf.WriteString("---\n")
	buf.WriteString("title: " + list.Name + "\n")
	buf.WriteString("draft: false\n")
	buf.WriteString("---\n")
	buf.WriteString("content heyyy")

	_, err := s.file.Write(buf.Bytes())
	return err
}
