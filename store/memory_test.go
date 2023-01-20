package store

import (
	"testing"

	"github.com/jsirianni/server/model"
	"github.com/stretchr/testify/require"
)

func TestNewMemory(t *testing.T) {
	m := NewMemory()

	require.NotNil(t, m)
	require.NotNil(t, m.data)
}

func TestCheckSubscription(t *testing.T) {
	m := newSeededMemoryStore()

	require.NoError(t, m.CheckSubscription("abc", "xyz"))
	require.NoError(t, m.CheckSubscription("go", "095"))
	require.Error(t, m.CheckSubscription("bad", "sub"), "expected an error when an invalid account id was given")
	require.Error(t, m.CheckSubscription("abc", "invalid"), "expected an error when a valid account is given with the wrong key")
}

func TestRegisterDevice(t *testing.T) {
	m := newSeededMemoryStore()

	_, ok := m.data["abc"].devices["test-device"]
	require.False(t, ok, "expected test-device to not exist")

	err := m.RegisterDevice("abc", "xyz", "test-device", model.Device{})
	require.NoError(t, err)

	_, ok = m.data["abc"].devices["test-device"]
	require.True(t, ok, "expected test-device to exist")

	err = m.RegisterDevice("abc", "ttt", "test-device", model.Device{})
	require.Error(t, err, "expected an error when registering a device using an invalid account key")
}

func TestAccounts(t *testing.T) {
	m := newSeededMemoryStore()

	accounts, err := m.Accounts()
	require.NoError(t, err)
	require.NotNil(t, accounts)
	require.Len(t, accounts, 2, "expected exactly two accounts, 'abc' and 'go'")
}

func TestDevices(t *testing.T) {
	m := newSeededMemoryStore()

	devices, err := m.Devices("abc")
	require.NoError(t, err)
	require.NotNil(t, devices)
	require.Len(t, devices, 2, "expected exactly two devices for account 'abc'")

	_, err = m.Devices("badaccount")
	require.Error(t, err, "expected an error when looking up devices for an account that does not exist")
}

func newSeededMemoryStore() *Memory {
	seed := map[string]account{}
	seed["abc"] = account{
		accountKey: "xyz",
		devices: map[string]model.Device{
			"device-a": {},
			"device-b": {},
		},
	}
	seed["go"] = account{
		accountKey: "095",
		devices: map[string]model.Device{
			"apple": {},
			"pear":  {},
		},
	}
	m := NewMemory()
	m.data = seed
	return m
}
