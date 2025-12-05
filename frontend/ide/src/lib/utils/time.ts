/**
 * Format an ISO date string to a local time string (HH:mm).
 * Uses the environment's local timezone.
 * @param isoString ISO 8601 date string
 * @returns Formatted time string, e.g., "14:30"
 */
export function formatLocalTime(isoString: string): string {
  if (!isoString) return '';
  try {
    const date = new Date(isoString);
    // Format: YYYY-MM-DD HH:mm
    const activeDate = new Date(date);
    const year = activeDate.getFullYear();
    const month = String(activeDate.getMonth() + 1).padStart(2, '0');
    const day = String(activeDate.getDate()).padStart(2, '0');
    const hours = String(activeDate.getHours()).padStart(2, '0');
    const minutes = String(activeDate.getMinutes()).padStart(2, '0');
    
    return `${year}-${month}-${day} ${hours}:${minutes}`;
  } catch {
    console.error('Invalid date string:', isoString);
    return isoString;
  }
}
