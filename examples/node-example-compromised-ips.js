// Read the NetAPI compromised IP feed with threat data (JSON Lines).
//
// The addresses were caught attacking NetAPI Nightwatch, NetAPI's own honeypot network.
// The feed is free, needs no token and is rebuilt every hour. Each line is one JSON
// object: ip, first_seen, last_seen, category (A infection attempt, B password guessing,
// C exploit), top_threat, threats (type, name, family, cve, mitre, first, last, days),
// country, asn, as_name, sensors.
//
// The response is streamed line by line, nothing is buffered in memory.
//
// Requires Node.js 18+ (built-in fetch). No dependencies.
// Usage: node node-example-compromised-ips.js

const { Readable } = require("node:stream");
const readline = require("node:readline");

const API_URL = "https://netapi.com/api2/";

const params = new URLSearchParams({
  method: "compromised",
  dataset_type: "ip", // "ip" = last 24 hours, "ip-all" = every address ever listed
  format: "jsonl",
  // Optional filters, the same as the list builder on https://netapi.com/compromised-ips/
  // category: "B,C",
  // country_not: "us",
  // min_sensors: "2",
});

function bump(counter, key) {
  counter.set(key, (counter.get(key) || 0) + 1);
}

function top(counter, n) {
  return [...counter.entries()].sort((a, b) => b[1] - a[1]).slice(0, n);
}

async function main() {
  const resp = await fetch(`${API_URL}?${params}`);
  if (!resp.ok) {
    // API errors are plain text, e.g. "405 Method Not Allowed: Invalid category: ..."
    throw new Error(`HTTP ${resp.status}: ${await resp.text()}`);
  }

  const categories = new Map();
  const threats = new Map();
  const countries = new Map();
  const listed = new Set();

  const lines = readline.createInterface({ input: Readable.fromWeb(resp.body), crlfDelay: Infinity });
  for await (const line of lines) {
    if (!line) continue;
    const record = JSON.parse(line);
    listed.add(record.ip);
    bump(categories, record.category);
    bump(countries, record.country || "??");
    for (const threat of record.threats) bump(threats, threat.type);
  }

  console.log(`${listed.size.toLocaleString("en-US")} addresses in the list`);
  console.log("By category:", Object.fromEntries([...categories.entries()].sort()));
  console.log("Top attack types:", top(threats, 5));
  console.log("Top countries:", top(countries, 5));

  const first = listed.values().next().value;
  for (const ip of ["8.8.8.8", first]) {
    if (ip) console.log(ip, listed.has(ip) ? "-> LISTED" : "-> not listed");
  }
}

main().catch((err) => {
  console.error(err.message);
  process.exitCode = 1;
});
