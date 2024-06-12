package common

import (
	"fmt"
	"time"
)

// AppName is the application name.
const AppName = "App"

// DefaultVersion is the application version.
const DefaultVersion = "0.1.0"

// DefaultDate when not passed in by the compiler.
var DefaultDate = time.Now().Format(time.RFC3339)

// DefaultMetricInterval is the metric interval.
const DefaultMetricInterval = time.Duration(1) * time.Minute

// VendorName is the vendor named used for versioning schemes that depend on a vendor name
// we use the github name for convince.
const VendorName = "synkube"

// BuildInfo will contains build info from https://goreleaser.com/cookbooks/using-main.version
// it is set at compile time by default. If it cannot be, we attempt to derive it at runtime.
type BuildInfo struct {
	version        string
	name           string
	description    string
	date           string
	metricInterval time.Duration
}

// NewBuildInfo creates a build info struct from buildtime data
// it sets sensible defaults.
func NewBuildInfo(version, name, date string) BuildInfo {
	return BuildInfo{
		version: version,
		name:    name,
		date:    date,
	}
}

// Version of the build.
func (b BuildInfo) Version() string {
	return b.version
}

// Name of the application.
func (b BuildInfo) Name() string {
	return b.name
}

// Name of the application.
func (b BuildInfo) Description() string {
	return b.name + ": " + b.VersionString()
}

// Date the application was built.
func (b BuildInfo) Date() string {
	return b.date
}

// MetricInterval the interval to record metrics at.
func (b BuildInfo) MetricInterval() time.Duration {
	return b.metricInterval
}

// VersionString pretty prints a version string with the info above.
func (b BuildInfo) VersionString() string {
	return fmt.Sprintf("%s: (date: %s) \n", b.version, b.date)
}
