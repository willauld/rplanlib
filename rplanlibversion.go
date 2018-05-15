
    package rplanlib

var (
    // See createLibRelease.ps1 for variable definition / values
    vermajor      string
    verminor      string
    verpatch      string
    verstr        string
)

var version = struct {
    major         string
    minor         string
    patch         string
    str           string
} {"0", 
    "3", "16",
    "alpha" }

