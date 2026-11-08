package ui

// UI layout constants to avoid magic numbers throughout the codebase.
const (
	// Header section height (logo + version)
	HeaderHeight = 3

	// Footer section height (help text)
	FooterHeight = 3

	// Minimum height for info box displaying command details
	InfoBoxMinHeight = 3

	// Height of search input field
	SearchInputHeight = 2

	// Height reserved for pagination display
	PaginationHeight = 1

	// Default number of items per page in paginated views
	DefaultPageSize = 15

	// Minimum recommended terminal width
	MinWindowWidth = 80

	// Minimum recommended terminal height
	MinWindowHeight = 20

	// Maximum number of tags to display before truncating
	MaxVisibleTags = 5

	// Maximum length for truncated strings
	MaxTruncatedLength = 100
)
