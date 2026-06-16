/**
 * @packageDocumentation Stage-3 class decorator and base class for `static Resource` ({@link ClassResource}).
 */

import { ResourceNameLogger } from "../shared/logging";
import { type ClassResource, createClassResource } from "./resource-template";

type AnyCtor = abstract new (...args: never[]) => unknown;

function attachResource(
	ctor: AnyCtor,
	resource: ClassResource,
	template: string,
): void {
	Object.defineProperty(ctor, "Resource", {
		value: resource,
		enumerable: false,
		configurable: true,
		writable: false,
	});
	ResourceNameLogger.debug(
		"resource.decorator",
		"attached static Resource",
		template,
	);
}

/**
 * Return an **abstract** base class with `static readonly Resource` fully typed — extend it so
 * `Artist.Resource.Parse` / `Artist.Resource.Generate` appear in IntelliSense.
 *
 * @param template - Resource pattern with `{placeholder}` segments (e.g. `"//music.example.com/artists/{artist_id}"`).
 *
 * @example
 * ```ts
 * class Artist extends resourceNameBase("//music.example.com/artists/{artist_id}") {}
 * Artist.Resource.Parse("//music.example.com/artists/radiohead");
 * Artist.Resource.Generate({ artist_id: "bjork" });
 * ```
 */
export function resourceNameBase(template: string) {
	const Resource = createClassResource(template);
	abstract class ResourceNameBase {
		static readonly Resource: ClassResource = Resource;
		protected constructor() {}
	}
	ResourceNameLogger.debug(
		"resource.base",
		"created resourceNameBase",
		template,
	);
	return ResourceNameBase;
}

/**
 * Attach `static Resource` at runtime. TypeScript does not infer decorator-added statics; add
 * `declare static readonly Resource: ClassResource` on the class, or use {@link resourceNameBase}.
 *
 * @param template - Resource pattern with `{placeholder}` segments.
 *
 * @remarks ECMAScript stage-3 class decorators (TypeScript 5+; keep `experimentalDecorators` disabled).
 *
 * @example
 * ```ts
 * @resourceName("//music.example.com/albums/{album_id}")
 * class Album {
 *   declare static readonly Resource: ClassResource;
 * }
 * ```
 */
export function resourceName(template: string) {
	const Resource = createClassResource(template);
	return function resourceNameDecorator<T extends AnyCtor>(
		target: T,
		_context: ClassDecoratorContext<T>,
	): void {
		ResourceNameLogger.decoratorBootstrapOnce();
		attachResource(target, Resource, template);
	};
}
