package password

// TODO
// This file defines the interface used to save
// a just-modified credential. This is actually
// the counterpart of the interface used to load
// a user for login, but this one defines more
// methods:
// - Load a credential by identifier only, since
//   it is assumed an admin is performing the
//   action, the credential is the logged user
//   requiring the action, or a nonce code was
//   appropriately validated.
// - Save a credential object (which must exist
//   beforehand, since we're not creating it).
