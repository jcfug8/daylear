package namer

// newReflectNamerConfig holds configuration options for creating a ReflectNamer.
type newReflectNamerConfig struct {
	disableStrictNoMissingStructKeys bool     // If true, disables strict checking for missing struct keys
	extraPatterns                    []string // Additional resource name patterns to support
}

// newReflectNamerOption defines a function that modifies newReflectNamerConfig.
type newReflectNamerOption func(config newReflectNamerConfig) newReflectNamerConfig

// WithExtraPatterns adds extra resource name patterns to the ReflectNamer configuration.
func WithExtraPatterns(patterns []string) newReflectNamerOption {
	return func(config newReflectNamerConfig) newReflectNamerConfig {
		config.extraPatterns = patterns
		return config
	}
}

// DisableStrictNoMissingStructKeys sets whether strict checking for missing struct keys is disabled.
// If set to true, the ReflectNamer will not require all pattern keys to be present in the struct.
func DisableStrictNoMissingStructKeys(disable bool) newReflectNamerOption {
	return func(config newReflectNamerConfig) newReflectNamerConfig {
		config.disableStrictNoMissingStructKeys = disable
		return config
	}
}

// formatReflectNamerConfig holds configuration for formatting resource names.
type formatReflectNamerConfig struct {
	patternIndex int // The index of the pattern to use for formatting
}

// formatReflectNamerOption defines a function that modifies formatReflectNamerConfig.
type formatReflectNamerOption func(config formatReflectNamerConfig) formatReflectNamerConfig

// AsPatternIndex sets the pattern index to use when formatting a resource name.
func AsPatternIndex(patternIndex int) formatReflectNamerOption {
	return func(config formatReflectNamerConfig) formatReflectNamerConfig {
		config.patternIndex = patternIndex
		return config
	}
}
