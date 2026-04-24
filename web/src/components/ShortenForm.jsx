import { useState } from "react";
import { shorten } from "../api.js";

export default function ShortenForm({ onCreated }) {
  const [url, setUrl] = useState("");
  const [alias, setAlias] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit(event) {
    event.preventDefault();
    if (loading) return;
    setError("");
    setLoading(true);
    try {
      const link = await shorten({ url: url.trim(), alias: alias.trim() });
      onCreated(link);
    } catch (err) {
      onCreated(null);
      setError(err.message);
    } finally {
      setLoading(false);
    }
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
