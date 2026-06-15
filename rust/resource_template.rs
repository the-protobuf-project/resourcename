//! Core template compiler and parse/generate implementation.

use regex::Regex;
use serde::Serialize;
use serde::de::DeserializeOwned;
use serde_json::{Map, Value};
use std::collections::BTreeMap;

use crate::ResourceNameError;
use crate::template_helpers::{
    compile_regex, ensure_balanced_braces, extract_placeholders, reject_duplicates, validate_input,
    validate_placeholder_names,
};

/// Immutable compiled template with parse/generate helpers.
#[derive(Debug, Clone)]
pub struct ResourceTemplate {
    template: String,
    placeholders: Vec<String>,
    regex: Regex,
}

impl ResourceTemplate {
    /// Compiles a template string.
    ///
    /// # Errors
    ///
    /// Returns [`ResourceNameError`] if the template is invalid.
    pub fn new(template: impl Into<String>) -> Result<Self, ResourceNameError> {
        let template = template.into();
        if template.is_empty() {
            return Err(ResourceNameError::EmptyTemplate);
        }
        ensure_balanced_braces(&template)?;
        let placeholders = extract_placeholders(&template);
        if placeholders.is_empty() {
            return Err(ResourceNameError::NoPlaceholders);
        }
        validate_placeholder_names(&placeholders)?;
        reject_duplicates(&placeholders)?;
        let regex = compile_regex(&template, &placeholders)?;
        Ok(Self {
            template,
            placeholders,
            regex,
        })
    }

    /// Returns original template source.
    #[must_use]
    pub fn template(&self) -> &str {
        &self.template
    }

    /// Returns placeholder names in declaration order.
    #[must_use]
    pub fn placeholders(&self) -> &[String] {
        &self.placeholders
    }

    /// Parses into `placeholder -> value`.
    ///
    /// # Errors
    ///
    /// Returns [`ResourceNameError::TemplateMismatch`] when no match.
    pub fn parse(
        &self,
        resource_name: &str,
    ) -> Result<BTreeMap<String, String>, ResourceNameError> {
        let captures = self.regex.captures(resource_name).ok_or_else(|| {
            ResourceNameError::TemplateMismatch {
                name: resource_name.to_string(),
                template: self.template.clone(),
            }
        })?;
        let mut values = BTreeMap::new();
        for (idx, placeholder) in self.placeholders.iter().enumerate() {
            let value = captures
                .get(idx + 1)
                .map_or_else(String::new, |m| m.as_str().to_string());
            values.insert(placeholder.clone(), value);
        }
        Ok(values)
    }

    /// Parses and deserializes into a typed struct.
    ///
    /// # Errors
    ///
    /// Returns [`ResourceNameError`] on parse or decode failure.
    pub fn parse_into<T>(&self, resource_name: &str) -> Result<T, ResourceNameError>
    where
        T: DeserializeOwned,
    {
        let values = self.parse(resource_name)?;
        let mut object = Map::new();
        for (k, v) in values {
            object.insert(k, Value::String(v));
        }
        Ok(serde_json::from_value(Value::Object(object))?)
    }

    /// Generates from a key/value map.
    ///
    /// # Errors
    ///
    /// Returns [`ResourceNameError`] for missing/extra/invalid keys.
    pub fn generate(&self, values: &BTreeMap<String, String>) -> Result<String, ResourceNameError> {
        validate_input(self.placeholders.as_slice(), values)?;
        let mut out = self.template.clone();
        for placeholder in &self.placeholders {
            let value =
                values
                    .get(placeholder)
                    .ok_or_else(|| ResourceNameError::MissingPlaceholders {
                        missing: placeholder.clone(),
                        required: self.placeholders.join(", "),
                        provided: values.keys().cloned().collect::<Vec<_>>().join(", "),
                    })?;
            out = out.replace(&format!("{{{placeholder}}}"), value);
        }
        Ok(out)
    }

    /// Generates from a serde-serializable object.
    ///
    /// # Errors
    ///
    /// Returns [`ResourceNameError`] when serialization is incompatible.
    pub fn generate_from<T>(&self, values: &T) -> Result<String, ResourceNameError>
    where
        T: Serialize,
    {
        let value = serde_json::to_value(values)?;
        let object = value
            .as_object()
            .ok_or(ResourceNameError::GenerateInputMustBeObject)?;
        let mut map = BTreeMap::new();
        for (k, v) in object {
            let string_value = match v {
                Value::String(s) => s.clone(),
                Value::Bool(b) => b.to_string(),
                Value::Number(n) => n.to_string(),
                _ => return Err(ResourceNameError::NonStringValue { field: k.clone() }),
            };
            map.insert(k.clone(), string_value);
        }
        self.generate(&map)
    }
}
