package router

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/RussellLuo/validating/v2"
	"github.com/rlaskowski/go-email/model"
)

type MutlipartController struct {
	Reader *multipart.Reader
}

func (m *MutlipartController) Message() (*model.Message, error) {
	messageForm, err := m.walk("message")

	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(messageForm)

	if err != nil {
		return nil, err
	}
	message, err := m.unmarshalMessage(b)

	if err != nil {
		return nil, err
	}

	return message, nil
}

func (m *MutlipartController) File() (*model.File, error) {
	fileForm, err := m.walk("file")

	if err != nil {
		return nil, err
	}

	err = m.validateFileForm(fileForm)
	if err != nil {
		return nil, err
	}

	file := model.NewFile(fileForm.FileName())

	file.Reader = fileForm

	return file, nil
}

func (m *MutlipartController) walk(name string) (*multipart.Part, error) {
	for {
		part, err := m.Reader.NextRawPart()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}

		if strings.Contains(name, part.FormName()) {
			return part, nil
		}
	}

	return nil, fmt.Errorf("Could not find %s form data", name)
}

func (m *MutlipartController) unmarshalMessage(data []byte) (*model.Message, error) {
	message := new(model.Message)

	err := json.Unmarshal(data, message)
	if err != nil {
		return nil, err
	}

	validateErr := validating.Validate(validating.Schema{
		validating.F("Sender", &message.Sender):       validating.Nonzero(),
		validating.F("Recipient", &message.Recipient): validating.Nonzero(),
		validating.F("Subject", &message.Subject):     validating.Nonzero(),
		validating.F("Content", &message.Content):     validating.Nonzero(),
	})

	if validateErr != nil {
		return nil, validateErr
	}
	return message, nil
}

func (m *MutlipartController) validateFileForm(fileForm *multipart.Part) error {
	if !(len(fileForm.FileName()) > 0) {
		return fmt.Errorf("File name not found")
	}

	name := strings.TrimSuffix(fileForm.FileName(), filepath.Ext(fileForm.FileName()))
	if len(name) > 50 {
		return fmt.Errorf("File name include more than 50 characters")
	}

	return nil
}
