"""Resource Name Template Module.

Provides bidirectional conversion between resource names and component values
using template patterns with placeholders.

Example:
    from google import resourcename

    @resourcename("//system.com/devices/{device_id}")
    class Device:
        pass

    # Parse a resource name
    parsed = Device.resourcename.parse("//system.com/devices/router-01")
    print(parsed)  # {'device_id': 'router-01'}

    # Generate a resource name
    name = Device.resourcename.generate(device_id="sensor-22")
    print(name)  # "//system.com/devices/sensor-22"
"""

import re
from collections.abc import Callable
from dataclasses import dataclass


def _extract_placeholders(template: str) -> list[str]:
    """Return a list of placeholder names found inside `{...}` within the template string."""
    return re.findall(r"\{([^{}]+)\}", template)


def _compile_regex(template: str, placeholders: list[str]) -> re.Pattern:
    """Convert a template with placeholders into a compiled regular-expression re.Pattern."""
    # Replace placeholders with temporary marker first
    pattern = template
    for ph in placeholders:
        pattern = pattern.replace("{" + ph + "}", "<<>>")

    # Escape special regex characters
    pattern = re.escape(pattern)

    # Replace markers with capture groups
    pattern = pattern.replace("<<>>", r"([^/]+)")

    return re.compile("^" + pattern + "$")


@dataclass(frozen=True)
class ResourceTemplate:
    """Represents an immutable resource-name template with placeholder support.

    Allows parsing a resource string into its component values and generating
    a resource string from provided values using the same template.
    """

    template: str

    def __post_init__(self):
        """Validate the template, extract placeholders, and build the internal matching regex."""
        if not self.template:
            msg = "Template cannot be empty"
            raise ValueError(msg)

        placeholders = _extract_placeholders(self.template)
        if not placeholders:
            msg = "Template must contain at least one placeholder"
            raise ValueError(msg)

        # Check for duplicate placeholders
        if len(placeholders) != len(set(placeholders)):
            duplicates = [p for p in set(placeholders) if placeholders.count(p) > 1]
            msg = f"Template contains duplicate placeholders: {duplicates}"
            raise ValueError(msg)

        regex = _compile_regex(self.template, placeholders)
        object.__setattr__(self, "placeholders", tuple(placeholders))
        object.__setattr__(self, "regex", regex)

    def parse(self, resource_name: str) -> dict[str, str]:
        """Parse a resource name and return a dictionary mapping placeholder names to extracted values.

        Args:
            resource_name: The resource name string to parse

        Returns:
            Dictionary mapping placeholder names to their extracted values

        Raises:
            ValueError: If the resource name doesn't match the template
        """
        match = self.regex.fullmatch(resource_name)
        if not match:
            msg = f"Resource name '{resource_name}' does not match template '{self.template}'"
            raise ValueError(msg)
        return {ph: match.group(i + 1) for i, ph in enumerate(self.placeholders)}

    def generate(self, **values) -> str:
        """Generate a resource name by substituting placeholder values into the template.

        Args:
            **values: Keyword arguments mapping placeholder names to values

        Returns:
            The generated resource name string

        Raises:
            ValueError: If required placeholders are missing or values contain invalid characters
        """
        # Check for missing placeholders
        missing = set(self.placeholders) - set(values.keys())
        if missing:
            msg = (
                f"Missing values for placeholders: {sorted(missing)}. "
                f"Required: {list(self.placeholders)}, Provided: {list(values.keys())}"
            )
            raise ValueError(msg)

        # Check for extra values
        extra = set(values.keys()) - set(self.placeholders)
        if extra:
            msg = f"Unexpected values provided: {sorted(extra)}. Expected only: {list(self.placeholders)}"
            raise ValueError(msg)

        # Validate values don't contain slashes
        invalid = {k: v for k, v in values.items() if "/" in v}
        if invalid:
            msg = f"Values contain invalid character '/': {invalid}"
            raise ValueError(msg)

        # Build result
        result = self.template
        for ph, value in values.items():
            result = result.replace("{" + ph + "}", value)
        return result

    def __repr__(self) -> str:
        """Return a readable string representation of the template object."""
        return f"ResourceTemplate('{self.template}')"


class ResourceNamespace:
    """Namespace object that provides access to parse/generate methods.

    and template metadata. Attached as Class.resourcename.
    """

    def __init__(self, template: ResourceTemplate):
        """Initialize the namespace with a resource template."""
        self._template = template

    def parse(self, resource_name: str) -> dict[str, str]:
        """Parse a resource name into component values."""
        return self._template.parse(resource_name)

    def generate(self, **values) -> str:
        """Generate a resource name from component values."""
        return self._template.generate(**values)

    @property
    def template(self) -> str:
        """Get the template string."""
        return self._template.template

    @property
    def placeholders(self) -> tuple[str, ...]:
        """Get the tuple of placeholder names."""
        return self._template.placeholders

    def __repr__(self) -> str:
        """Return a readable string representation of the namespace object."""
        return f"ResourceNamespace(template='{self._template.template}')"


def resourcename(template: str) -> Callable[[type], type]:
    """Decorator that attaches a .resourcename namespace to a class.

    Usage:
        @resourcename("//system.com/devices/{device_id}")
        class Device:
            pass

        # Access via .resourcename namespace
        parsed = Device.resourcename.parse("//system.com/devices/router-01")
        name = Device.resourcename.generate(device_id="router-01")

        # Access metadata
        print(Device.resourcename.template)
        print(Device.resourcename.placeholders)

    Args:
        template: The resource name template with {placeholder} markers

    Returns:
        A decorator function that adds a .resourcename namespace to the class
    """
    tmpl = ResourceTemplate(template)
    namespace = ResourceNamespace(tmpl)

    def decorator(cls: type) -> type:
        cls.resourcename = namespace
        return cls

    return decorator
