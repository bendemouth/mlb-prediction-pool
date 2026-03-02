export function toProfilePath(username: string): string {
  return `/profile/${encodeURIComponent(username)}`;
}

export function normalizeUsername(username: string): string {
  return decodeURIComponent(username).trim().toLowerCase();
}