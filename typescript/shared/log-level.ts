/**
 * @packageDocumentation Severity levels for `shared/logging` (package diagnostics).
 *
 * @remarks Higher numeric **weight** in {@link LOG_LEVEL_WEIGHT} means more important when filtering.
 */
export enum LogLevel {
	/** No log output. */
	SILENT = "silent",
	/** Verbose diagnostics (default in non-production). */
	DEBUG = "debug",
	/** Operational milestones. */
	INFO = "info",
	/** Recoverable issues. */
	WARN = "warn",
	/** Failures and thrown-error context. */
	ERROR = "error",
}

/** @internal Order used for min-level filtering (must match {@link LogLevel}). */
export const LOG_LEVEL_WEIGHT: Record<LogLevel, number> = {
	[LogLevel.SILENT]: 100,
	[LogLevel.DEBUG]: 10,
	[LogLevel.INFO]: 20,
	[LogLevel.WARN]: 30,
	[LogLevel.ERROR]: 40,
};
