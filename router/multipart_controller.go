package router

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/RussellLuo/validating"
	"github.com/rlaskowski/go-email/model"
)

type mutlipartController struct {
	mutlipartReader *multipart.Reader
}

func (m *mutlipartController) Message() (*model.Message, error) {
	messageForm := m.walk("message")
	var messageModel *model.Message

	if messageForm == nil {
		return nil, fmt.Errorf("No message parameter")
	}

	messageModel = new(model.Message)
	message, err := ioutil.ReadAll(messageForm)

	if err != nil {
		return nil, err
	}
	err = m.unmarshalMessage([]byte(message), messageModel)

	if err != nil {
		return nil, err
	}
	return messageModel, nil
}

func (m *mutlipartController) File() (*multipart.Part, error) {
	fileForm := m.walk("file")

	if fileForm == nil {
		return nil, fmt.Errorf("No file parameter")
	}

	err := m.validateFileForm(fileForm)
	if err != nil {
		return nil, err
	}
	return fileForm, nil
}

func (m *mutlipartController) walk(name string) *multipart.Part {
	for {
		part, err := m.mutlipartReader.NextRawPart()
		if err == io.EOF || err != nil {
			break
		}

		if strings.Contains(part.FormName(), name) {
			return part
		}
	}
	return nil
}

func (m *mutlipartController) unmarshalMessage(data []byte, message *model.Message) error {
	err := json.Unmarshal(data, message)
	if err != nil {
		return err
	}

	validateErr := validating.Validate(validating.Schema{
		validating.F("Sender", &message.Sender):       validating.Nonzero(),
		validating.F("Recipient", &message.Recipient): validating.Nonzero(),
		validating.F("Subject", &message.Subject):     validating.Nonzero(),
		validating.F("Content", &message.Content):     validating.Nonzero(),
	})

	if validateErr != nil {
		return validateErr
	}
	return nil
}

func (m *mutlipartController) validateFileForm(fileForm *multipart.Part) error {
	if !(len(fileForm.FileName()) > 0) {
		return fmt.Errorf("File name not found")
	}

	name := strings.TrimSuffix(fileForm.FileName(), filepath.Ext(fileForm.FileName()))
	if len(name) > 50 {
		return fmt.Errorf("File name include more than 50 characters")
	}

	return nil
}
