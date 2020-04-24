package superuser

// This trait allows credentials to
// become superusers. Superusers can
// do ANYTHING.
type SuperuserCapable interface {
	IsSuperuser() bool
}
