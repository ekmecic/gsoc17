package golibbeaglebone

// ExportState is state in which the device is in, either exported or unexported.
type ExportState int

const (
	// Exported means the device will be made ready to use.
	Exported ExportState = iota
	// UnExported means the device will be unavailable for use.
	// This is the default state when the device is first created.
	UnExported
)
