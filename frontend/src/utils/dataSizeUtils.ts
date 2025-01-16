/**
 * Constants for file size conversions
 */
const BYTES_IN_MB = 1024 * 1024;      // 1 MB = 1,048,576 bytes
const BYTES_IN_GB = BYTES_IN_MB * 1024; // 1 GB = 1,073,741,824 bytes

/**
 * Converts bytes to gigabytes with specified decimal places
 * @param bytes - Number of bytes to convert
 * @param decimals - Number of decimal places (default: 2)
 * @returns Formatted string with GB unit or number if returnNumber is true
 */
export function bytesToGB(bytes: number, decimals: number = 2): string {
  if (bytes === 0) return '0 GB';
  
  const size = bytes / BYTES_IN_GB;
  return `${size.toFixed(decimals)} GB`;
}

/**
 * Converts bytes to megabytes with specified decimal places
 * @param bytes - Number of bytes to convert
 * @param decimals - Number of decimal places (default: 2)
 * @returns Formatted string with MB unit or number if returnNumber is true
 */
export function bytesToMB(bytes: number, decimals: number = 2): string {
  if (bytes === 0) return '0 MB';
  
  const size = bytes / BYTES_IN_MB;
  return `${size.toFixed(decimals)} MB`;
}

/**
 * Converts megabytes to bytes
 * @param mb - Number of megabytes to convert
 * @returns Number of bytes
 */
export function mbToBytes(mb: number): number {
  return mb * BYTES_IN_MB;
}

/**
 * Converts gigabytes to bytes
 * @param gb - Number of gigabytes to convert
 * @returns Number of bytes
 */
export function gbToBytes(gb: number): number {
  return gb * BYTES_IN_GB;
}

/**
 * Formats a file size in bytes to the most appropriate unit
 * @param bytes - Number of bytes to format
 * @param decimals - Number of decimal places (default: 2)
 * @returns Formatted string with appropriate unit (B, KB, MB, GB, TB)
 */
export function formatFileSize(bytes: number, decimals: number = 2): string {
  if (bytes === 0) return '0 Bytes';

  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  const size = bytes / Math.pow(k, i);
  
  return `${size.toFixed(decimals)} ${sizes[i]}`;
}

/**
 * Converts a formatted file size string back to bytes
 * @param sizeStr - Formatted size string (e.g., "1.5 GB", "500 MB")
 * @returns Number of bytes or null if invalid format
 */
export function parseFileSize(sizeStr: string): number | null {
  const matches = sizeStr.match(/^([\d.]+)\s*(B|KB|MB|GB|TB)$/i);
  if (!matches) return null;

  const size = parseFloat(matches[1]);
  const unit = matches[2].toUpperCase();

  const units: { [key: string]: number } = {
    'B': 1,
    'KB': 1024,
    'MB': BYTES_IN_MB,
    'GB': BYTES_IN_GB,
    'TB': BYTES_IN_GB * 1024
  };

  return size * (units[unit] || 0);
}

// Example usage types
interface FileSizeExample {
  originalSize: number;
  formatted: string;
  inMB: string;
  inGB: string;
}

/**
 * Creates an example of all file size conversions
 * @param bytes - Number of bytes to use in example
 * @returns Object with various conversions
 */
export function getFileSizeExample(bytes: number): FileSizeExample {
  return {
    originalSize: bytes,
    formatted: formatFileSize(bytes),
    inMB: bytesToMB(bytes),
    inGB: bytesToGB(bytes)
  };
}