$ErrorActionPreference = 'Stop'

$moduleDir = (& go list -m -f '{{.Dir}}' github.com/qtgolang/SunnyNet).Trim()
if (-not $moduleDir) {
    throw 'SunnyNet module directory was not found'
}

$headerPath = Join-Path $moduleDir 'src\iphlpapi\c_iphlpapi_tcp.h'
$sourcePath = Join-Path $moduleDir 'src\iphlpapi\c_iphlpapi_tcp.c'
if (-not (Test-Path -LiteralPath $headerPath) -or -not (Test-Path -LiteralPath $sourcePath)) {
    throw "SunnyNet Windows sources were not found under $moduleDir"
}

# SunnyNet v1.4.4 redeclares MIB_TCPROW2/MIB_TCPTABLE2, which current MinGW
# already provides. It also uses a typedef with the same name as GetTcpTable2.
# Patch the downloaded module cache until the upstream dependency includes the
# equivalent fix. The operation is idempotent and deliberately version-local.
foreach ($path in @($headerPath, $sourcePath)) {
    (Get-Item -LiteralPath $path).IsReadOnly = $false
}

$header = [IO.File]::ReadAllText($headerPath)
$header = [regex]::Replace(
    $header,
    '(?s)// 定义 MIB_TCPROW2 结构体.*?} MIB_TCPTABLE2, \*PMIB_TCPTABLE2;\s*',
    ''
)
$header = $header.Replace(
    'typedef DWORD (WINAPI * GetTcpTable2)(PMIB_TCPTABLE2 TcpTable, PULONG SizePointer, BOOL Order);',
    'typedef DWORD (WINAPI * GETTCPTABLE2)(PMIB_TCPTABLE2 TcpTable, PULONG SizePointer, BOOL Order);'
)
[IO.File]::WriteAllText($headerPath, $header)

$source = [IO.File]::ReadAllText($sourcePath)
$source = $source.Replace('GetTcpTable2 pGetTcpTable2;', 'GETTCPTABLE2 pGetTcpTable2;')
$source = $source.Replace(
    'pGetTcpTable2 = (GetTcpTable2) GetProcAddress( hModule, "GetTcpTable2" );',
    'pGetTcpTable2 = (GETTCPTABLE2) GetProcAddress( hModule, "GetTcpTable2" );'
)
[IO.File]::WriteAllText($sourcePath, $source)

if ($header.Contains('typedef struct _MIB_TCPROW2') -or
    -not $header.Contains('typedef DWORD (WINAPI * GETTCPTABLE2)') -or
    -not $source.Contains('GETTCPTABLE2 pGetTcpTable2;')) {
    throw 'SunnyNet Windows compatibility patch verification failed'
}

Write-Host "Patched SunnyNet Windows sources: $moduleDir"
