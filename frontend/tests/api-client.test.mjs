import assert from "node:assert/strict";
import { afterEach, beforeEach, test } from "node:test";
import { ApiError, request } from "../src/api/client.ts";

const originalFetch = globalThis.fetch;
let values;
let reloads;
beforeEach(() => {
  values = new Map([["access_token", "test-token"], ["refresh_token", "test-refresh"], ["user", "{}"], ["unrelated", "keep"]]);
  reloads = 0;
  globalThis.localStorage = {
    getItem: (key) => values.get(key) ?? null,
    removeItem: (key) => values.delete(key),
  };
  globalThis.window = { location: { reload: () => reloads++ } };
});
afterEach(() => {
  globalThis.fetch = originalFetch;
  delete globalThis.localStorage;
  delete globalThis.window;
});

for (const body of [{ error: "invalid authorization token" }, { responseCode: "01", responseMessage: "expired" }, { responseCode: "41", responseMessage: "unauthorized" }]) {
  test(`HTTP 401 clears only this session: ${JSON.stringify(body)}`, async () => {
    globalThis.fetch = async () => Response.json(body, { status: 401 });
    await assert.rejects(request("/claims"), (err) => err instanceof ApiError && err.status === 401 && err.message.length > 0);
    assert.deepEqual([...values], [["unrelated", "keep"]]);
    assert.equal(reloads, 1);
  });
}

test("HTML gateway failures remain useful API errors", async () => {
  globalThis.fetch = async () => new Response("<html>Bad gateway</html>", { status: 502 });
  await assert.rejects(request("/claims"), (err) => err instanceof ApiError && err.status === 502 && err.message.includes("502"));
  assert.equal(reloads, 0);
  assert.equal(values.get("access_token"), "test-token");
});

test("HTTP failures cannot masquerade as success envelopes", async () => {
  globalThis.fetch = async () => Response.json({ responseCode: "00", responseData: {} }, { status: 500 });
  await assert.rejects(request("/claims"), ApiError);
});

test("successful data and authorization header are preserved", async () => {
  globalThis.fetch = async (url, options) => {
    assert.equal(url, "/api/v1/claims");
    assert.equal(options.headers.Authorization, "Bearer test-token");
    return Response.json({ responseCode: "00", responseData: ["test-claim"] });
  };
  assert.deepEqual(await request("/claims"), ["test-claim"]);
});
