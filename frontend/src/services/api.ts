function apiUrl(path: string): string {
  const baseUrl = process.env.REACT_APP_API_BASE_URL || '';
  return `${baseUrl}${path}`;
}

export default apiUrl;