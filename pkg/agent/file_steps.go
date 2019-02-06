package agent

import (
	"crypto/rand"
	"os"

	"github.com/nicktitle/rc/pkg/fileformat"

	"github.com/pkg/errors"
)

func createFile(path string) (Event, error) {
	f, err := os.Create(path)
	if err != nil {
		return Event{}, errors.Wrapf(err, "failed to create file at path: %s", path)
	}
	defer f.Close()

	f.Sync()
	return fileEvent(path, fileformat.FileCreate), nil
}

func modifyFile(path string) (Event, error) {
	var e Event

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return e, errors.Wrapf(err, "failed to open file for modify at path: %s", path)
	}
	defer f.Close()

	modBytes := make([]byte, modByteLen)
	if _, err := rand.Read(modBytes); err != nil {
		return e, errors.Wrap(err, "failed to generate modbytes")
	}

	if _, err := f.Write(modBytes); err != nil {
		return e, errors.Wrap(err, "failed to write new bytes to file")
	}

	// sync is important, since for large files, the defer up top can fire
	// before the file is finished being written to
	f.Sync()

	return fileEvent(path, fileformat.FileModify), nil
}

func deleteFile(path string) (Event, error) {
	if err := os.Remove(path); err != nil {
		return Event{}, errors.Wrapf(err, "failed to delete file at path: %s", path)
	}
	return fileEvent(path, fileformat.FileDelete), nil
}

func fileEvent(path string, kind fileformat.StepKind) Event {
	return Event{
		"path":     path,
		"activity": kind,
	}
}
