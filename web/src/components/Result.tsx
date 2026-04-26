import { useEffect, useState } from "react";
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
