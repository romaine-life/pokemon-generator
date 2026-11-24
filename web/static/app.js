// The generator page: asks the API for a creature and draws its card. The address bar always
// names the creature on screen (?seed=…&type=…&stage=…), so reloading or sharing shows it again.
"use strict";

const $ = (id) => document.getElementById(id);
const bodyNames = {
  biped: "Upright", quadruped: "Four-legged", serpent: "Serpent", avian: "Winged", blob: "Blob",
  insect: "Insect", aquatic: "Finned", golem: "Golem", spirit: "Floating",
};
const RECENT_KEY = "pokemon-generator.recent";
let current = null; // the creature on screen
let request = 0; // the newest load; older answers are dropped
// The release, put in every address the page asks for, so answers cached by the browser are
// never from an older release (the server caches those for a day).
const RELEASE = document.body.dataset.version || "dev";

function query(p) {
  const q = new URLSearchParams();
  if (p.seed) q.set("seed", p.seed);
  if (p.type) q.set("type", p.type);
  if (p.stage) q.set("stage", p.stage);
  return q.toString();
}

function fromLocation() {
  const q = new URLSearchParams(location.search);
  return { seed: q.get("seed") || "", type: q.get("type") || "", stage: Number(q.get("stage")) || 0 };
}

async function load(p, push) {
  const mine = ++request;
  const card = $("card");
  card.classList.add("loading");
  $("status").hidden = true;
  try {
    const res = await fetch("/api/monster?" + query(p) + "&v=" + encodeURIComponent(RELEASE));
    const body = await res.json().catch(() => ({}));
    if (!res.ok) throw new Error(body.error || `the server answered ${res.status}`);
    if (mine !== request) return;
    current = { m: body, stage: p.stage || 0 };
    render(body);
    const addr = { seed: body.seed, type: body.typeLock || "", stage: p.stage || 0 };
    const url = "/?" + query(addr);
    if (push) history.pushState(addr, "", url);
    else history.replaceState(addr, "", url);
    remember(body, addr);
  } catch (err) {
    if (mine !== request) return;
    showError(err, p, push);
  } finally {
    if (mine === request) card.classList.remove("loading");
  }
}

function showError(err, p, push) {
  const s = $("status");
  s.replaceChildren(
    document.createTextNode(`Couldn't find that Pokémon: ${err.message}.`),
    document.createElement("br"),
  );
  const retry = document.createElement("button");
  retry.textContent = "Try again";
  retry.onclick = () => load(p, push);
  s.append(retry);
  s.hidden = false;
}

function drawSprite(canvas, sprite, scale) {
  canvas.width = sprite.width * scale;
  canvas.height = sprite.height * scale;
  const ctx = canvas.getContext("2d");
  ctx.clearRect(0, 0, canvas.width, canvas.height);
  for (let y = 0; y < sprite.height; y++) {
    for (let x = 0; x < sprite.width; x++) {
      const i = sprite.pixels[y * sprite.width + x];
      if (i === 0) continue;
      ctx.fillStyle = sprite.palette[i];
      ctx.fillRect(x * scale, y * scale, scale, scale);
    }
  }
}

function typeBadge(t, small) {
  const span = document.createElement("span");
  span.className = small ? "type small" : "type";
  span.style.background = t.color;
  span.textContent = t.name;
  return span;
}

function statColor(v) {
  if (v < 50) return "#f2665a";
  if (v < 75) return "#f5a142";
  if (v < 95) return "#f2cf3c";
  if (v < 120) return "#8fd14f";
  return "#33c4b0";
}

function render(m) {
  const root = document.documentElement;
  root.style.setProperty("--type", m.types[0].color);
  root.style.setProperty("--type2", (m.types[1] || m.types[0]).color);
  document.title = `${m.name} · #${String(m.number).padStart(4, "0")} · Pokémon Generator`;

  const canvas = $("sprite");
  drawSprite(canvas, m.sprite, 8);
  canvas.setAttribute("aria-label", `Sprite of ${m.name}`);
  canvas.classList.remove("pop");
  void canvas.offsetWidth; // restart the entrance animation
  canvas.classList.add("pop");

  $("card").classList.toggle("shiny", m.shiny);
  $("shiny").hidden = !m.shiny;
  $("legendary").hidden = !m.legendary;
  $("number").textContent = "#" + String(m.number).padStart(4, "0");
  $("name").textContent = m.name;
  $("category").textContent = `The ${m.category} Pokémon`;
  $("types").replaceChildren(...m.types.map((t) => typeBadge(t)));
  $("height").textContent = `${m.heightM} m`;
  $("weight").textContent = `${m.weightKg} kg`;
  $("body").textContent = bodyNames[m.bodyPlan] || m.bodyPlan;
  $("entry").textContent = m.entry;
  $("seed").value = m.seed;
  $("type").value = m.typeLock || "";

  const stats = $("stats");
  stats.replaceChildren();
  for (const s of m.stats) {
    const tr = stats.insertRow();
    tr.innerHTML = `<th></th><td class="v"></td><td><div class="bar"></div></td>`;
    tr.cells[0].textContent = s.name;
    tr.cells[1].textContent = s.value;
    const bar = tr.querySelector(".bar");
    bar.style.setProperty("--c", statColor(s.value));
    requestAnimationFrame(() => (bar.style.width = `${Math.min(100, (s.value / 180) * 100)}%`));
  }
  const total = stats.insertRow();
  total.className = "total";
  total.innerHTML = `<th>Total</th><td class="v">${m.total}</td><td></td>`;

  $("abilities").replaceChildren(...m.abilities.map((a) => {
    const li = document.createElement("li");
    const b = document.createElement("b");
    b.textContent = a.name;
    li.append(b);
    if (a.hidden) {
      const tag = document.createElement("span");
      tag.className = "hidden-tag";
      tag.textContent = "hidden";
      li.append(tag);
    }
    const d = document.createElement("span");
    d.className = "d";
    d.textContent = a.description;
    li.append(d);
    return li;
  }));

  const tbody = $("moves").tBodies[0];
  tbody.replaceChildren();
  for (const mv of m.moves) {
    const tr = tbody.insertRow();
    tr.innerHTML = `<td><b></b><span class="desc"></span></td><td></td><td></td><td class="n"></td><td class="n"></td><td class="n"></td>`;
    tr.querySelector("b").textContent = mv.name;
    tr.querySelector(".desc").textContent = mv.description;
    tr.cells[1].append(typeBadge(mv.type, true));
    tr.cells[2].textContent = mv.category;
    tr.cells[3].textContent = mv.power || "—";
    tr.cells[4].textContent = mv.accuracy;
    tr.cells[5].textContent = mv.pp;
  }

  const evo = $("evolution");
  evo.replaceChildren();
  for (const e of m.evolution) {
    const li = document.createElement("li");
    if (e.how) {
      const how = document.createElement("span");
      how.className = "how";
      how.textContent = e.how;
      li.append(how);
    }
    const btn = document.createElement("button");
    btn.type = "button";
    btn.setAttribute("aria-current", String(e.stage === m.stage));
    const c = document.createElement("canvas");
    drawSprite(c, e.sprite, 2);
    const label = document.createElement("span");
    label.textContent = e.name;
    btn.append(c, label);
    btn.onclick = () => load({ seed: m.seed, type: m.typeLock || "", stage: e.stage }, true);
    li.append(btn);
    evo.append(li);
  }
  if (m.evolution.length === 1) {
    const li = document.createElement("li");
    li.className = "note";
    li.textContent = m.legendary ? "Legendaries don't evolve." : "It doesn't evolve.";
    evo.append(li);
  }

  $("download").href = "/sprite.png?" + query({ seed: m.seed, type: m.typeLock, stage: m.stage }) + "&scale=8&download=1&v=" + encodeURIComponent(RELEASE);
}

function remember(m, addr) {
  let list = [];
  try { list = JSON.parse(localStorage.getItem(RECENT_KEY)) || []; } catch { list = []; }
  const key = query({ ...addr, stage: m.stage });
  list = list.filter((r) => r.key !== key);
  list.unshift({ key, name: m.name, number: m.number });
  list = list.slice(0, 16);
  try { localStorage.setItem(RECENT_KEY, JSON.stringify(list)); } catch { /* private mode: no history */ }
  showRecent(list);
}

function showRecent(list) {
  $("recent").replaceChildren(...list.map((r) => {
    const li = document.createElement("li");
    const a = document.createElement("a");
    a.href = "/?" + r.key;
    const img = document.createElement("img");
    img.src = "/sprite.png?" + r.key + "&scale=2&v=" + encodeURIComponent(RELEASE);
    img.alt = "";
    const span = document.createElement("span");
    span.textContent = r.name;
    a.append(img, span);
    a.onclick = (ev) => {
      ev.preventDefault();
      const q = new URLSearchParams(r.key);
      load({ seed: q.get("seed"), type: q.get("type") || "", stage: Number(q.get("stage")) || 0 }, true);
    };
    li.append(a);
    return li;
  }));
}

function surprise() {
  load({ seed: "", type: $("type").value }, true);
}

async function fillTypes() {
  try {
    const res = await fetch("/api/types?v=" + encodeURIComponent(RELEASE));
    if (!res.ok) return;
    const sel = $("type");
    for (const t of await res.json()) {
      const o = document.createElement("option");
      o.value = t.id;
      o.textContent = t.name;
      sel.append(o);
    }
    sel.value = fromLocation().type;
  } catch { /* the lock stays at Any type */ }
}

function start() {
  $("controls").addEventListener("submit", (ev) => {
    ev.preventDefault();
    const seed = $("seed").value.trim();
    if (seed) load({ seed, type: $("type").value }, true);
    else surprise();
  });
  $("random").onclick = surprise;
  $("type").addEventListener("change", surprise);
  $("share").onclick = async () => {
    const btn = $("share");
    try {
      await navigator.clipboard.writeText(location.href);
      btn.textContent = "Copied!";
    } catch {
      btn.textContent = "Copy failed";
    }
    setTimeout(() => (btn.textContent = "Copy link"), 1500);
  };
  document.addEventListener("keydown", (ev) => {
    if (ev.target.closest("input, select, textarea") || ev.ctrlKey || ev.metaKey || ev.altKey) return;
    if (ev.key === " " && !ev.target.closest("button, a")) {
      ev.preventDefault();
      surprise();
    }
    if (current && (ev.key === "ArrowLeft" || ev.key === "ArrowRight")) {
      const m = current.m;
      const stage = m.stage + (ev.key === "ArrowLeft" ? -1 : 1);
      if (stage >= 1 && stage <= m.stages) load({ seed: m.seed, type: m.typeLock || "", stage }, true);
    }
  });
  window.addEventListener("popstate", () => load(fromLocation(), false));

  try { showRecent(JSON.parse(localStorage.getItem(RECENT_KEY)) || []); } catch { /* no history yet */ }
  fillTypes();
  load(fromLocation(), false);
}

start();
