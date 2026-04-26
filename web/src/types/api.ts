// Shared API types. Mirrors the Go DTOs in api/internal/httpapi/types.go.
// Keep this as the single source of truth for request/response shapes
// used by the web client.

export interface ShortenRequest {
  url: string;
  alias?: string;
}

export interface ShortLink {
  code: string;
  short_url: string;
  url: string;
}

export interface ApiError {
  error: string;
}
