import { useEffect, useState } from "react";

export default function Result({ link }) {
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    setCopied(false);
  }, [link]);

  async function handleCopy() {
    try {
      await navigator.clipboard.writeText(link.short_url);
      setCopied(true);
    } catch {
      setCopied(false);
    }
  }

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
    </div>
  );
}
