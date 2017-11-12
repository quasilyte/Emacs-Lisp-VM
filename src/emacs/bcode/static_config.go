package bcode

const (
	// safetyCheck enables some additional checks and assetions
	// that are not expected to fail in valid code,
	// but may do so in development builds.
	//
	// This can lead to performance degradation.
	safetyCheck bool = true
)
