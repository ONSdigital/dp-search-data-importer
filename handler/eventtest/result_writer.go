// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package eventtest

import (
	"github.com/ONSdigital/dp-search-data-importer/handler"
)

var ()

// Ensure, that ResultWriterMock does implement ResultWriter.
// If this is not the case, regenerate this file with moq.
var _ handler.ResultWriter = &ResultWriterMock{}

// ResultWriterMock is a mock implementation of handler.ResultWriter.
//
//     func TestSomethingThatUsesResultWriter(t *testing.T) {
//
//         // make and configure a mocked handler.ResultWriter
//         mockedResultWriter := &ResultWriterMock{
//         }
//
//         // use mockedResultWriter in code that requires handler.ResultWriter
//         // and then make assertions.
//
//     }
type ResultWriterMock struct {
	// calls tracks calls to the methods.
	calls struct {
	}
}
