package rplanlib

type Version struct {
	Version   string
	BuildTime string
}

var version = Version{
	Version:   "0.3-g-rc2",
	BuildTime: "",
}
