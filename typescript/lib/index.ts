/**
 * @packageDocumentation Public API for `@the-protobuf-project/resourcename`.
 *
 * @remarks
 * - **Templates:** {@link ResourceTemplate} compiles `{placeholder}` patterns; segments are `[^/]+`.
 * - **Static API:** {@link createClassResource} returns {@link ClassResource} (`Parse` / `Generate` / metadata).
 * - **Classes:** {@link resourceNameBase} (typed `extends`) or {@link resourceName} (decorator + `declare static`).
 * - **Diagnostics:** {@link ResourceNameLogger} / {@link PackageLogger} in `shared/` — mostly `debug`, two `info` bootstraps, `warn` / `error` for failures.
 */

export {
	LOG_LEVEL_WEIGHT,
	LogLevel,
	PackageLogger,
	ResourceNameLogger,
} from "../shared/index";
export { resourceName, resourceNameBase } from "./decorator";
export {
	type ClassResource,
	createClassResource,
	ResourceTemplate,
} from "./resource-template";
