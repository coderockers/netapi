// Download a zone domain list from NetAPI and stream it line by line.
//
// The response is a gzip-compressed CSV. It is piped through zlib, so even
// the .com list can be processed without buffering it in memory.
//
// Requires Node.js 18+ (built-in fetch). No dependencies.
// Usage: node node-example.js

const { Readable } = require("node:stream");
const { createGunzip } = require("node:zlib");
const readline = require("node:readline");

const API_URL = "https://netapi.com/api2/";
const API_TOKEN = "YOUR_API_TOKEN"; // https://netapi.com/dashboard/

const params = new URLSearchParams({
  method: "download",
  zone_tld: "net",        // any TLD from ?method=zones, or "all-zones"
  dataset_type: "list",   // "list" or "dataset"
  filter_type: "active",  // "active", "new" or "deleted"
  token: API_TOKEN,
});

async function main() {
  const resp = await fetch(`${API_URL}?${params}`);
  if (!resp.ok) {
    // Errors are short plain-text messages, e.g. "403 Forbidden: No active/paid plan."
    throw new Error(`HTTP ${resp.status}: ${await resp.text()}`);
  }

  const lines = readline.createInterface({
    input: Readable.fromWeb(resp.body).pipe(createGunzip()),
    crlfDelay: Infinity,
  });

  let n = 0;
  for await (const line of lines) {
    if (n === 0) {
      console.log("Header:", line);
    } else {
      console.log(line);
    }
    if (++n > 10) break; // remove to process the whole file
  }
}

main().catch((err) => {
  console.error(err.message);
  process.exitCode = 1;
});
