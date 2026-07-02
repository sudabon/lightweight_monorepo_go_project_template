import { useCallback, useEffect, useState } from "react";

import { fetchHealth } from "../lib/api";
import type { HealthResponse } from "../types/health";

export type LoadState =
  | { type: "loading" }
  | { type: "ready"; data: HealthResponse }
  | { type: "error"; message: string };

export function useHealth() {
  const [state, setState] = useState<LoadState>({ type: "loading" });

  const refresh = useCallback(async (signal?: AbortSignal) => {
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
  }, []);

  useEffect(() => {
    const controller = new AbortController();
    void refresh(controller.signal);
    return () => controller.abort();
  }, [refresh]);

  return { state, refresh };
}
