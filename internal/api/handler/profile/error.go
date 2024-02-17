package profile

const (
	AlreadyExists        string = "User profile already exists"
	DoesNotExist         string = "User profile does not exist"
	GamerTagAlreadyInUse string = "Gamer tag already in use"
)

// TODO: Move the below errors to the referral handler
const (
	AlreadyReferred      string = "User already referred"
	ReferrerDoesNotExist string = "Referrer does not exist"
	SelfReferralError    string = "A user cannot refer themselves"
)
