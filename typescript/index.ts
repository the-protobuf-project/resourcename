/**
 * @packageDocumentation Entry point for `@the-protobuf-project/resourcename`.
 *
 * @remarks
 * Use the default aggregate for a single clean import:
 * ```ts
 * import resourcename from "@the-protobuf-project/resourcename";
 * const t = new resourcename.ResourceTemplate("//music.example.com/artists/{artist_id}");
 * ```
 * or named imports (`import { ResourceTemplate } from "..."`) — both resolve to {@link ./lib/index}.
 */

export * from "./lib/index";

import {
	createClassResource,
	LOG_LEVEL_WEIGHT,
	LogLevel,
	PackageLogger,
	ResourceNameLogger,
	ResourceTemplate,
	resourceName,
	resourceNameBase,
} from "./lib/index";

/**
 * Default aggregate of the runtime API, so consumers can write
 * `import resourcename from "@the-protobuf-project/resourcename"` and use
 * `resourcename.ResourceTemplate`, `resourcename.resourceName`, etc.
 */
const resourcename = {
	ResourceTemplate,
	createClassResource,
	resourceName,
	resourceNameBase,
	PackageLogger,
	ResourceNameLogger,
	LogLevel,
	LOG_LEVEL_WEIGHT,
} as const;

export default resourcename;
