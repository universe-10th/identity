package types


/**
 * Login check stages (to add custom checks).
 */
type LoginStage int
const (
	BeforePasswordCheck = LoginStage(0)
	AfterPasswordCheck = LoginStage(1)
)
