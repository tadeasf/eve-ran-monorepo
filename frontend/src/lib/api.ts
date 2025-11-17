// API configuration
export const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
export const API_KEY = process.env.NEXT_PUBLIC_API_KEY || '';

// Helper function to get headers with API key
export function getApiHeaders(): HeadersInit {
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
  };

  if (API_KEY) {
    headers['X-API-Key'] = API_KEY;
  }

  return headers;
}

// Helper function for fetch with API key
export async function apiFetch(
  endpoint: string,
  options: RequestInit = {}
): Promise<Response> {
  const url = `${API_URL}${endpoint}`;
  const headers = {
    ...getApiHeaders(),
    ...options.headers,
  };

  return fetch(url, {
    ...options,
    headers,
  });
}

// Helper function for fetch with JSON response
export async function apiFetchJson<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const response = await apiFetch(endpoint, options);
  
  if (!response.ok) {
    throw new Error(`API error: ${response.status} ${response.statusText}`);
  }
  
  return response.json();
}
