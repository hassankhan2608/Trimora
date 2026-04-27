// Shared API types. Mirrors the Go DTOs in api/internal/httpapi/types.go.
// Keep this as the single source of truth for request/response shapes
// used by the web client.

export type ExpiresIn = "1h" | "1d" | "7d" | "30d";

export interface ShortenRequest {
  url: string;
  alias?: string;
  expires_in?: ExpiresIn;
}

export interface ShortLink {
  code: string;
  short_url: string;
  url: string;
  expires_at?: string;
}

export interface ApiError {
  error: string;
}
