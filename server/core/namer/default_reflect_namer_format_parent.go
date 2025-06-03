package namer

import "go.einride.tech/aip/resourcename"

// MustFormatParent calls FormatParent and panics if an error occurs.
// It is a convenience method for callers who expect the parent resource to always be valid for formatting.
func (n *defaultReflectNamer[T]) MustFormatParent(in T, options ...formatReflectNamerOption) string {
	formatted, err := n.FormatParent(in, options...)
	if err != nil {
		panic(err)
	}
	return formatted
}

// FormatParent formats the parent portion of the resource name using the AIP parent pattern specified by the options.
// By default, it uses the parent pattern at index 0. If patternIndex is -1, it tries all parent patterns
// and returns the first one that matches. Returns an error if no parent pattern matches or if
// required parent fields are missing/invalid.
func (n *defaultReflectNamer[T]) FormatParent(in T, options ...formatReflectNamerOption) (string, error) {
	config := formatReflectNamerConfig{
		patternIndex: 0,
	}
	for _, option := range options {
		config = option(config)
	}

	// Build a list of pattern indices to try
	tryPatterns := []int{}
	if config.patternIndex == -1 {
		for i := range n.patternsDetails {
			tryPatterns = append(tryPatterns, i)
		}
	} else {
		tryPatterns = append(tryPatterns, config.patternIndex)
	}

	var lastErr error
	for _, idx := range tryPatterns {
		patternDetails := n.patternsDetails[idx]
		// If there is no parent pattern, treat as root (empty string)
		if patternDetails.parentPattern == "" {
			return "", nil
		}
		vars, err := n.extractPatternValues(patternDetails, in)
		if err == nil {
			return resourcename.Sprint(patternDetails.parentPattern, vars...), nil
		}
		lastErr = err
	}
	if lastErr == nil {
		lastErr = ErrNoPatternFound{
			msg: "no suitable parent pattern found for type",
		}
	}

	return "", lastErr
}
