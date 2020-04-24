package staff

// This trait allows credentials to
// become staff users. Such users
// can perform administrative actions.
type StaffCapable interface {
	IsStaff() bool
}
