package invalid

type API interface {
	// ValidMethod is a fake method for testing the
	// plugin comment checker with a valid comment.
	//
	// Minimum server version: 1.2.3
	ValidMethod()

	// InvalidMethod is a fake method for testing the
	// plugin comment checker with an invalid comment.
	InvalidMethod()
}
