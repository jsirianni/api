package store

import "github.com/jsirianni/server/model"

type Store interface {
	// CheckSubscription returns an error if the given account
	// is invalid.
	CheckSubscription(accountID, accountKey string) error

	// RegisterDevice takes an accountID, accountKey, deviceID, deviceInfo and stores
	// the device if the account is valid.
	RegisterDevice(accountID, accountKey, deviceID string, deviceInfo model.Device) error

	// Accounts returns all accounts
	Accounts() ([]string, error)

	// Devices returns all devices for a given account
	Devices(accountID string) ([]model.Device, error)
}
