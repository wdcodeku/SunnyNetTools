//go:build !windows

package Service

func (g *AppMain) RequestDeviceElevation(mode int) string {
	return "ready"
}
