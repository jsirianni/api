package store

import (
	"fmt"
	"sync"

	"github.com/jsirianni/server/model"
)

func NewMemory() *Memory {
	return &Memory{
		data: make(map[string]account),
	}
}

// Memory is an in memory store.
type Memory struct {
	// accountID is used as the index
	data map[string]account

	// TODO(jsirianni): Consider using a RWMutex
	m sync.Mutex
}

type account struct {
	accountKey string

	// deviceID is used as the index
	devices map[string]model.Device
}

var _ Store = (*Memory)(nil)

// CheckSubscription returns an error if the given account
// is invalid.
func (m *Memory) CheckSubscription(accountID, accountKey string) error {
	m.m.Lock()
	defer m.m.Unlock()

	if err := m.validateAccount(accountID, accountKey); err != nil {
		return fmt.Errorf("subscription validation failed: %v", err)
	}
	return nil
}

// RegisterDevice takes an accountID, accountKey, deviceID, deviceInfo and stores
// the device if the account is valid.
func (m *Memory) RegisterDevice(accountID, accountKey, deviceID string, deviceInfo model.Device) error {
	m.m.Lock()
	defer m.m.Unlock()

	if err := m.validateAccount(accountID, accountKey); err != nil {
		return fmt.Errorf("subscription validation failed: %v", err)
	}

	// assume the account exists because validation passed
	// and this function holds a lock against the map.
	account := m.data[accountID]

	// Okay to overwrite existing deviceInfo if
	// the deviceID already exists
	account.devices[deviceID] = deviceInfo

	// Persist the modified account to the map
	m.data[accountID] = account

	return nil
}

// Accounts returns all accounts
func (m *Memory) Accounts() ([]string, error) {
	m.m.Lock()
	defer m.m.Unlock()

	accounts := []string{}
	for accountID, _ := range m.data {
		accounts = append(accounts, accountID)
	}

	return accounts, nil
}

// Devices returns all devices for a given account
func (m *Memory) Devices(accountID string) ([]model.Device, error) {
	m.m.Lock()
	defer m.m.Unlock()

	account, ok := m.data[accountID]
	if !ok {
		return nil, fmt.Errorf("account with id %s does not exist", accountID)
	}

	devices := []model.Device{}
	for _, device := range account.devices {
		devices = append(devices, device)
	}

	return devices, nil
}

func (m *Memory) validateAccount(accountID, accountKey string) error {
	notFound := fmt.Errorf("account does not exist or account key is invalid: %s", accountID)

	account, ok := m.data[accountID]
	if !ok {
		return notFound
	}

	if account.accountKey != accountKey {
		return notFound
	}

	return nil
}
