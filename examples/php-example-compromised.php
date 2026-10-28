<?php
/**
 * Read the NetAPI compromised-domain / compromised-IP feed.
 *
 * The feed is free and needs no token. The first line is a "#" comment with
 * the list name and date; every other line is one domain or IP address.
 *
 * Run from the command line: php php-example-compromised.php
 */

$apiUrl = 'https://netapi.com/api2/';

$params = [
    'method'       => 'compromised',
    'dataset_type' => 'url', // 'url', 'ip', 'url-all' or 'ip-all'
];

$ch = curl_init($apiUrl . '?' . http_build_query($params));
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
$response = curl_exec($ch);
$httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
curl_close($ch);

if ($httpCode !== 200) {
    fwrite(STDERR, "HTTP $httpCode: $response\n");
    exit(1);
}

$entries = [];
foreach (explode("\n", $response) as $line) {
    $line = trim($line);
    if ($line === '' || $line[0] === '#') {
        continue;
    }
    $entries[] = $line;
}

printf("%s entries in the current list\n", number_format(count($entries)));
echo "First 10:\n";
foreach (array_slice($entries, 0, 10) as $entry) {
    echo "  $entry\n";
}

$blocked = array_flip($entries);
foreach (['example.com', $entries[0]] as $domain) {
    echo $domain . (isset($blocked[$domain]) ? ' -> LISTED' : ' -> not listed') . "\n";
}
