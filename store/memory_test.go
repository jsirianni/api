package store

import (
	"testing"

	"github.com/jsirianni/server/model"
	"github.com/stretchr/testify/require"
)

func TestNewMemory(t *testing.T) {
	m := NewMemory()

	require.NotNil(t, m)
	require.NotNil(t, m.accounts)
	require.NotNil(t, m.devices)
}

func TestCheckSubscription(t *testing.T) {
	m := NewTestingMemory()

	require.NoError(t, m.CheckSubscription("abc", "xyz"))
	require.NoError(t, m.CheckSubscription("go", "095"))
	require.Error(t, m.CheckSubscription("bad", "sub"), "expected an error when an invalid account id was given")
	require.Error(t, m.CheckSubscription("abc", "invalid"), "expected an error when a valid account is given with the wrong key")
}

func TestRegisterDevice(t *testing.T) {
	m := NewTestingMemory()

	// Add to device map
	d := []model.Device{
		{
			AccountID: "abc",
			ID:        "test-device",
			Hostname:  "orig",
		},
	}
	m.devices["abc"] = d

	err := m.RegisterDevice("abc", "xyz", model.Device{
		ID:        "test-device",
		AccountID: "abc",
		Hostname:  "new",
	})
	require.NoError(t, err)

	found := false
	for _, d := range m.devices["abc"] {
		if d.ID == "test-device" {
			found = true
			require.Equal(t, "new", d.Hostname, "expected RegisterDevice to update existing device")
		}
	}
	require.True(t, found, "expected found to be true, because we seeded the 'test-device'. This should never fail.")

	err = m.RegisterDevice("invalidaccount", "ttt", model.Device{})
	require.Error(t, err, "expected an error when registering a device using an invalid account key")
	require.ErrorContains(t, err, "subscription validation failed")
}

func TestAccount(t *testing.T) {
	m := NewTestingMemory()

	account, err := m.Account("abc")
	require.NoError(t, err)
	require.Equal(t, model.Account{
		ID:     "abc",
		Key:    "xyz",
		Active: true,
	}, account)

	account, err = m.Account("invalid")
	require.Error(t, err)
	require.Equal(t, model.Account{}, account)
}

func TestDevices(t *testing.T) {
	m := NewTestingMemory()

	devices, err := m.Devices("abc")
	require.NoError(t, err)
	require.NotNil(t, devices)
	require.Len(t, devices, 2, "expected exactly two devices for account 'abc'")

	_, err = m.Devices("badaccount")
	require.Error(t, err, "expected an error when looking up devices for an account that does not exist")
}

func TestDevice(t *testing.T) {
	m := NewTestingMemory()

	device, err := m.Device("abc", "device-a")
	require.NoError(t, err)
	// TODO(jsirianni): Add fields to device type
	require.Equal(t, model.Device{
		ID:        "device-a",
		AccountID: "abc",
		Hostname:  "testname",
	}, device)

	_, err = m.Device("badaccount", "")
	require.Error(t, err, "expected an error when looking up devices for an account that does not exist")
	require.ErrorContains(t, err, "account with id badaccount does not exist")

	_, err = m.Device("abc", "invalid")
	require.Error(t, err)
	require.ErrorContains(t, err, "account with id abc does not have device with id invalid")
}
