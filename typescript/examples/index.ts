/**
 * `@resourceName` adds `static Resource` at **runtime** only. TypeScript does not infer that,
 * so add a **types-only** `declare static` line — then `Artist.Resource.Parse` shows up in the IDE.
 *
 * Alternative with no `declare`: `class Artist extends resourceNameBase("...") {}`
 *
 * From `examples/`: `bun ./index.ts` (run `npm run build` in the package root first so `dist/` is current).
 */

import {
	type ClassResource,
	resourceName,
} from "@the-protobuf-project/resourcename";

@resourceName("//music.example.com/artists/{artist_id}")
class Artist {
	/**
	 * Types only (no emit). The decorator assigns the real object at runtime.
	 * Without this line, `Artist.Resource` is a type error even though it works at runtime.
	 */
	declare static readonly Resource: ClassResource;
}

const parsed = Artist.Resource.Parse("//music.example.com/artists/radiohead");
console.log("Parsed resource name:", parsed.artist_id);

const name = Artist.Resource.Generate({ artist_id: "bjork" });
console.log("Generated resource name:", name);
