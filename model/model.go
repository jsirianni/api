package model

// Account represents a user account.
type Account struct {
	// The account id.
	ID string

	// The account's authentication key.
	Key string

	// Active represents whether or not the account
	// has an active subscription.
	Active bool
}

// Device represents an enduser device.
type Device struct {
	// AccountID is the account the device is assosiated with
	AccountID string

	// ID is the id of the device
	ID string

	// The hostname of the device
	Hostname string
}
