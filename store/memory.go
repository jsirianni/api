package store

import (
	"fmt"
	"sync"

	"github.com/jsirianni/server/model"
)

// NewMemory returns a new memory store.
func NewMemory() *Memory {
	return &Memory{
		accounts: []model.Account{},
		devices:  make(map[string][]model.Device),
	}
}

// NewTestingMemory returns a new memory store with seeded data
// for development use.
func NewTestingMemory() *Memory {
	m := NewMemory()
	m.accounts = []model.Account{
		{
			ID:     "abc",
			Key:    "xyz",
			Active: true,
		},
		{
			ID:     "go",
			Key:    "095",
			Active: false,
		},
	}
	m.devices = map[string][]model.Device{
		"abc": {
			{
				ID:        "device-a",
				AccountID: "abc",
				Hostname:  "testname",
			},
			{
				ID:        "device-b",
				AccountID: "abc",
				Hostname:  "testname-b",
			},
		},
	}
	return m
}

// Memory is an in memory store.
type Memory struct {
	// accounts are stored as slices
	// because we do not care about performance for
	// testing.
	accounts []model.Account

	// device map is indexed with account id
	devices map[string][]model.Device

	mu sync.Mutex
}

var _ Store = (*Memory)(nil)

// CheckSubscription returns an error if the given account
// is invalid.
func (m *Memory) CheckSubscription(accountID, accountKey string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, err := m.validateAccount(accountID, accountKey); err != nil {
		return fmt.Errorf("subscription validation failed for account with id %s: %v", accountID, err)
	}
	return nil
}

// RegisterDevice takes an accountID, accountKey, deviceInfo and stores
// the device if the account is valid.
func (m *Memory) RegisterDevice(accountID, accountKey string, device model.Device) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, err := m.validateAccount(accountID, accountKey); err != nil {
		return fmt.Errorf("subscription validation failed: %v", err)
	}

	// If account not in device map, index it and add the device.
	if _, ok := m.devices[accountID]; !ok {
		m.devices[accountID] = []model.Device{device}
	}

	// If account in device map, check if device already exists
	for i, d := range m.devices[accountID] {
		// This should always be the case because devices are indexed
		// under accountID, but check to be sure.
		if d.AccountID == accountID {
			if d.ID == device.ID {
				// replace device in slice.
				m.devices[accountID][i] = device
				return nil
			}
		}
	}

	// If device did not exist already, append it
	m.devices[accountID] = append(m.devices[accountID], device)

	return nil
}

// Account returns an account
func (m *Memory) Account(accountID string) (model.Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, a := range m.accounts {
		if a.ID == accountID {
			return a, nil
		}
	}

	return model.Account{}, fmt.Errorf("account with id %s does not exist", accountID)
}

// Devices returns all devices for a given account
func (m *Memory) Devices(accountID string) ([]model.Device, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	accountDevices, ok := m.devices[accountID]
	if !ok {
		return nil, fmt.Errorf("account with id %s does not exist in device map", accountID)
	}

	devices := []model.Device{}
	for _, device := range accountDevices {
		if device.AccountID == accountID {
			devices = append(devices, device)
		}
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("account with id %s does not have any devices", accountID)
	}

	return devices, nil
}

// Device returns a device for a given account
func (m *Memory) Device(accountID, deviceID string) (model.Device, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	accountDevices, ok := m.devices[accountID]
	if !ok {
		return model.Device{}, fmt.Errorf("account with id %s does not exist in device map", accountID)
	}

	for _, device := range accountDevices {
		if device.AccountID == accountID {
			if device.ID == deviceID {
				return device, nil
			}
		}
	}

	return model.Device{}, fmt.Errorf("account with id %s does not have device with id %s", accountID, deviceID)
}

func (m *Memory) validateAccount(id, key string) (model.Account, error) {
	account := model.Account{}
	found := false

	for _, a := range m.accounts {
		if a.ID == id && a.Key == key {
			account = a
			found = true
			break
		}
	}

	if !found {
		return model.Account{}, fmt.Errorf("account does not exist or account key is invalid: %s", id)
	}

	return account, nil
}
