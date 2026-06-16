/**
 * Single default import — `import resourcename from "@the-protobuf-project/resourcename"` —
 * then use it directly. `resourceNameBase` gives a typed `static Resource` with no `declare`
 * line and no extra type imports.
 *
 * From `examples/`: `bun ./index.ts` (run `bun run build` in the package root first so `dist/` is current).
 */

import resourcename from "@the-protobuf-project/resourcename";

class Artist extends resourcename.resourceNameBase(
	"//music.example.com/artists/{artist_id}",
) {}

const parsed = Artist.Resource.Parse("//music.example.com/artists/radiohead");
console.log("Parsed resource name:", parsed.artist_id);

const name = Artist.Resource.Generate({ artist_id: "bjork" });
console.log("Generated resource name:", name);
