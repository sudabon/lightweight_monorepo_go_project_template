import { useEffect, useState } from "react";

import { fetchHealth } from "./lib/api";
import type { HealthResponse } from "./types/health";

type LoadState =
  | { type: "loading" }
  | { type: "ready"; data: HealthResponse }
  | { type: "error"; message: string };

export default function App() {
  const [state, setState] = useState<LoadState>({ type: "loading" });

  async function refresh(signal?: AbortSignal) {
    setState({ type: "loading" });
    try {
      const data = await fetchHealth(signal);
      setState({ type: "ready", data });
    } catch (error) {
      if (signal?.aborted) {
        return;
      }
      setState({
        type: "error",
        message: error instanceof Error ? error.message : "Unexpected error",
      });
    }
  }

  useEffect(() => {
    const controller = new AbortController();
    void refresh(controller.signal);
    return () => controller.abort();
  }, []);

  const statusText =
    state.type === "ready" ? state.data.status : state.type === "loading" ? "checking" : "offline";

  return (
    <main className="app-shell">
      <section className="status-panel" aria-labelledby="app-title">
        <div className="brand-row">
          <div className="brand-mark" aria-hidden="true" />
          <div>
            <p className="eyebrow">Go + React</p>
            <h1 id="app-title">Lightweight Monorepo Template</h1>
          </div>
        </div>

        <div className="status-grid">
          <div>
            <p className="label">Backend</p>
            <p className="status-value">{statusText}</p>
          </div>
          <span className={`status-dot status-dot--${state.type}`} aria-hidden="true" />
        </div>

        {state.type === "error" ? <p className="error-text">{state.message}</p> : null}

        <button type="button" className="refresh-button" onClick={() => void refresh()}>
          Refresh
        </button>
      </section>
    </main>
  );
}
