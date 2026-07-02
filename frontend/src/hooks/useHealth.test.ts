import { renderHook, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { fetchHealth } from "../lib/api";
import { useHealth } from "./useHealth";

vi.mock("../lib/api", () => ({
  fetchHealth: vi.fn(),
}));

const mockedFetchHealth = vi.mocked(fetchHealth);

describe("useHealth", () => {
  beforeEach(() => {
    mockedFetchHealth.mockReset();
  });

  it("sets ready when health fetch succeeds", async () => {
    mockedFetchHealth.mockResolvedValue({ status: "ok" });

    const { result } = renderHook(() => useHealth());

    await waitFor(() => expect(result.current.state.type).toBe("ready"));
    expect(result.current.state).toEqual({ type: "ready", data: { status: "ok" } });
  });

  it("sets error when health fetch fails", async () => {
    mockedFetchHealth.mockRejectedValue(new Error("failed"));

    const { result } = renderHook(() => useHealth());

    await waitFor(() => expect(result.current.state.type).toBe("error"));
    expect(result.current.state).toEqual({ type: "error", message: "failed" });
  });
});
