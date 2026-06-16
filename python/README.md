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
import resourcename

# Class-based API
t = resourcename.ResourceTemplate("//music.example.com/artists/{artist_id}")
print(t.parse("//music.example.com/artists/radiohead"))  # {'artist_id': 'radiohead'}
print(t.generate(artist_id="bjork"))                     # "//music.example.com/artists/bjork"

# Decorator API
@resourcename.resource("//music.example.com/artists/{artist_id}")
class Artist:
    pass

print(Artist.resourcename.parse("//music.example.com/artists/radiohead"))  # {'artist_id': 'radiohead'}
print(Artist.resourcename.generate(artist_id="bjork"))                     # "//music.example.com/artists/bjork"
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
