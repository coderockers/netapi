<?php
/*
 * NetAPI JSON API (/api-json/): JSON answers for integrations and AI tools.
 * Free methods need no token; the paid ones (and higher limits) take the token as a Bearer header.
 * PHP 8+, cURL extension.
 */

const TOKEN = 'YOUR_API_TOKEN';
const BASE = 'https://netapi.com/api-json/';

/**
 * Run one method; throws on an API error (errors are JSON too: {"error": {"code", "message"}}).
 */
function call(string $method, array $params = [], ?string $token = null): array
{
    $url = BASE . '?' . http_build_query(['method' => $method] + array_filter($params, static fn($v) => $v !== null));
    $headers = ['Accept: application/json'];
    if ($token !== null) {
        $headers[] = 'Authorization: Bearer ' . $token;
    }
    $curl = curl_init($url);
    curl_setopt_array($curl, [CURLOPT_RETURNTRANSFER => true, CURLOPT_HTTPHEADER => $headers, CURLOPT_TIMEOUT => 60]);
    $body = curl_exec($curl);
    $status = (int)curl_getinfo($curl, CURLINFO_HTTP_CODE);
    curl_close($curl);
    $data = json_decode((string)$body, true) ?? [];
    if ($status !== 200) {
        if ($status === 429) {
            fwrite(STDERR, 'rate limited, retry after ' . ($data['error']['retry_after'] ?? '?') . " s\n");
        }
        throw new RuntimeException($status . ' ' . ($data['error']['code'] ?? '') . ': ' . ($data['error']['message'] ?? $body));
    }
    return $data;
}

// 1. compromised check (free)
$check = call('compromised-check', ['value' => 'example.com']);
echo 'example.com listed now: ', var_export($check['listed_now'], true), ' | ever: ', var_export($check['listed_ever'], true), PHP_EOL;

// 2. newly registered domains containing "crypto" (free)
$fresh = call('search-new', ['q' => 'crypto', 'days' => 3, 'limit' => 5]);
foreach ($fresh['results'] as $row) {
    echo $row['added_on'], ' ', $row['domain'], PHP_EOL;
}

// 3. TLD statistics (free)
$stats = call('tld-stats', ['tld' => 'de']);
echo '.de domains: ', $stats['domains'], ' | registry: ', $stats['registry'], PHP_EOL;

// 4. domain lookup (active plan required)
$lookup = call('lookup-domain', ['domain' => 'example.com'], TOKEN);
echo 'nameservers: ', implode(', ', $lookup['dns'] ?? []), ' | expires: ', $lookup['expires_on'] ?? '-', PHP_EOL;
