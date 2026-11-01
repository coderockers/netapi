<?php
/**
 * Read the NetAPI compromised IP feed with threat data (JSON Lines).
 *
 * The addresses were caught attacking NetAPI Nightwatch, NetAPI's own honeypot network.
 * The feed is free, needs no token and is rebuilt every hour. Each line is one JSON
 * object: ip, first_seen, last_seen, category (A infection attempt, B password guessing,
 * C exploit), top_threat, threats (type, name, family, cve, mitre, first, last, days),
 * country, asn, as_name, sensors.
 *
 * The response is streamed to a temporary file and read line by line.
 *
 * Run from the command line: php php-example-compromised-ips.php
 */

$apiUrl = 'https://netapi.com/api2/';

$params = [
    'method'       => 'compromised',
    'dataset_type' => 'ip', // 'ip' = last 24 hours, 'ip-all' = every address ever listed
    'format'       => 'jsonl',
    // Optional filters, the same as the list builder on https://netapi.com/compromised-ips/
    // 'category'    => 'B,C',
    // 'country_not' => 'us',
    // 'min_sensors' => 2,
];

$tmp = tmpfile();
$ch = curl_init($apiUrl . '?' . http_build_query($params));
curl_setopt($ch, CURLOPT_FILE, $tmp);
curl_setopt($ch, CURLOPT_TIMEOUT, 120);
$ok = curl_exec($ch);
$httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
$error = curl_error($ch);
curl_close($ch);

rewind($tmp);
if ($ok === false) {
    fwrite(STDERR, "Request failed: $error\n");
    exit(1);
}
if ($httpCode !== 200) {
    // API errors are plain text, e.g. "405 Method Not Allowed: Invalid category: ..."
    fwrite(STDERR, "HTTP $httpCode: " . stream_get_contents($tmp) . "\n");
    exit(1);
}

$categories = [];
$threats = [];
$countries = [];
$listed = [];

while (($line = fgets($tmp)) !== false) {
    $line = trim($line);
    if ($line === '') {
        continue;
    }
    $record = json_decode($line, true, 512, JSON_THROW_ON_ERROR);
    $listed[$record['ip']] = true;
    $categories[$record['category']] = ($categories[$record['category']] ?? 0) + 1;
    $country = $record['country'] ?: '??';
    $countries[$country] = ($countries[$country] ?? 0) + 1;
    foreach ($record['threats'] as $threat) {
        $threats[$threat['type']] = ($threats[$threat['type']] ?? 0) + 1;
    }
}
fclose($tmp);

ksort($categories);
arsort($threats);
arsort($countries);

printf("%s addresses in the list\n", number_format(count($listed)));
echo 'By category: ' . json_encode($categories) . "\n";
echo 'Top attack types: ' . json_encode(array_slice($threats, 0, 5, true)) . "\n";
echo 'Top countries: ' . json_encode(array_slice($countries, 0, 5, true)) . "\n";

foreach (['8.8.8.8', array_key_first($listed)] as $ip) {
    if ($ip !== null) {
        echo $ip . (isset($listed[$ip]) ? ' -> LISTED' : ' -> not listed') . "\n";
    }
}
