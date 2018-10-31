param(
    [string]$vermajor = "0", 
    [string]$verminor = "5",
    [string]$verStr = "alpha",
    [switch]$resetPatch = $false,
    [switch]$tagRepo = $false
)

# file has current patch number mananged by this script
$verPatchFile = ".patchNum" 

#Test-Path $verPatchFile -PathType Leaf
if ($resetPatch -or !(Test-Path $verPatchFile -PathType Leaf)) {
    $zero = 0
    Set-Content $verPatchFile -value $zero 
}
$verpatch = Get-Content .patchNum
[int]$verpatchInt = $verpatch
Write-Output "verpatch is: $verpatch, Int $verpatchInt"

$newPatchNum = 1 + $verpatchInt
[string]$newPatchStr = $newPatchNum

Set-Content $verPatchFile -value $newPatchStr
$verpatch = Get-Content .patchNum
[int]$verpatchInt = $verpatch
Write-Output "verpatch is: $verpatch, Int $verpatchInt"

$rplanlibversionFile = "
    package rplanlib

var (
    // See createLibRelease.ps1 for variable definition / values
    vermajor      string
    verminor      string
    verpatch      string
    verstr        string
)

var Version = struct {
    Major         string
    Minor         string
    Patch         string
    Str           string
} {""$vermajor"", ""$verminor"", ""$verpatch"", ""$verstr"" }
"
#Write-Output "version file will be: $rplanlibversionFile"

$rplanlibversionFileName = "rplanlibversion.go"


if ($tagRepo) {
    Write-Output "Tagging the repo"

    $tagStr = "v$vermajor.$verminor.$verpatch"

    Set-Content $rplanlibversionFileName -value $rplanlibversionFile
    git commit -m "Updating version file for $tagStr" $rplanlibversionFileName 

    git tag -a $tagStr -m "Tag version $tagStr"
    git push origin $tagStr

    # to delete tags if needed:
    #git tag --delete $tagStr
    #git push --delete origin $tagStr
}
else {
    Write-Output "No changes made to repo, Need to use -tagrepo for changes"
}

$libgitver = git describe --always

Write-Output "rplanlib git hash: $libgitver"

