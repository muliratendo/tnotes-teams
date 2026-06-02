/** Standard API envelope from the Go backend. */
export interface APIResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

export class APIError extends Error {
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.name = 'APIError';
    this.status = status;
  }
}

function authHeaders(): HeadersInit {
  const token = localStorage.getItem('jwt_token');
  return token
    ? { Authorization: `Bearer ${token}` }
    : {};
}

/** Parses JSON and unwraps `{ success, data }` responses. */
export async function apiRequest<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const res = await fetch(path, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...authHeaders(),
      ...(options.headers || {}),
    },
  });

  const json: APIResponse<T> = await res.json();

  if (!res.ok || !json.success) {
    throw new APIError(json.error || 'Request failed', res.status);
  }

  return json.data as T;
}

/** Download export files with JWT auth. */
export async function downloadExport(path: string, filename: string): Promise<void> {
  const res = await fetch(path, { headers: authHeaders() });
  if (!res.ok) {
    const json = await res.json().catch(() => ({}));
    throw new APIError((json as APIResponse<unknown>).error || 'Export failed', res.status);
  }
  const blob = await res.blob();
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  a.click();
  URL.revokeObjectURL(url);
}
