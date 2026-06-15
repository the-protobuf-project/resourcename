/**
 * @packageDocumentation Bidirectional resource-name templates with `{placeholder}` segments.
 *
 * @remarks
 * - Ports the Python reference under `resource-name/py-examplke/`: each `{placeholder}` segment matches `[^/]+`.
 * - Pair with {@link resourceName} / {@link resourceNameBase} and {@link ClassResource} for the static API.
 */

import { ResourceNameLogger } from "../shared/logging";

function extractPlaceholders(template: string): string[] {
	const out: string[] = [];
	for (const m of template.matchAll(/\{([^{}]+)\}/g)) {
		const cap = m[1];
		if (cap !== undefined) {
			out.push(cap);
		}
	}
	return out;
}

function compileRegex(template: string, placeholders: string[]): RegExp {
	let pattern = template;
	for (const ph of placeholders) {
		pattern = pattern.split(`{${ph}}`).join("<<>>");
	}
	pattern = pattern.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
	pattern = pattern.replace(/<<>>/g, "([^/]+)");
	return new RegExp(`^${pattern}$`);
}

/**
 * Immutable template: parse full resource strings and generate names from segment values.
 *
 * @throws {Error} When the template is empty, has no placeholders, or has duplicate placeholder names.
 */
export class ResourceTemplate {
	readonly placeholders: readonly string[];
	private readonly regex: RegExp;
	readonly template: string;

	constructor(template: string) {
		ResourceNameLogger.bootstrapOnce();
		if (!template) {
			ResourceNameLogger.warn(
				"resource.template",
				"rejected empty template",
				"",
			);
			throw new Error("Template cannot be empty");
		}
		const placeholders = extractPlaceholders(template);
		if (placeholders.length === 0) {
			ResourceNameLogger.warn(
				"resource.template",
				"rejected template without placeholders",
				template.slice(0, 80),
			);
			throw new Error("Template must contain at least one placeholder");
		}
		const unique = new Set(placeholders);
		if (unique.size !== placeholders.length) {
			const duplicates = placeholders.filter(
				(p, i) => placeholders.indexOf(p) !== i,
			);
			const dupList = [...new Set(duplicates)];
			ResourceNameLogger.warn(
				"resource.template",
				"duplicate placeholders",
				dupList.join(","),
			);
			throw new Error(
				`Template contains duplicate placeholders: ${dupList.join(", ")}`,
			);
		}
		this.template = template;
		this.placeholders = Object.freeze([...placeholders]);
		this.regex = compileRegex(template, [...placeholders]);
		ResourceNameLogger.debug(
			"resource.template",
			"compiled template",
			template,
		);
	}

	/**
	 * Parse a resource name into placeholder → value (each value is a string segment).
	 *
	 * @param resourcename - Full resource string to match against this template.
	 * @throws {Error} When `resourcename` does not match {@link ResourceTemplate.template}.
	 */
	parse(resourcename: string): Record<string, string> {
		const match = this.regex.exec(resourcename);
		if (!match) {
			ResourceNameLogger.warn(
				"resource.parse",
				"resource name does not match template",
				`name=${resourcename} template=${this.template}`,
			);
			throw new Error(
				`Resource name '${resourcename}' does not match template '${this.template}'`,
			);
		}
		const out: Record<string, string> = {};
		for (let i = 0; i < this.placeholders.length; i++) {
			const ph = this.placeholders[i];
			if (ph !== undefined) {
				out[ph] = match[i + 1] ?? "";
			}
		}
		ResourceNameLogger.debug(
			"resource.parse",
			"parsed resource name",
			resourcename,
		);
		return out;
	}

	/**
	 * Substitute segments into the template. Values must not contain `'/'` (same as the Python helper).
	 *
	 * @param values - Map of placeholder name → segment value.
	 * @throws {Error} On missing/extra keys or values containing `'/'`.
	 */
	generate(values: Record<string, string>): string {
		const missing = this.placeholders.filter((p) => !(p in values));
		if (missing.length > 0) {
			const msg =
				`Missing values for placeholders: ${missing.sort().join(", ")}. ` +
				`Required: [${this.placeholders.join(", ")}], provided: [${Object.keys(values).join(", ")}]`;
			ResourceNameLogger.warn("resource.generate", "missing placeholders", msg);
			throw new Error(msg);
		}
		const extra = Object.keys(values).filter(
			(k) => !this.placeholders.includes(k),
		);
		if (extra.length > 0) {
			const msg = `Unexpected values: ${extra.sort().join(", ")}. Expected only: [${this.placeholders.join(", ")}]`;
			ResourceNameLogger.warn(
				"resource.generate",
				"unexpected placeholder keys",
				msg,
			);
			throw new Error(msg);
		}
		const invalid = Object.entries(values).filter(([, v]) => v.includes("/"));
		if (invalid.length > 0) {
			const msg = `Values must not contain '/': ${JSON.stringify(Object.fromEntries(invalid))}`;
			ResourceNameLogger.warn("resource.generate", "value contains slash", msg);
			throw new Error(msg);
		}
		let result = this.template;
		for (const ph of this.placeholders) {
			const v = values[ph];
			if (v === undefined) {
				ResourceNameLogger.warn("resource.generate", "undefined segment", ph);
				throw new Error(`Missing value for '${ph}'`);
			}
			result = result.split(`{${ph}}`).join(v);
		}
		ResourceNameLogger.debug(
			"resource.generate",
			"generated resource name",
			result,
		);
		return result;
	}
}

/**
 * Static `Device.Resource` API: {@link ClassResource.Parse} / {@link ClassResource.Generate} plus metadata.
 */
export type ClassResource = {
	readonly Template: string;
	readonly Placeholders: readonly string[];
	Parse(resourcename: string): Record<string, string>;
	Generate(values: Record<string, string>): string;
};

/**
 * Build the object used as `static Resource` on {@link resourceNameBase}.
 */
export function createClassResource(template: string): ClassResource {
	const tmpl = new ResourceTemplate(template);
	return {
		Template: tmpl.template,
		Placeholders: tmpl.placeholders,
		Parse: (s) => tmpl.parse(s),
		Generate: (values) => tmpl.generate(values),
	};
}
