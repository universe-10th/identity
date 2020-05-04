Requirements
------------

This module has no requirements.

Usage
-----

**Credentials and Sources**

You need to correctly define classes implementing the following interfaces:

  - `credentials.Credential`: They will be, essentially, the users. There are complementary credentials that will also
    be useful depending on the context and needs of the system, like:
    - `credentials/traits/superuser.SuperuserCapable`: Such users _may_ become superusers in the appropriate context.
    - `credentials/traits/staff.StaffCapable`: Such users _may_ become staff users (while not being superusers).
    - `credentials/traits/scoped.Scoped`: Such users _may_ have scopes (self-identified permissions). Scopes are objects
      satisfying the `credentials/traits/scoped.Scope` interface, but the inner match will only make use of the `Key` in
      the scopes, not their references.
    - `credentials/traits/recoverable.Recoverable`: Such users may be password-reset by their owner when their password
      is lost.
    - `credentials/traits/indexed.Indexed`: Such users know their index (inner key) the sources use to retrieve them.
    - `credentials/traits/identified.Identified`: Such users know their identification the sources use to log them in.
    - `credentials/traits/deniable.Activable`: Such users know whether they must be considered active or inactive. They
      also have a mean to set such state.
    - `credentials/traits/deniable.Punishable`: Such users know whether they must be considered banned/restricted. They
      also have a mean to set such state.
  - `credentials.Broker`: They are means to get the credentials from an underlying store. This interface will seldom
    implemented, for there will exist common implementations (e.g. gorm, json, ...). **Notes**: when implementing your
    own broker, remember to return `nil, nil` in `ByIdentifier` if a credential was not found by its identifier.
    
Once these two interfaces (and the desired complementary ones) are implemented, a `credentials.Source` object must be
created via `credentials.NewSource(aBrokerInstance, YourUserType{})` (you can use any primitive-derived or struct type
as a `credentials.Credential` provided it is implemented correctly).

**Login pipeline**

Login process is implemented as a pipeline. After the credential is successfully retrieved it traverses a non-empty
pipeline which satisfies the `realm/login.PipelineStep` interface, which will tell if the login must be considered as
failed (by returning an error). Default implementations already exist for these pipeline steps like:

  - `realm/login/activity.ActivityStep`: Fails the login if the underlying credential implements the
    `credentials/traits/deniable.Activable` interface and the `Active` method returns false. Returns
    `realm.ErrLoginFailed` in that case.
  - `realm/login/password/PasswordCheckingStep` performs the actual password check. On failure, it will also return the
    `realm.ErrLoginFailed` error.
  - `realm/login/punish.PunishmentCheckStep` performs an "is punished" check. On failure, it will return a custom error
    of type `realm/login/punish.PunishmentCheckStep`. This only applies to credentials satisfying the punishable
    interface (`credentials/traits/deniable.Punishable`), while non-implementors will always pass.

When the pipeline is specified (as a variadic `...realm/login.PipelineStep` argument), the `ActivityStep` and the
`PasswordCheckingStep` must always run first (in the order you prefer, but before any other pipeline step). Most likely,
you will always want the `PasswordCheckingStep` interface in your pipeline, but for external logins it may be a
different case.

**Realms**

Realms are created by calling `realm.NewRealm(a source instance, ...pipeline step instances)`. They have methods like:

  - `user, err := Login(identifier, password)`: Attempts a login. Returns `realm.ErrLoginFailed` if no credential was
    found by the given identifier, or whatever the underlying source or pipeline step(s) return as an error.
  - `err := SetPassword(credential, password)`: Attempts a password change. The credential is then saved via the
    underlying source. Returns whatever the source returns on save, or the credential's hasher returns on hashing.
  - `err := UnsetPassword(credential)`: Attempts a password clear on a credential. Password-cleared credentials will
    always fail to login. Returns whatever the source returns on save, since the credential will also be saved in this
    case.
  - `err := ChangePassword(credential, current, new)`: Attempts a user-commanded password change. Aside from all the
    possible outcomes of `SetPassword`, it will also fail returning `realm.ErrBadCurrentPassword` if the current
    password is invalid.
  - `err := PreparePasswordReset(credential, token, duration)`: Sets a recovery token on a credential, and attempts a
    save of it. It will fail with `realm.ErrNotRecoverable` if the credential does not implement the 
    `credentials/traits/recoverable.Recoverable` interface, and will also return whatever the underlying source returns
    when attempting to save the credential.
  - `err := CancelPasswordReset(credential)`: Clears a recovery token on a credential. It returns whatever
    `PreparePasswordReset` would return for that credential and a token.
  - `err := ConfirmPasswordReset(credential, token, newPassword)`: Confirms a recovery process (password reset) on a
    credential. If the credential does not implement the `Recoverable` interface, it will return the
    `realm.ErrNotRecoverable` error. If the credential did not prepare any token, or the token is expired (considering
    the `duration` a parameter in `PreparePasswordReset` always sets a deadline for the token starting at the issue
    time) then `realm.ErrBadToken` will be returned. Otherwise, the same error results in the `SetPassword` may be
    returned.

**Authorization requirements**

Any object satisfying the `authreqs.AuthorizationRequirement` may be used to check if a credentials satisfies it, like:

    isSatisfied := myReq.SatisfiedBy(myCredential)

Several implementations are provided out of the box in this module:

  - `authreqs/superuser.RequireSuperuser`: Satisfies when the credential implements the
    `credentials/traits/superuser.SuperuserCapable`, and by invoking `Superuser()` on it, it returns true.
  - `authreqs/staff.RequireStaff`: Satisfies when the credential implements the
    `credentials/traits/staff.StaffCapable`, and by invoking `Staff()` on it, it returns true.
  - `authreqs/scoped.RequireScopesAmong(...specs)`: Satisfies when the credential implements the
    `credentials/traits/scoped.Scoped` interface and the scopes on it satisfy at least one of the scope specs among the
    arguments. The scope specs can be recursive structures involving:
    - An instance (implementor) of `credentials/traits/scoped.Scope`, which will satisfy when the credential has a scope
      with its same key.
    - The result of `authreqs/scoped.Any(...specs)`, which will recursively satisfy when at least one of the specs
      satisfies for the credential, and does not satisfy otherwise.
    - The result of `authreqs/scoped.All(...specs)`, which will recursively not satisfy when at least one of the specs
      does not satisfy for the credential, and satisfies otherwise.
  - `authreqs/compound.Admin(...specs)`: Applies the following rule:
  
        `RequireSuperuser.SatisfiedBy(c) || RequireStaff.SatisfiedBy(c) && RequireScopesAmong(...specs).SatisfiedBy(c)` 

     Intended for system administrators (superuser(s), limited-scope staff members).
  - `authreqs/compouns.TryAll{requirement1, requirement2, ...}`: Combines all the given authorization requirements in
    one, by satisfying when [at least] one of them satisfies for the current credential.

Custom hashers
--------------

Such like credentials and brokers, custom hashers may be created by implementing the `hashing.HashingEngine` interface.
They will have a `Name()` which should not collide with other implementations and may be per-instance.

A convenience multi-hasher is provided by creating one with `hashing.NewMultipleHashingEngine(...hashers)` or with a
_default_ engine with: `hashing.NewMultipleHashingEngineWithDefault(defaultEngine, ...engines)`. For this latter case,
the default engine must exist among the given `...engines`. For all the cases, the engines must be never nil, they
cannot be instances of `hashing.MultipleHashingEngine`, and at least one must be specified. For the first, the first
engine among the arguments will be used as the _default_ one.

This engine will attempt to check hashes like "foo:bar3435FSEF#" by using a registered hashing engine with name "foo"
to check a hash "bar3435FSEF#". Hashes with no "foo:" part will be attempted to check by the default engine. Finally,
hashing a password will always involve the default hashing engine.

This hasher (hashing engine) is intended  to have several changing hashing engines being used.