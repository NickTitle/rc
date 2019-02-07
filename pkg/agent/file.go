package agent

import (
	"crypto/rand"
	"os"

	"github.com/nicktitle/rc/pkg/fileformat"

	"github.com/pkg/errors"
)

func createFile(settings fileformat.FileSettings) (Event, error) {
	f, err := os.Create(settings.Path)
	if err != nil {
		return Event{}, errors.Wrapf(err, "failed to create file at path: %s", settings.Path)
	}
	defer f.Close()

	f.Sync()
	return fileEvent(settings, fileformat.FileCreate), nil
}

func modifyFile(settings fileformat.FileSettings) (Event, error) {
	var e Event

	f, err := os.OpenFile(settings.Path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return e, errors.Wrapf(err, "failed to open file for modify at path: %s", settings.Path)
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

	return fileEvent(settings, fileformat.FileModify), nil
}

func deleteFile(settings fileformat.FileSettings) (Event, error) {
	if err := os.Remove(settings.Path); err != nil {
		return Event{}, errors.Wrapf(err, "failed to delete file at path: %s", settings.Path)
	}
	return fileEvent(settings, fileformat.FileDelete), nil
}

func fileEvent(settings fileformat.FileSettings, kind fileformat.StepKind) Event {
	return Event{
		"path":     settings.Path,
		"activity": kind,
	}
}
