Requirements
============

In order to use respective hashing algorithms:

  - `bcrypt`: `golang.org/x/crypto/bcrypt` (any version will do).
  - `argon2`: `golang.org/x/crypto/argon2` (any version will do).

For the credentials implementations:

  - `gorm`: `github.com/jinzhu/gorm` version `v1.9.8`.

You have to install them externally.

Usage
-----

You need to correctly define classes implementing the following interfaces:

  - `stub.Credential`: They will be, essentially, the users.
  - `stub.Scope`: They will be, essentially, the requirements or permissions.
  - `stub.Source`: It will lookup credentials, essentially, by their identification.
  
And optional interfaces like:

  - `stub.WithSuperUserFlag`: Credentials also implementing interfaces can be
    considered as super-users if their method returns true.

Assuming you have those interfaces correctly implemented, you can invoke:

  - `Login(realm string, managers CredentialsMultiManager, identification interface{}, password string)`:
    Tries to log a credential in. It may fail due to password mismatch, empty password,
    or another log in restriction failure. If the log in operation is successful, it returns
    the logged in credential's realm, and the logged in credential.
  - `SetPassword(credential stub.Credential, password string)`: Sets a new password on the object.
    It does by hashing the password (according to the credential's hashing engine) and
    stores it by the same mean the Credential provides to store the password hash.
  - `ClearPassword(credential stub.Credential)`: Actually tells the credential to clear its password.
  - `Authorize(credential stub.Credential, requirement stub.AuthorizationRequirement)`:
    Checks whether a specific credential (this usually apples to logged ones) is authorized
    by that requirement (which could be a single or complex one).

For login to work, a `CredentialsMultiManager` must be created. It is just a kind of map with string keys,
and values of type `CredentialsManager`. You'll be using it indirectly through the `Login` function, and
middleware implementations/plugins will make use of it.

Configuring hashers
-------------------

Hashers (hashing engines) are the means of safely hashing, keeping and comparing passwords.

This package currently supports the following engines:

  - `bcrypt` (path: `implementations/hashing/bcrypt`).
  - `argon2` (path: `implementations/hashing/argon2`).
  - Create your own implementing `stub.PasswordHashingEngine`.

And also a mixed one. This mixed can take arbitrary hashing engines, a default one, and
read many different passwords (provided the appropriate hasher is registered in the mixed
one), and generate new passwords using the default one.

To get a default hasher quickly, you can use:

  - `bcrypt.Default`
  - `argon2.Default`
  
Or perhaps customized hashers by invoking:

  - `bcrypt.New(...)`
  - `argon2.New(...)`

If you understand how do those algorithms work.

For the mixed case (path: `implementations/hashing/multiple`), you use:

  - `multiple.New(hasher1, hasher2, ...)` - Takes the first hasher as default.
  - `multiple.NewWithDefault(hasher2, hasher1, hasher2, ...)` - Takes `hasher2` as default.

Considering that mixed hashers cannot be added to other mixed hashers, and that
you cannot add more than one hasher of the same type.

Keep the hasher of your choice preserved inside a global variable or something like that.
You'll use it later, when creating a custom Credential.

Default Implementation
----------------------

For `gorm` there is a default implementation at `implementations/persistence/gorm`.
Let's call that package `gorm_impl`, to not collision with the `gorm`'s package name.

It consists on three parts:
  - `gorm_impl.ModelBackedScope`: It is the type of a stored scope in database.
  - `gorm_impl.User`: It is an abstract class with some user functionalities.
    You still have to inherit (compose it) this class and define the
    `HashingEngine()` to make it return your newly created / obtained hasher.
  - `gorm_impl.Lookup(*gorm.DB)`: A gorm lookup mechanism that allows us to
    search for those credentials.

When you're ready, you can install them: `db.AutoMigrate(&ModelBackedScope{}, &MyUserSubclass{})`.
If you want, you can peek the source code of those two classes and implement types by yourself.

And that's it! Migrate them, run them, and have your system around it.