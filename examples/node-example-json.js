#!/usr/bin/env node
// NetAPI JSON API (/api-json/): JSON answers for integrations and AI tools.
// Free methods need no token; the paid ones (and higher limits) take the token as a Bearer header.
// Node.js 18+ (global fetch), no dependencies.

const TOKEN = 'YOUR_API_TOKEN';
const BASE = 'https://netapi.com/api-json/';

async function call(method, params = {}, token = null) {
  const url = new URL(BASE);
  url.searchParams.set('method', method);
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined && value !== null) url.searchParams.set(key, String(value));
  }
  const headers = { Accept: 'application/json' };
  if (token) headers.Authorization = `Bearer ${token}`;
  const response = await fetch(url, { headers });
  const body = await response.json(); // errors are JSON too: {"error": {"code", "message"}}
  if (!response.ok) {
    if (response.status === 429) console.error(`rate limited, retry after ${body.error.retry_after} s`);
    throw new Error(`${response.status} ${body.error.code}: ${body.error.message}`);
  }
  return body;
}

(async () => {
  // 1. compromised check (free)
  const check = await call('compromised-check', { value: 'example.com' });
  console.log('example.com listed now:', check.listed_now, '| ever:', check.listed_ever);

  // 2. newly registered domains containing "crypto" (free)
  const fresh = await call('search-new', { q: 'crypto', days: 3, limit: 5 });
  for (const row of fresh.results) console.log(row.added_on, row.domain);

  // 3. TLD statistics (free)
  const stats = await call('tld-stats', { tld: 'de' });
  console.log('.de domains:', stats.domains, '| registry:', stats.registry);

  // 4. domain lookup (active plan required)
  const lookup = await call('lookup-domain', { domain: 'example.com' }, TOKEN);
  console.log('nameservers:', lookup.dns, '| expires:', lookup.expires_on);
})().catch((error) => {
  console.error(error.message);
  process.exit(1);
});
