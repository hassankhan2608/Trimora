import { useState, type FormEvent } from "react";
import { shorten } from "../api";
import type { ExpiresIn, ShortLink } from "../types/api";

interface Props {
  onCreated: (link: ShortLink | null) => void;
}

const expiryChoices: { value: "" | ExpiresIn; label: string }[] = [
  { value: "", label: "Never" },
  { value: "1h", label: "1 hour" },
  { value: "1d", label: "1 day" },
  { value: "7d", label: "7 days" },
  { value: "30d", label: "30 days" },
];

export default function ShortenForm({ onCreated }: Props) {
  const [url, setUrl] = useState("");
  const [alias, setAlias] = useState("");
  const [expiresIn, setExpiresIn] = useState<"" | ExpiresIn>("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (loading) return;
    setError("");
    setLoading(true);
    void (async () => {
      try {
        const payload: Parameters<typeof shorten>[0] = { url: url.trim() };
        const a = alias.trim();
        if (a) payload.alias = a;
        if (expiresIn) payload.expires_in = expiresIn;
        const link = await shorten(payload);
        onCreated(link);
      } catch (err) {
        onCreated(null);
        setError(err instanceof Error ? err.message : "Something went wrong");
      } finally {
        setLoading(false);
      }
    })();
  }

  return (
    <form className="form" onSubmit={handleSubmit} noValidate>
      <label className="field">
        <span className="field__label">URL</span>
        <input
          className="field__input"
          type="url"
          inputMode="url"
          placeholder="https://example.com/long/path"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          required
          autoFocus
        />
      </label>

      <label className="field">
        <span className="field__label">
          Custom alias <span className="field__hint">optional</span>
        </span>
        <div className="field__group">
          <span className="field__prefix">trimora.app/</span>
          <input
            className="field__input field__input--joined"
            type="text"
            placeholder="my-link"
            value={alias}
            onChange={(e) => setAlias(e.target.value)}
            maxLength={32}
            pattern="[A-Za-z0-9_-]*"
          />
        </div>
      </label>

      <label className="field">
        <span className="field__label">
          Expires <span className="field__hint">optional</span>
        </span>
        <select
          className="field__input field__select"
          value={expiresIn}
          onChange={(e) => setExpiresIn(e.target.value as "" | ExpiresIn)}
        >
          {expiryChoices.map((c) => (
            <option key={c.value || "never"} value={c.value}>
              {c.label}
            </option>
          ))}
        </select>
      </label>

      {error && (
        <p className="form__error" role="alert">
          {error}
        </p>
      )}

      <button className="button" type="submit" disabled={loading}>
        {loading ? "Shortening…" : "Shorten link"}
      </button>
    </form>
  );
}
