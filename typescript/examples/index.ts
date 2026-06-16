/**
 * `@resourceName` adds `static Resource` at **runtime** only. TypeScript does not infer that,
 * so add a **types-only** `declare static` line — then `Device.Resource.Parse` shows up in the IDE.
 *
 * Alternative with no `declare`: `class Device extends resourceNameBase("...") {}`
 *
 * From `examples/`: `bun ./index.ts` (run `npm run build` in the package root first so `dist/` is current).
 */

import {
	type ClassResource,
	resourceName,
} from "@the-protobuf-project/resourcename";

@resourceName("//system.com/devices/{device_id}")
class Device {
	/**
	 * Types only (no emit). The decorator assigns the real object at runtime.
	 * Without this line, `Device.Resource` is a type error even though it works at runtime.
	 */
	declare static readonly Resource: ClassResource;
}

const parsed = Device.Resource.Parse("//system.com/devices/router-01");
console.log("Parsed resource name:", parsed.device_id);

const name = Device.Resource.Generate({ device_id: "sensor-22" });
console.log("Generated resource name:", name);
