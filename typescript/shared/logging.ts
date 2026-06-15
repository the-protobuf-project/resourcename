/**
 * @packageDocumentation Environment-aware logging for `@protobuf_project/resourcename`.
 *
 * @remarks
 * - **Browser:** `console` only; ANSI colors off unless `RESOURCE_NAME_COLOR=always`.
 * - **Node/Bun:** reads `process.env`; colors follow `RESOURCE_NAME_COLOR` (`auto` = TTY).
 * - **`RESOURCE_NAME_LOG_LEVEL`:** `silent` | `debug` | `info` | `warn` | `error` (default `debug` when `NODE_ENV !== production`, else `info`).
 */

import { LOG_LEVEL_WEIGHT, LogLevel } from "./log-level";

type PackageMeta = {
	name: string;
	version: string;
};

type LoggerRuntimeConfig = {
	packageName: string;
	version: string;
	environment: string;
	colorMode: "auto" | "always" | "never";
	minLevel: LogLevel;
};

const ANSI = {
	reset: "\x1b[0m",
	dim: "\x1b[2m",
	gray: "\x1b[90m",
	white: "\x1b[97m",
	cyan: "\x1b[36m",
	magenta: "\x1b[35m",
	green: "\x1b[32m",
	yellow: "\x1b[33m",
	red: "\x1b[31m",
	blue: "\x1b[34m",
} as const;

const defaultMeta = loadDefaultMeta();
let runtimeConfig: LoggerRuntimeConfig = {
	packageName: defaultMeta.name,
	version: defaultMeta.version,
	environment: resolveEnvironment(),
	colorMode: resolveColorMode(),
	minLevel: resolveMinLogLevel(),
};

let infoBootstrapped = false;
let decoratorInfoBootstrapped = false;

/** @internal */
function resolveCallerLocation(): string {
	const stackRaw = new Error().stack ?? "";
	const lines = stackRaw.split("\n").map((line) => line.trim());
	for (const line of lines) {
		if (!line.includes(":")) continue;
		if (line.includes("shared/logging.ts")) continue;
		if (!line.includes(".ts")) continue;
		const match = line.match(/((\/|[A-Za-z]:\\).+?):(\d+):(\d+)/);
		if (!match?.[1] || !match[3]) continue;
		const filePath = match[1];
		const lineNumber = match[3];
		return `${basename(filePath)}:${lineNumber}`;
	}
	return "unknown:0";
}

function timestampWithOffset(): string {
	const date = new Date();
	const y = date.getFullYear();
	const m = `${date.getMonth() + 1}`.padStart(2, "0");
	const d = `${date.getDate()}`.padStart(2, "0");
	const hh = `${date.getHours()}`.padStart(2, "0");
	const mm = `${date.getMinutes()}`.padStart(2, "0");
	const ss = `${date.getSeconds()}`.padStart(2, "0");
	const offsetMinutes = -date.getTimezoneOffset();
	const sign = offsetMinutes >= 0 ? "+" : "-";
	const abs = Math.abs(offsetMinutes);
	const offH = `${Math.floor(abs / 60)}`.padStart(2, "0");
	const offM = `${abs % 60}`.padStart(2, "0");
	return `${y}-${m}-${d}T${hh}:${mm}:${ss}${sign}${offH}:${offM}`;
}

function levelColor(level: LogLevel): string {
	if (level === LogLevel.ERROR) return ANSI.red;
	if (level === LogLevel.WARN) return ANSI.yellow;
	if (level === LogLevel.INFO) return ANSI.green;
	if (level === LogLevel.DEBUG) return ANSI.gray;
	return ANSI.white;
}

function basename(pathValue: string): string {
	const parts = pathValue.split(/[\\/]/);
	const last = parts[parts.length - 1];
	return last ?? "unknown";
}

/** Safe `process.env` read (Node/Bun); empty string in browsers unless injected by bundler. */
function readEnv(key: string): string {
	const proc = (
		globalThis as {
			process?: { env?: Record<string, string | undefined> };
		}
	).process;
	const value = proc?.env?.[key];
	return typeof value === "string" ? value : "";
}

function isBrowser(): boolean {
	const g = globalThis as { window?: unknown; document?: unknown };
	return g.window !== undefined && g.document !== undefined;
}

function isStdoutColorSupported(): boolean {
	const proc = (
		globalThis as {
			process?: { stdout?: { isTTY?: boolean } };
		}
	).process;
	return proc?.stdout?.isTTY === true;
}

function loadDefaultMeta(): PackageMeta {
	const nameFromEnv = readEnv("RESOURCE_NAME_PACKAGE_NAME");
	const versionFromEnv = readEnv("RESOURCE_NAME_PACKAGE_VERSION");
	return {
		name: nameFromEnv || "@protobuf_project/resourcename",
		version: versionFromEnv || "1.0.0",
	};
}

function resolveEnvironment(): string {
	return readEnv("NODE_ENV") || "development";
}

function resolveColorMode(): "auto" | "always" | "never" {
	const mode = readEnv("RESOURCE_NAME_COLOR");
	if (mode === "always") return "always";
	if (mode === "never") return "never";
	return "auto";
}

function parseLogLevel(raw: string): LogLevel | undefined {
	switch (raw.trim().toLowerCase()) {
		case "silent":
			return LogLevel.SILENT;
		case "debug":
			return LogLevel.DEBUG;
		case "info":
			return LogLevel.INFO;
		case "warn":
		case "warning":
			return LogLevel.WARN;
		case "error":
			return LogLevel.ERROR;
		default:
			return undefined;
	}
}

function resolveMinLogLevel(): LogLevel {
	const fromEnv = parseLogLevel(readEnv("RESOURCE_NAME_LOG_LEVEL"));
	if (fromEnv !== undefined) return fromEnv;
	return resolveEnvironment() === "production" ? LogLevel.INFO : LogLevel.DEBUG;
}

function colorizeLine(
	stamp: string,
	level: LogLevel,
	levelText: string,
	source: string,
	env: string,
	target: string,
	payload: string,
	packageName: string,
	packageVersion: string,
): string {
	const colorMode = runtimeConfig.colorMode;
	const shouldColor =
		colorMode === "always" ||
		(colorMode === "auto" && !isBrowser() && isStdoutColorSupported());
	if (!shouldColor) {
		return `${stamp} ${levelText} <${source}> ${packageName} (${packageVersion} | ${env}): [${target}] ${payload}`;
	}
	const lvlColor = levelColor(level);
	const ts = `${ANSI.dim}${stamp}${ANSI.reset}`;
	const lvl = `${lvlColor}${levelText}${ANSI.reset}`;
	const src = `${ANSI.cyan}<${source}>${ANSI.reset}`;
	const pkg = `${ANSI.magenta}${packageName}${ANSI.reset}`;
	const meta = `${ANSI.gray}(${packageVersion} | ${env})${ANSI.reset}`;
	const tag = `${ANSI.blue}[${target}]${ANSI.reset}`;
	const body = `${ANSI.white}${payload}${ANSI.reset}`;
	return `${ts} ${lvl} ${src} ${pkg} ${meta}: ${tag} ${body}`;
}

/**
 * Low-level logger: level gating, formatted single line, `console` sink.
 *
 * @remarks Prefer {@link ResourceNameLogger} for call sites inside this package.
 */
export const PackageLogger = {
	shouldEmit(minLevel: LogLevel, messageLevel: LogLevel): boolean {
		if (minLevel === LogLevel.SILENT) return false;
		return LOG_LEVEL_WEIGHT[messageLevel] >= LOG_LEVEL_WEIGHT[minLevel];
	},

	log(level: LogLevel, target: string, message: string, details = ""): boolean {
		if (!this.shouldEmit(runtimeConfig.minLevel, level)) {
			return false;
		}
		const stamp = timestampWithOffset();
		const source = resolveCallerLocation();
		const env = runtimeConfig.environment;
		const levelText = level.toUpperCase();
		const payload = details ? `${message} ${details}` : message;
		const line = colorizeLine(
			stamp,
			level,
			levelText,
			source,
			env,
			target,
			payload,
			runtimeConfig.packageName,
			runtimeConfig.version,
		);
		if (level === LogLevel.ERROR) console.error(line);
		else if (level === LogLevel.WARN) console.warn(line);
		else console.log(line);
		return true;
	},

	configure(options: Partial<LoggerRuntimeConfig>): LoggerRuntimeConfig {
		runtimeConfig = { ...runtimeConfig, ...options };
		return runtimeConfig;
	},

	minLevel(): LogLevel {
		return runtimeConfig.minLevel;
	},
};

/**
 * Convenience surface for `@protobuf_project/resourcename`: mostly **debug**, sparse **info**,
 * **warn** / **error** for validation failures.
 */
export const ResourceNameLogger = {
	debug(target: string, message: string, details = ""): boolean {
		return PackageLogger.log(LogLevel.DEBUG, target, message, details);
	},

	info(target: string, message: string, details = ""): boolean {
		return PackageLogger.log(LogLevel.INFO, target, message, details);
	},

	warn(target: string, message: string, details = ""): boolean {
		return PackageLogger.log(LogLevel.WARN, target, message, details);
	},

	error(target: string, message: string, details = ""): boolean {
		return PackageLogger.log(LogLevel.ERROR, target, message, details);
	},

	configure(options: Partial<LoggerRuntimeConfig>): LoggerRuntimeConfig {
		return PackageLogger.configure(options);
	},

	shouldEmit(minLevel: LogLevel, messageLevel: LogLevel): boolean {
		return PackageLogger.shouldEmit(minLevel, messageLevel);
	},

	/**
	 * One-time **info** on first `ResourceTemplate` construction.
	 * @internal
	 */
	bootstrapOnce(): void {
		if (infoBootstrapped) return;
		infoBootstrapped = true;
		PackageLogger.log(
			LogLevel.INFO,
			"resource.runtime",
			"resource-name runtime ready",
			`minLevel=${PackageLogger.minLevel()} env=${runtimeConfig.environment}`,
		);
	},

	/**
	 * Second one-time **info** when a stage-3 `resourceName` class decorator runs.
	 * @internal
	 */
	decoratorBootstrapOnce(): void {
		if (decoratorInfoBootstrapped) return;
		decoratorInfoBootstrapped = true;
		PackageLogger.log(
			LogLevel.INFO,
			"resource.decorator",
			"class decorator pipeline active",
			"",
		);
	},
};
