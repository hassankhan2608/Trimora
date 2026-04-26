import { useState } from "react";
import ShortenForm from "./components/ShortenForm";
import Result from "./components/Result";
import type { ShortLink } from "./types/api";

export default function App() {
  const [link, setLink] = useState<ShortLink | null>(null);

  return (
    <div className="page">
      <header className="page__header">
        <a className="brand" href="/">
          <span className="brand__mark">tr</span>
          <span className="brand__name">trimora</span>
        </a>
        <a
          className="ghost-link"
          href="https://github.com"
          target="_blank"
          rel="noreferrer"
        >
          GitHub
        </a>
      </header>

      <main className="page__main">
        <section className="hero">
          <h1 className="hero__title">
            Short links, <em>quietly done.</em>
          </h1>
          <p className="hero__subtitle">
            Paste a URL, pick a custom alias if you like, and share something
            shorter.
          </p>
        </section>

        <section className="card" aria-label="Create a short link">
          <ShortenForm onCreated={setLink} />
          {link && <Result link={link} />}
        </section>
      </main>

      <footer className="page__footer">
        <span>Trimora · a small URL shortener</span>
      </footer>
    </div>
  );
}
