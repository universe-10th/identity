package password

// TODO
// This file contains definitions and pipelines to change
// the password of any credential (since all credentials
// support hashing passwords and also support setting a
// custom hashed password).
//
// Two actions will be defined here: Set a password, and
// Unset the password (which means: set a blank password).
// Unsetting the password could be reverted by recovering
// the credential via the recovery pipeline.
// To make this action effective, a "source" must also be
// defined in this context, which involves saving back a
// credential.
