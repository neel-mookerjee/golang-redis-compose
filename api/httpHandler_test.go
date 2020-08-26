package main

import (
	"testing"
)

type HandlerWrapperTestSuite struct {
	handler HandlerWrapper
}

func NewHandlerWrapperTestSuite() *HandlerWrapperTestSuite {
	return &HandlerWrapperTestSuite{handler: HandlerWrapper{}}
}

// test uniqueness of the endpoint ids
func TestHandlerWrapper_GenerateUniqueId(t *testing.T) {
	ts := NewHandlerWrapperTestSuite()
	uId1, err1 := ts.handler.GenerateUniqueId()
	if err1 != nil {
		t.Errorf("Expected successful execution. Received error: %v", err1)
	}
	uId2, err2 := ts.handler.GenerateUniqueId()
	if err2 != nil {
		t.Errorf("Expected successful execution. Received error: %v", err2)
	}
	if uId1 == uId2 {
		t.Errorf("Expected uId1 and uId2 to be different. Found same: uId1=%s, uId2=%s", uId1, uId2)
	}
}
