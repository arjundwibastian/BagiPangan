import assert from "node:assert/strict";
import { test } from "node:test";
import { createServer, loadConfigFromFile } from "vite";

test("development routes remain local with specific user-request routing first", async () => {
  const { config } = await loadConfigFromFile({ command: "serve", mode: "development" });
  assert.equal(config.define["import.meta.env.VITE_API_BASE_URL"], '""');
  const routes = config.server.proxy;
  const targets = Object.entries(routes).map(([pattern, route]) => [pattern, route.target]);
  assert.deepEqual(targets, [
    ["^/api/v1/users/[^/]+/food-requests", "http://127.0.0.1:18083"],
    ["/api/v1/food-listings", "http://127.0.0.1:18082"],
    ["/api/v1/food-requests", "http://127.0.0.1:18083"],
    ["/api/v1/claims", "http://127.0.0.1:18084"],
    ["/api/v1/users", "http://127.0.0.1:18081"],
    ["/api/v1/auth", "http://127.0.0.1:18081"],
  ]);
});

test("development transformation overrides an old remote API environment value", async () => {
  const previous = process.env.VITE_API_BASE_URL;
  process.env.VITE_API_BASE_URL = "https://remote-must-never-be-used.example.test";
  let server;
  try {
    server = await createServer({ envDir: false, server: { middlewareMode: true, hmr: false } });
    const transformed = await server.transformRequest("/src/api/client.ts");
    assert.ok(transformed);
    assert.ok(!transformed.code.includes("remote-must-never-be-used"));
  } finally {
    await server?.close();
    if (previous === undefined) delete process.env.VITE_API_BASE_URL;
    else process.env.VITE_API_BASE_URL = previous;
  }
});
