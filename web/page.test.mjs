// The page script finds its elements by id; each one it asks for must be in index.html, or
// the page breaks with nothing on screen to say so.
import { test } from "node:test";
import assert from "node:assert/strict";
import { readFileSync } from "node:fs";

const html = readFileSync(new URL("static/index.html", import.meta.url), "utf8");
const script = readFileSync(new URL("static/app.js", import.meta.url), "utf8");

test("every element app.js looks up is in index.html", () => {
  const ids = [...script.matchAll(/\$\("([\w-]+)"\)/g)].map((m) => m[1]);
  assert.ok(ids.length > 10, `found only ${ids.length} lookups`);
  for (const id of new Set(ids)) {
    assert.match(html, new RegExp(`id="${id}"`), `index.html has no element with id ${id}`);
  }
});

test("the page loads its script, styles and release version", () => {
  for (const ref of ["/static/app.js", "/static/style.css", 'data-version="{{.Version}}"']) {
    assert.ok(html.includes(ref), `index.html lacks ${ref}`);
  }
});
