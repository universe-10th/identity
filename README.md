Requirements
------------

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

Assuming you have those interfaces correctly implemented, you can create what I called a `realm`
(which is barely more than a relationship between a `Credential` type and a `Source` engine).

Realms provide two methods of interest:

  - `Unmarshal(pk interface{}) (stub.Credential, error)`: Retrieves a credential of the realm's type,
    given a key (which is defined by the `Credential` as the primary one, but that's not necessarily
    true at database level: it may be another candidate key there).
  - `Lookup(identification interface{}) (stub.Credential, error)`: This is similar, but searches for
    the credential's identification field (and case-sensitivity settings) instead of strictly against
    the primary key (or non-human-friendly candidate key).

But you'll rarely use these methods: the nearest thing you will do is implementing the underlying
`Source` you'll use. The realm is used indirectly, however: you'll instantiate something I added as
`MultiRealm` (which is, just to start, a map with string keys, and `Realm` values). Once you instantiate
this multi-realm object, you'll have access to these 3 methods you may use directly:

  - `Unmarshal(realm string, pk interface{}) (stub.Credential, error)`: Given a realm key (which must
    exist in the map) it unmarshals a given key into the corresponding credential in the realm. This
    calls `Unmarshal` in the corresponding `Realm` object.
  - `Lookup(realm string, identification interface{}) (stub.Credential, error)`: This method is quite
    similar but, instead of looking by its key, the credentials lookup is done by identification. This
    calls `Lookup` in the corresponding `Realm` object.
  - `MultiLookup(identification interface{}) (string, stub.Credential, error)`: This method performs
    a lookup of the given key by testing each realm (order is not guaranteed!) until it finds a match.
    The matched realm and credential are returned, or an error if there is no match.
  - `Login(realm string, identification interface{}, password string) (string, stub.Credential, error)`:
    Performs a lookup for a credential given its identification and realm (if realm is `"""`, performs
    a multi-lookup instead), and tries to log-in given its password. If there is a matched credential,
    its password also matches, and it _can login_ (which is defined by the `Credential` itself) then
    the credential and the realm are returned. If no credential and/or password is matched, or it cannot
    login (according to the said criteria) then an error is returned, and no credential/realm.

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
  - `gorm_impl.NewSource(*gorm.DB)`: A gorm lookup mechanism that allows us to
    search for those credentials (satisfies the `Source` interface).

When you're ready, you can install them: `db.AutoMigrate(&ModelBackedScope{}, &MyUserSubclass{})`.
If you want, you can peek the source code of those two classes and implement types by yourself.

And that's it! Migrate them, run them, and have your system around it.

Web Helpers
-----------

For Iris framework, there is one adapter at `github.com/universe-10th/identity/plugins/iris`.

Such adapter works by instantiating it against a current `Session` (which is either a regular
Iris session or a JWT-based one when using module `github.com/universe-10th/iris-jwt-sessions`).
Say this one is the adapter:

    import (
        ...
        webRealms "github.com/universe-10th/identity/plugins/iris"
        ...
    )
    
    myMultiRealm := MultiRealm{... whatever realms you declare here ...}
    newWebRealm := webRealms.Factory(myMultiRealm)

    func MyHandler(session *sessions.Session) string, int {
        adapter := newWebRealm(session)
    }

You can call these methods:

  - `Login(realm string, identification interface{}, password string) (string, stub.Credential, error)`:
    Tries to perform a login against that credential (using the aforementioned `Login` method in multi-realms).
    On successful login, it stores the credential's key and realm in the involved session.
  - `Logout()`: Just clears any realm / key values from the involved session.
  - `Current() (string, stub.Credential)`: Gets the current credential. It takes the stored key and realm from
    the involved session to know how to unmarshal and return the credential. If any error occurs, or any data
    is missing, a silent logout will be performed. **Caution**: This method actually hits the database, as the
    `Login` method does. Avoid calling this method if you already called `Login`: just keep the result value of
    `Login` and you'll save extra unnecessary access to database(s) or storage(s).

Notes
-----

**TESTS TO BE ADDED!!!**