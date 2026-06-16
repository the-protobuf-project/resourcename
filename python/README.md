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

@resourcename("//music.example.com/artists/{artist_id}")
class Artist:
    pass

# Parse a resource name
parsed = Artist.resourcename.parse("//music.example.com/artists/radiohead")
print(parsed)  # {'artist_id': 'radiohead'}

# Generate a resource name
name = Artist.resourcename.generate(artist_id="bjork")
print(name)    # "//music.example.com/artists/bjork"
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
