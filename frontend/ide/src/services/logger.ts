/**
 * Multiverse IDE Frontend Logger
 *
 * Provides structured logging for the frontend application with:
 * - Log levels (debug, info, warn, error)
 * - Trace ID support for correlation with backend logs
 * - Structured context data
 * - Console output with formatting
 */

export type LogLevel = 'debug' | 'info' | 'warn' | 'error';

export interface LogEntry {
  timestamp: string;
  level: LogLevel;
  message: string;
  traceId?: string;
  component?: string;
  context?: Record<string, unknown>;
}

const LOG_LEVEL_PRIORITY: Record<LogLevel, number> = {
  debug: 0,
  info: 1,
  warn: 2,
  error: 3,
};

const LOG_LEVEL_STYLES: Record<LogLevel, string> = {
  debug: 'color: #8be9fd',
  info: 'color: #50fa7b',
  warn: 'color: #ffb86c',
  error: 'color: #ff5555; font-weight: bold',
};

/**
 * Logger class for structured frontend logging
 */
class Logger {
  private static minLevel: LogLevel = 'info';
  private static traceId: string | null = null;
  private static component: string | null = null;
  private static enabled: boolean = true;

  /**
   * Set the minimum log level
   */
  static setLevel(level: LogLevel): void {
    this.minLevel = level;
  }

  /**
   * Set the current trace ID for log correlation
   */
  static setTraceId(id: string | null): void {
    this.traceId = id;
  }

  /**
   * Get the current trace ID
   */
  static getTraceId(): string | null {
    return this.traceId;
  }

  /**
   * Set the current component name for log context
   */
  static setComponent(name: string | null): void {
    this.component = name;
  }

  /**
   * Enable or disable logging
   */
  static setEnabled(enabled: boolean): void {
    this.enabled = enabled;
  }

  /**
   * Create a child logger with a specific component name
   */
  static withComponent(component: string): ComponentLogger {
    return new ComponentLogger(component);
  }

  /**
   * Check if a log level should be output
   */
  private static shouldLog(level: LogLevel): boolean {
    return this.enabled && LOG_LEVEL_PRIORITY[level] >= LOG_LEVEL_PRIORITY[this.minLevel];
  }

  /**
   * Format and output a log entry
   */
  private static log(level: LogLevel, message: string, context?: Record<string, unknown>, component?: string): void {
    if (!this.shouldLog(level)) return;

    const entry: LogEntry = {
      timestamp: new Date().toISOString(),
      level,
      message,
      traceId: this.traceId ?? undefined,
      component: component ?? this.component ?? undefined,
      context,
    };

    // Format for console output
    const prefix = `[${entry.timestamp.split('T')[1].slice(0, 12)}]`;
    const levelStr = level.toUpperCase().padEnd(5);
    const componentStr = entry.component ? `[${entry.component}]` : '';
    const traceStr = entry.traceId ? `[${entry.traceId.slice(0, 8)}]` : '';

    const logFn = level === 'error' ? console.error :
                  level === 'warn' ? console.warn :
                  level === 'debug' ? console.debug :
                  console.log;

    if (context && Object.keys(context).length > 0) {
      logFn(
        `%c${prefix} ${levelStr} ${componentStr}${traceStr} ${message}`,
        LOG_LEVEL_STYLES[level],
        context
      );
    } else {
      logFn(
        `%c${prefix} ${levelStr} ${componentStr}${traceStr} ${message}`,
        LOG_LEVEL_STYLES[level]
      );
    }
  }

  /**
   * Log a debug message
   */
  static debug(message: string, context?: Record<string, unknown>): void {
    this.log('debug', message, context);
  }

  /**
   * Log an info message
   */
  static info(message: string, context?: Record<string, unknown>): void {
    this.log('info', message, context);
  }

  /**
   * Log a warning message
   */
  static warn(message: string, context?: Record<string, unknown>): void {
    this.log('warn', message, context);
  }

  /**
   * Log an error message
   */
  static error(message: string, context?: Record<string, unknown>): void {
    this.log('error', message, context);
  }

  /**
   * Log with component context (internal use)
   */
  static logWithComponent(level: LogLevel, component: string, message: string, context?: Record<string, unknown>): void {
    this.log(level, message, context, component);
  }
}

/**
 * Component-specific logger instance
 */
class ComponentLogger {
  constructor(private component: string) {}

  debug(message: string, context?: Record<string, unknown>): void {
    Logger.logWithComponent('debug', this.component, message, context);
  }

  info(message: string, context?: Record<string, unknown>): void {
    Logger.logWithComponent('info', this.component, message, context);
  }

  warn(message: string, context?: Record<string, unknown>): void {
    Logger.logWithComponent('warn', this.component, message, context);
  }

  error(message: string, context?: Record<string, unknown>): void {
    Logger.logWithComponent('error', this.component, message, context);
  }
}

// Initialize logger based on environment
if (typeof window !== 'undefined') {
  // Enable debug level in development
  const isDev = import.meta.env?.DEV ?? false;
  if (isDev) {
    Logger.setLevel('debug');
  }
}

export { Logger, ComponentLogger };
export default Logger;
