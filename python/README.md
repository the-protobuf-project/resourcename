# The Protobuf Project Google Resource Name Parser

A Python library for parsing and manipulating Google Cloud Resource Names (GRNs) according to the [Resource Name Proposal](https://aip.dev/122).

## Features

- Define resource name templates using class decorators
- Parse resource names into their components
- Generate resource names from components
- Type-safe resource name handling

## Installation

### Using pip

```bash
pip install -e .
```

## Usage

```python
from resourcename import resourcename

@resourcename("//system.com/devices/{device_id}")
class Device:
    pass

# Parse a resource name
parsed = Device.resourcename.parse("//system.com/devices/router-01")
print(parsed)  # {'device_id': 'router-01'}

# Generate a resource name
name = Device.resourcename.generate(device_id="sensor-22")
print(name)    # "//system.com/devices/sensor-22"
```

## Development

### Setup

1. Clone the repository
2. Create a virtual environment:
   ```bash
   python -m venv venv
   source venv/bin/activate  # On Windows use `venv\Scripts\activate`
   ```
3. Install development dependencies:
   ```bash
   pip install -e ".[dev]"
   ```

### Running Tests

```bash
pytest
```
