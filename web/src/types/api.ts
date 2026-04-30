// Mirrors api/internal/httpapi/types.go.

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
