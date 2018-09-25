# How to run this script
#   powershell -ExecutionPolicy ByPass -File install.ps1
# Or from the web
#   Set-ExecutionPolicy Bypass -Scope Process -Force
#   iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/cisco-sso/kdk/master/files/install.ps1'))

#####################################################################
# Constants

$version_url = 'https://api.github.com/repos/cisco-sso/kdk/releases/latest'
$tmp_dir = "${ENV:TEMP}\kdk-install"
$install_dir = "C:\Users\${ENV:USERNAME}\AppData\Local\Microsoft\WindowsApps"
$os = "windows"

#####################################################################
# Functions

Function Get-Kdk-Latest-Version($url) {
	# Enable SSL for the WebRequest
	$AllProtocols = [System.Net.SecurityProtocolType]'Ssl3,Tls,Tls11,Tls12'
	[System.Net.ServicePointManager]::SecurityProtocol = $AllProtocols

	# Make the webrequest and parse the tag_name from the JSON
	$version = (Invoke-WebRequest $url | ConvertFrom-Json | Select tag_name).tag_name
	if (-Not ($version)) {
		echo "FATAL: Unable to find latest release version.  WebAPI request must have failed"
		exit
	}
	return $version
}

Function Get-Cpu-Arch() {
	$arch = "386"
	if ($ENV:PROCESSOR_ARCHITECTURE -eq "AMD64") {
		$arch = "amd64"
	}
	return $arch
}

Function Expand-Tar($tarFile, $dest) {
	if (-not (Get-Command Expand-7Zip -ErrorAction Ignore)) {
		Install-Package -Scope CurrentUser -Force 7Zip4PowerShell > $null
	}
	Expand-7Zip -ArchiveFileName "$tarFile" -TargetPath "$dest"
}

Function Expand-Gzip($infile) {
    $outFile = $infile.Substring(0, $infile.LastIndexOfAny('.'))
    $input = New-Object System.IO.FileStream $inFile, ([IO.FileMode]::Open), ([IO.FileAccess]::Read), ([IO.FileShare]::Read)
    $output = New-Object System.IO.FileStream $outFile, ([IO.FileMode]::Create), ([IO.FileAccess]::Write), ([IO.FileShare]::None)
    $gzipStream = New-Object System.IO.Compression.GzipStream $input, ([IO.Compression.CompressionMode]::Decompress)

    $buffer = New-Object byte[](1024)
    while($true){
        $read = $gzipstream.Read($buffer, 0, 1024)
        if ($read -le 0){break}
        $output.Write($buffer, 0, $read)
        }

    $gzipStream.Close()
    $output.Close()
    $input.Close()
}

#####################################################################
# Code

$version = Get-Kdk-Latest-Version $version_url
echo "Latest version: ${version}"

$arch = Get-Cpu-Arch
echo "Operating System: ${os}"
echo "CPU Architecture: ${arch}"
echo "Will install to: $install_dir\kdk.exe"


$dist_tar="kdk-${version}-${os}-${arch}.tar"
$dist_tgz="kdk-${version}-${os}-${arch}.tar.gz"
$download_url="https://github.com/cisco-sso/kdk/releases/download/${version}/${dist_tgz}"

# Create the tmp directory, redirecting to /dev/null
echo "Creating temporary directory"
echo "  $tmp_dir"
New-Item -Force -ItemType directory -Path $tmp_dir | out-null

# Enter the tmp directory
cd $tmp_dir

# Download the tgz
echo "Downloading from url"
echo "  $download_url"
Invoke-WebRequest -Uri $download_url -OutFile $dist_tgz

# Extract the gzip
echo "Extracting gzip file"
echo "  $tmp_dir\$dist_tgz"
Expand-Gzip $tmp_dir\$dist_tgz

# Extract the tar file contents
echo "Extracting tar file"
echo "  $tmp_dir\$dist_tar"
Expand-Tar $dist_tar $tmp_dir

# Copy the kdk binary to the install location
echo "Copying binary"
echo "  from $tmp_dir\kdk.exe"
echo "    to $install_dir\kdk.exe"
Copy-Item -Force -Path $tmp_dir\kdk.exe -Destination $install_dir\kdk.exe

# Get out of the tmp_dir
cd $install_dir

# Clean up the tmp_dir
echo "Cleaning up by removing directory"
echo "  $tmp_dir"
Remove-Item -Force -Recurse -Path $tmp_dir

echo ""
echo "DONE"
