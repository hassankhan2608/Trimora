import { useEffect, useMemo, useState } from "react";
import type { ShortLink } from "../types/api";

interface Props {
  link: ShortLink;
}

export default function Result({ link }: Props) {
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    setCopied(false);
  }, [link]);

  function handleCopy() {
    void navigator.clipboard
      .writeText(link.short_url)
      .then(() => setCopied(true))
      .catch(() => setCopied(false));
  }

  const expiry = useMemo(() => formatExpiry(link.expires_at), [link.expires_at]);

  return (
    <div className="result" aria-live="polite">
      <div className="result__label">Your short link</div>
      <div className="result__row">
        <a
          className="result__link"
          href={link.short_url}
          target="_blank"
          rel="noreferrer"
        >
          {link.short_url}
        </a>
        <button
          type="button"
          className="button button--ghost"
          onClick={handleCopy}
        >
          {copied ? "Copied" : "Copy"}
        </button>
      </div>
      <div className="result__target" title={link.url}>
        → {link.url}
      </div>
      {expiry && <div className="result__expiry">Expires {expiry}</div>}
    </div>
  );
}

function formatExpiry(iso?: string): string | null {
  if (!iso) return null;
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) return null;
  const diffMs = date.getTime() - Date.now();
  const absolute = date.toLocaleString(undefined, {
    dateStyle: "medium",
    timeStyle: "short",
  });
  if (diffMs <= 0) return `on ${absolute} (expired)`;
  return `${relative(diffMs)} · ${absolute}`;
}

function relative(ms: number): string {
  const minutes = Math.round(ms / 60000);
  if (minutes < 60) return `in ${minutes} minute${minutes === 1 ? "" : "s"}`;
  const hours = Math.round(minutes / 60);
  if (hours < 48) return `in ${hours} hour${hours === 1 ? "" : "s"}`;
  const days = Math.round(hours / 24);
  return `in ${days} day${days === 1 ? "" : "s"}`;
}
