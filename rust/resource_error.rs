//! Error types for resource-name parsing and generation.

/// Errors produced while validating templates, parsing resource names, or generating values.
#[derive(Debug, thiserror::Error)]
pub enum ResourceNameError {
    /// Template was empty.
    #[error("Template cannot be empty")]
    EmptyTemplate,
    /// Template had no `{placeholder}` entries.
    #[error("Template must contain at least one placeholder")]
    NoPlaceholders,
    /// Template used the same placeholder more than once.
    #[error("Template contains duplicate placeholders: {0}")]
    DuplicatePlaceholders(String),
    /// Input resource string does not match the compiled template.
    #[error("Resource name '{name}' does not match template '{template}'")]
    TemplateMismatch { name: String, template: String },
    /// Missing expected keys while generating a resource string.
    #[error(
        "Missing values for placeholders: {missing}. Required: [{required}], provided: [{provided}]"
    )]
    MissingPlaceholders {
        /// Missing key names.
        missing: String,
        /// Ordered placeholder list from the template.
        required: String,
        /// Keys present in the input.
        provided: String,
    },
    /// Extra keys were provided that are not in template placeholders.
    #[error("Unexpected values: {unexpected}. Expected only: [{expected}]")]
    UnexpectedValues {
        /// Comma-separated unexpected key names.
        unexpected: String,
        /// Comma-separated expected key names.
        expected: String,
    },
    /// One or more values contained `'/'`, which is invalid for one path segment.
    #[error("Values must not contain '/': {0}")]
    SlashInValue(String),
    /// Typed generation requires a JSON-like object map.
    #[error("Expected key/value object for generation")]
    GenerateInputMustBeObject,
    /// Typed generation accepts only scalar values for placeholders.
    #[error("Non-scalar value for placeholder '{field}'")]
    NonStringValue { field: String },
    /// Template has invalid brace structure (`{` and `}` are not balanced).
    #[error("Template has unbalanced braces")]
    UnbalancedBraces,
    /// Placeholder name contains unsupported characters.
    #[error("Invalid placeholder name '{0}'")]
    InvalidPlaceholder(String),
    /// Wrapped serde error text.
    #[error("Serde error: {0}")]
    Serde(String),
}

impl From<serde_json::Error> for ResourceNameError {
    fn from(value: serde_json::Error) -> Self {
        Self::Serde(value.to_string())
    }
}
