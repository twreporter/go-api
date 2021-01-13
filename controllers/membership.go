package controllers

import (
	"github.com/twreporter/go-api/storage"
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
