// Read the NetAPI compromised-domain / compromised-IP feed.
//
// The feed is free and needs no token. The first line is a "#" comment with
// the list name and date; every other line is one domain or IP address.
//
// Requires Node.js 18+ (built-in fetch). No dependencies.
// Usage: node node-example-compromised.js

const API_URL = "https://netapi.com/api2/";

const params = new URLSearchParams({
  method: "compromised",
  dataset_type: "url", // "url", "ip", "url-all" or "ip-all"
});

async function main() {
  const resp = await fetch(`${API_URL}?${params}`);
  if (!resp.ok) {
    throw new Error(`HTTP ${resp.status}: ${await resp.text()}`);
  }

  const text = await resp.text();
  const entries = text
    .split("\n")
    .map((l) => l.trim())
    .filter((l) => l && !l.startsWith("#"));

  console.log(`${entries.length.toLocaleString("en-US")} entries in the current list`);
  console.log("First 10:", entries.slice(0, 10));

  const blocked = new Set(entries);
  for (const domain of ["example.com", entries[0]]) {
    console.log(domain, blocked.has(domain) ? "-> LISTED" : "-> not listed");
  }
}

main().catch((err) => {
  console.error(err.message);
  process.exitCode = 1;
});
