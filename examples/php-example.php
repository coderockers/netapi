<?php
/**
 * Download a zone domain list from NetAPI and print the first rows.
 *
 * The response is a gzip-compressed CSV; it is decompressed with gzdecode().
 * For very large zones (.com, all-zones) save the file to disk and read it
 * with gzopen()/gzgets() instead of loading it into memory.
 *
 * Run from the command line: php php-example.php
 */

$apiUrl = 'https://netapi.com/api2/';
$token = 'YOUR_API_TOKEN'; // https://netapi.com/dashboard/

$params = [
    'method'       => 'download',
    'zone_tld'     => 'net',      // any TLD from ?method=zones, or 'all-zones'
    'dataset_type' => 'list',     // 'list' or 'dataset'
    'filter_type'  => 'active',   // 'active', 'new' or 'deleted'
    'token'        => $token,
];

$ch = curl_init($apiUrl . '?' . http_build_query($params));
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
curl_setopt($ch, CURLOPT_FOLLOWLOCATION, true);
$response = curl_exec($ch);
$httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
curl_close($ch);

if ($httpCode !== 200) {
    // Errors are short plain-text messages, e.g. "403 Forbidden: No active/paid plan."
    fwrite(STDERR, "HTTP $httpCode: $response\n");
    exit(1);
}

$csv = gzdecode($response);
if ($csv === false) {
    fwrite(STDERR, "Could not decompress the response\n");
    exit(1);
}

$lines = explode("\n", $csv);
echo "Header: " . $lines[0] . "\n";
foreach (array_slice($lines, 1, 10) as $line) {
    $row = str_getcsv($line);
    echo $row[0] . "\n";
}
