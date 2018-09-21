package controllers

import (
	"twreporter.org/go-api/storage"
	//log "github.com/Sirupsen/logrus"
)

// NewMembershipController ...
func NewMembershipController(s storage.MembershipStorage) *MembershipController {
	return &MembershipController{s}
}

// MembershipController ...
type MembershipController struct {
	Storage storage.MembershipStorage
}

// Close is the method of Controller interface
func (mc *MembershipController) Close() error {
	err := mc.Storage.Close()
	if err != nil {
		return err
	}
	return nil
}
