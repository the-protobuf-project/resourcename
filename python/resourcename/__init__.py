"""Google AIP-122 resource-name templates.

Import the package and use it directly:

    import resourcename

    t = resourcename.ResourceTemplate("//music.example.com/artists/{artist_id}")
    t.generate(artist_id="radiohead")
    t.parse("//music.example.com/artists/radiohead")
"""

from .template import ResourceTemplate, resource

__all__ = ["ResourceTemplate", "resource"]
