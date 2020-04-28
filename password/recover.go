package password

// TODO
// This file contains definitions & pipelines to change
// the password of credentials supporting the interface:
// credentials.traits.recoverable.Recoverable.
//
// Two actions will be defined here: Set a recovery code,
// attempt recovery (which involves: testing the recovery
// code against the input, and then set the new input
// password). Successful recoveries reset the formerly set
// password and also the case of a password being unset.
// To make this action effective, a "source" must also be
// defined in this context, which involves saving back a
// credential.
