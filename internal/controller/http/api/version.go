package api

type BuildInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Built   string `json:"built"`
}

var buildInfo = BuildInfo{Version: "6.1.0", Commit: "none", Built: "23.10.2026"}

func SetBuildInfo(version, commit, built string) {
	buildInfo = BuildInfo{Version: version, Commit: commit, Built: built}
}
