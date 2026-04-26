import type { ApiError, ShortLink, ShortenRequest } from "./types/api";

const baseURL = (import.meta.env.VITE_API_BASE_URL ?? "").replace(/\/$/, "");

export async function shorten(input: ShortenRequest): Promise<ShortLink> {
  const body: ShortenRequest = { url: input.url };
  if (input.alias) body.alias = input.alias;

  const res = await fetch(`${baseURL}/api/shorten`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });

  let data: ShortLink | ApiError | null = null;
  try {
    data = (await res.json()) as ShortLink | ApiError;
  } catch {
    // empty body
  }

  if (!res.ok) {
    const message =
      data && "error" in data ? data.error : `Request failed (${res.status})`;
    throw new Error(message);
  }
  if (!data || !("short_url" in data)) {
    throw new Error("Malformed response from server");
  }
  return data;
}
