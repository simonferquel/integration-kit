package integrationkit

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"

	"github.com/docker/docker/client"
)

func findInTuples(key string, tuples [][2]string) bool {
	for _, t := range tuples {
		if t[0] == key {
			return true
		}
	}
	return false
}
func isSwarmClassic(info *types.Info) bool {
	return findInTuples("Strategy", info.SystemStatus) && findInTuples("Filters", info.SystemStatus) && findInTuples("Nodes", info.SystemStatus)
}
func findSwarmClassicSupportedPlatforms(info *types.Info) []Platform {
	result := []Platform{}
	for _, t := range info.SystemStatus {
		if t[0] == "  └ Labels" {
			if indexOfOsType := strings.Index(t[1], "ostype="); indexOfOsType != -1 {
				res := t[1][indexOfOsType+7:]
				res = res[:strings.Index(res, ",")]
				result = append(result, Platform{OS: res})
			}
		}
	}
	return result
}

// DiscoverNode calls necessary docker APIs to get Node information to populate a cluster description
func DiscoverNode(ctx context.Context, dc *client.Client) (*Node, error) {
	dc.NegotiateAPIVersion(ctx)
	info, err := dc.Info(ctx)
	if err != nil {
		return nil, err
	}
	v, err := dc.ServerVersion(ctx)
	if err != nil {
		return nil, err
	}
	n := Node{
		Experimental:             info.ExperimentalBuild,
		HostPlatform:             Platform{Arch: info.Architecture, OS: info.OSType},
		IsSwarmClassicController: isSwarmClassic(&info),
		IsSwarmManager:           info.Swarm.ControlAvailable,
		Name:                     info.Name,
		SupportedPlatforms:       []Platform{{Arch: info.Architecture, OS: info.OSType}},
	}
	if n.HostPlatform.OS == "windows" {
		if strings.Contains(info.Driver, "lcow") {
			n.SupportedPlatforms = append(n.SupportedPlatforms, Platform{Arch: info.Architecture, OS: "lcow"})
		}
	}
	if n.IsSwarmClassicController {
		if n.HostPlatform.OS == "" {
			n.HostPlatform.OS = info.OperatingSystem
			n.SupportedPlatforms[0].OS = info.OperatingSystem
		}
		n.SupportedPlatforms = append(n.SupportedPlatforms, findSwarmClassicSupportedPlatforms(&info)...)
	}
	n.MinAPIVersion, err = ParseAPIVersion(v.MinAPIVersion)
	if err != nil {
		return nil, err
	}
	n.MaxAPIVersion, err = ParseAPIVersion(v.APIVersion)

	if err != nil {
		return nil, err
	}
	return &n, nil
}
