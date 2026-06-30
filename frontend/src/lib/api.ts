import type { HealthResponse } from "../types/health";

const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080").replace(/\/$/, "");

export async function fetchHealth(signal?: AbortSignal): Promise<HealthResponse> {
  const response = await fetch(`${API_BASE_URL}/health`, { signal });

  if (!response.ok) {
    throw new Error(`Health check failed with ${response.status}`);
  }

  return response.json() as Promise<HealthResponse>;
}
