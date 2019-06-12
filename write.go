package gomessagestore

import (
	"context"

	//	"github.com/blackhatbrigade/gomessagestore/repository"
	"github.com/sirupsen/logrus"
)

type writer struct {
	atPosition *int64
}

//WriteOption provide optional arguments to the Write function
type WriteOption func(w *writer)

func checkWriteOptions(opts ...WriteOption) *writer {
	w := &writer{}
	for _, option := range opts {
		option(w)
	}
	return w
}

//Write Writes a Message to the message store.
func (ms *msgStore) Write(ctx context.Context, message Message, opts ...WriteOption) error {
	envelope, err := Message.ToEnvelope(message)
	if err != nil {
		logrus.WithError(err).Error("Write: Validation Error")

		return err
	}

	writeOptions := checkWriteOptions(opts...)
	if writeOptions.atPosition != nil {
		err = ms.repo.WriteMessageWithExpectedPosition(ctx, envelope, *writeOptions.atPosition)
	} else {
		err = ms.repo.WriteMessage(ctx, envelope)
	}
	if err != nil {
		logrus.WithError(err).Error("Write: Error writing message")

		return err
	}
	return nil
}

//AtPosition allows for writing messages using an expected position
func AtPosition(position int64) WriteOption {
	return func(w *writer) {
		w.atPosition = &position
	}
}