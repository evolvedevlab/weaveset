package store

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/evolvedevlab/weaveset/data"
	"github.com/evolvedevlab/weaveset/util"
)

type FilepathGeneratorFunc func(*data.List) string

var DefaultFilepathGeneratorFunc = func(dirPath string) FilepathGeneratorFunc {
	return func(list *data.List) string {
		filename := fmt.Sprintf("%s-%s.md", util.GenerateSlug(list.Name, "-"), list.ID)
		return filepath.Join(dirPath, filename)
	}
}

type Storer interface {
	Save(*data.List) error
}

type FileSystem struct {
	DirPath           string
	FilepathGenerator FilepathGeneratorFunc
}

func NewFileSystem(dirPath string) *FileSystem {
	return &FileSystem{
		DirPath:           dirPath,
		FilepathGenerator: DefaultFilepathGeneratorFunc(dirPath),
	}
}

func (s *FileSystem) Save(list *data.List) error {
	filepath := s.FilepathGenerator(list)

	file, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}
	defer file.Close()

	if err := os.Rename(file.Name(), filepath); err != nil {
		return err
	}

	return s.writeContent(file, list)
}

func (s *FileSystem) writeContent(w io.Writer, list *data.List) error {
	buf := new(bytes.Buffer)

	// frontmatter
	// open

	buf.WriteString("---\n")
	buf.WriteString(fmt.Sprintf("title: %q\n", list.Name))
	buf.WriteString(fmt.Sprintf("description: %q\n", list.Description))
	buf.WriteString("date: " + list.CreatedAt.Format(time.RFC3339) + "\n")
	if len(list.Items) > 0 && len(list.Items[0].Images) > 0 {
		buf.WriteString("image: " + list.Items[0].Images[0] + "\n")
	}
	buf.WriteString("tags:\n  - dark\n")
	buf.WriteString("categories:\n  - books\n")
	buf.WriteString("draft: false\n")
	buf.WriteString("toc: true\n")

	// close
	buf.WriteString("---\n")

	// body
	for _, item := range list.Items {
		if item.Position > 0 {
			buf.WriteString(fmt.Sprintf("## %d. %s\n\n", item.Position, item.Title))
		} else {
			buf.WriteString(fmt.Sprintf("## %s\n\n", item.Title))
		}
		if len(item.By) > 0 {
			buf.WriteString("**By:** " + strings.Join(item.By, ", ") + "\n\n")
		}
		if item.AvgRating > 0 {
			buf.WriteString(fmt.Sprintf("⭐ %.2f (%d ratings)\n\n", item.AvgRating, item.TotalRatings))
		}
		buf.WriteString(fmt.Sprintf(
			"[View on Goodreads](https://www.goodreads.com/book/show/%s)\n\n",
			item.ID,
		))
		if len(item.Images) > 0 {
			buf.WriteString(fmt.Sprintf(`![%s](%s "{width='250',height='350'}")`, item.Title, item.Images[0]))
			buf.WriteString("\n\n")
		}

		// Divider
		buf.WriteString("---\n\n")
	}

	_, err := w.Write(buf.Bytes())
	return err
}
