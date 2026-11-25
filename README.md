# pokemon-generator

A website that makes up a new Pokémon-style creature every time you ask: its sprite,
name, types, stats, abilities, moves, dex entry and evolution line, all generated.
Any word is a seed, so typing your name always gives the same creature. It runs on
romaine.life's AKS cluster at https://pokemon-generator.romaine.life
([SES-221](https://openproject.home.romaine.life/wp/SES-221)).

Every creature is original: the sprites are procedural pixel art drawn from body plans,
and the names, abilities, moves and dex entries are assembled from word lists here.

## Using it

- **Surprise me!** (or Space) finds a random creature; the address bar then names it
  (`/?seed=tavoki42`), so reloading or sharing the link shows the same one.
- **Find the Pokémon of…** takes any text: `nelson`, `Nelson` and ` nelson ` are the same seed.
- **Type** locks the primary type for every creature found until it's set back to Any type.
- **Evolution line** shows every stage; click one (or ←/→) to see it. `&stage=` in the
  address picks a stage; without it the seed picks one.
- **Copy link** copies the address; **Download sprite** saves the PNG.
- **Recently found** keeps the last 16 in the browser's local storage.

A shared link previews the creature (its name, dex entry and sprite) wherever link previews
are drawn.

## How it fits together

| Path | Holds |
| --- | --- |
| `gen/monster.go` | `Generate`: a seed becomes a family (body plan, types, number of stages, legendary, shiny, name fragment, stat shape, abilities), then the requested stage |
| `gen/types.go`, `archetypes.go` | The eighteen types and nine body plans, and the words each brings |
| `gen/sprite.go`, `canvas.go` | The sprite: a 40×40 grid drawn from shapes per body plan, mirror-symmetric except the tail |
| `gen/color.go`, `png.go` | Palette, shading and outlines; the PNG |
| `gen/text.go` | Moves and dex entries |
| `internal/server` | The HTTP routes |
| `web/static` | The page (`index.html` is a template the server fills with share-card tags) |
| `web/package.json`, `web/page.test.mjs` | The page's checks: `npm run lint` checks `app.js` parses, `npm test` that every element it looks up is in `index.html` |
| `cmd/pokemon-generator` | `main`, and `--self-test`, which the image check runs |

Every random choice comes from SHA-256 of the seed, the type lock and what is being chosen,
so the same address always gives the same creature. Sprite traits are keyed by name
(`ge.f("legH", …)`) rather than drawn in sequence, so every stage of an evolution line
reads the same body and they look related; later stages are drawn larger, with smaller
heads and longer horns and ears. A drawing that leaves the grid is redrawn smaller.

`gen` is a public package: pokosmari imports it at a pinned module version to make its creatures, so a
change here reaches that game only when it takes the new version up.

`gen.Version` names the generator; what a seed gives can change between releases. The page
puts the release in every address it asks for (`v=`), so the browser's cached answers
(kept a day) never outlive a release.

## API

| Call | Answers |
| --- | --- |
| `GET /api/monster?seed=&type=&stage=` | The creature as JSON. No seed: a random one, whose `seed` the answer names. 400 for an unknown type or a stage outside 1–3 |
| `GET /api/types` | The eighteen types: `id`, `name`, badge `color` |
| `GET /sprite.png?seed=&type=&stage=&scale=` | The sprite, each pixel `scale` (1–16, default 8) pixels square; `download=1` saves it under the creature's name |
| `GET /healthz` | `{"status": "ok", "version", "generator"}` |

Each request is logged as one JSON line on stdout.

## Running and testing

```sh
go test ./...
(cd web && npm ci && npm run lint && npm test)
go run ./cmd/pokemon-generator          # http://localhost:8080
go run ./cmd/pokemon-generator --self-test
```

`PREVIEW_OUT=sheet.png go test ./gen -run Preview` draws a contact sheet of 54
creatures, for looking at a sprite change; add `PREVIEW_LINES=1` for evolution lines side
by side.

## How it ships

CI comes from the pipeline-templates suites `go`, `web-app` (with `web-image`) and
`cluster-service`, selected by this repository's entry in pipeline-templates' `fleet.json`.
A pull request runs the Go tests and the page's checks, renders and checks `k8s/`, and builds the image and runs
its `--self-test`. A merge publishes `romainecr.azurecr.io/pokemon-generator` and opens a
pin PR that rewrites the image in `k8s/values.yaml` to the new tag and digest. Argo CD runs
`k8s/` from its Application in infra-bootstrap's `k8s/apps`, in the `pokemon-generator`
namespace; external-dns publishes `pokemon-generator.romaine.life` from the HTTPRoute, and
cert-manager issues its certificate. The pod keeps no state and calls nothing, so it has no
Azure identity and mounts no service account token.
