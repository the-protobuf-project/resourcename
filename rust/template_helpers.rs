use regex::Regex;
use std::collections::{BTreeMap, BTreeSet};

use crate::ResourceNameError;

pub fn reject_duplicates(placeholders: &[String]) -> Result<(), ResourceNameError> {
    let mut seen = BTreeSet::new();
    let mut duplicates = BTreeSet::new();
    for placeholder in placeholders {
        if !seen.insert(placeholder.as_str()) {
            duplicates.insert(placeholder.clone());
        }
    }
    if duplicates.is_empty() {
        return Ok(());
    }
    Err(ResourceNameError::DuplicatePlaceholders(
        duplicates.into_iter().collect::<Vec<_>>().join(", "),
    ))
}

pub fn extract_placeholders(template: &str) -> Vec<String> {
    let mut placeholders = Vec::new();
    let mut start = 0usize;
    while let Some(open_rel) = template[start..].find('{') {
        let open = start + open_rel;
        let Some(close_rel) = template[open + 1..].find('}') else {
            break;
        };
        let close = open + 1 + close_rel;
        let name = &template[open + 1..close];
        if !name.is_empty() {
            placeholders.push(name.to_string());
        }
        start = close + 1;
    }
    placeholders
}

pub fn ensure_balanced_braces(template: &str) -> Result<(), ResourceNameError> {
    let mut depth = 0i32;
    for ch in template.chars() {
        if ch == '{' {
            depth += 1;
        } else if ch == '}' {
            depth -= 1;
            if depth < 0 {
                return Err(ResourceNameError::UnbalancedBraces);
            }
        }
    }
    if depth == 0 {
        return Ok(());
    }
    Err(ResourceNameError::UnbalancedBraces)
}

pub fn validate_placeholder_names(placeholders: &[String]) -> Result<(), ResourceNameError> {
    let valid = Regex::new(r"^[A-Za-z_][A-Za-z0-9_]*$").expect("valid placeholder regex");
    for name in placeholders {
        if !valid.is_match(name) {
            return Err(ResourceNameError::InvalidPlaceholder(name.clone()));
        }
    }
    Ok(())
}

pub fn compile_regex(template: &str, placeholders: &[String]) -> Result<Regex, ResourceNameError> {
    let mut pattern = template.to_string();
    for placeholder in placeholders {
        pattern = pattern.replace(&format!("{{{placeholder}}}"), "<<>>");
    }
    let escaped = regex::escape(&pattern);
    let final_pattern = escaped.replace("<<>>", "([^/]+)");
    Regex::new(&format!("^{final_pattern}$")).map_err(|e| ResourceNameError::Serde(e.to_string()))
}

pub fn validate_input(
    placeholders: &[String],
    values: &BTreeMap<String, String>,
) -> Result<(), ResourceNameError> {
    let expected: BTreeSet<&str> = placeholders.iter().map(String::as_str).collect();
    let provided: BTreeSet<&str> = values.keys().map(String::as_str).collect();
    let missing = expected.difference(&provided).copied().collect::<Vec<_>>();
    if !missing.is_empty() {
        return Err(ResourceNameError::MissingPlaceholders {
            missing: missing.join(", "),
            required: placeholders.join(", "),
            provided: values.keys().cloned().collect::<Vec<_>>().join(", "),
        });
    }
    let extra = provided.difference(&expected).copied().collect::<Vec<_>>();
    if !extra.is_empty() {
        return Err(ResourceNameError::UnexpectedValues {
            unexpected: extra.join(", "),
            expected: placeholders.join(", "),
        });
    }
    let bad = values
        .iter()
        .filter_map(|(k, v)| v.contains('/').then_some((k, v)))
        .map(|(k, v)| format!("{k}={v}"))
        .collect::<Vec<_>>();
    if bad.is_empty() {
        return Ok(());
    }
    Err(ResourceNameError::SlashInValue(bad.join(", ")))
}
