package integrationkit

import "testing"

func sliceContains(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func sliceContainsAll(slice []string, vals []string) bool {
	for _, v := range vals {
		if !sliceContains(slice, v) {
			return false
		}
	}
	return true
}

func slicesContainsSame(slice1, slice2 []string) bool {
	return sliceContainsAll(slice1, slice2) && sliceContainsAll(slice2, slice1)
}

func TestNodePredicates(t *testing.T) {
	cluster := Cluster{
		HasSwarm:        true,
		HasSwarmClassic: true,
		Nodes: []Node{
			Node{
				Name:                     "linux-amd64-manager",
				Experimental:             false,
				HostPlatform:             Platform{OS: "linux", Arch: "amd64"},
				IsSwarmClassicController: true,
				IsSwarmManager:           true,
				MinAPIVersion:            APIVersion{1, 0},
				MaxAPIVersion:            APIVersion{5, 0},
				SupportedPlatforms:       nil,
			},
			Node{
				Name:                     "win-amd64-manager",
				Experimental:             false,
				HostPlatform:             Platform{OS: "windows", Arch: "amd64"},
				IsSwarmClassicController: true,
				IsSwarmManager:           true,
				MinAPIVersion:            APIVersion{1, 0},
				MaxAPIVersion:            APIVersion{5, 0},
				SupportedPlatforms:       nil,
			},
			Node{
				Name:                     "win-amd64-lcow",
				Experimental:             false,
				HostPlatform:             Platform{OS: "windows", Arch: "amd64"},
				IsSwarmClassicController: false,
				IsSwarmManager:           false,
				MinAPIVersion:            APIVersion{1, 0},
				MaxAPIVersion:            APIVersion{5, 0},
				SupportedPlatforms:       []Platform{Platform{OS: "lcow", Arch: "amd64"}},
			},
			Node{
				Name:                     "linux-arm",
				Experimental:             false,
				HostPlatform:             Platform{OS: "linux", Arch: "arm"},
				IsSwarmClassicController: true,
				IsSwarmManager:           true,
				MinAPIVersion:            APIVersion{1, 0},
				MaxAPIVersion:            APIVersion{5, 0},
				SupportedPlatforms:       nil,
			},
			Node{
				Name:                     "linux-newversion",
				Experimental:             false,
				HostPlatform:             Platform{OS: "linux", Arch: "amd64"},
				IsSwarmClassicController: true,
				IsSwarmManager:           true,
				MinAPIVersion:            APIVersion{4, 0},
				MaxAPIVersion:            APIVersion{7, 0},
				SupportedPlatforms:       nil,
			},
			Node{
				Name:                     "linux-experimental",
				Experimental:             true,
				HostPlatform:             Platform{OS: "linux", Arch: "amd64"},
				IsSwarmClassicController: true,
				IsSwarmManager:           true,
				MinAPIVersion:            APIVersion{1, 0},
				MaxAPIVersion:            APIVersion{5, 0},
				SupportedPlatforms:       nil,
			},
		},
	}

	cases := []struct {
		name      string
		predicate NodePredicate
		matches   []string
	}{
		{
			name:      "experimental",
			predicate: IsExperimental,
			matches:   []string{"linux-experimental"},
		},
		{
			name:      "not experimental",
			predicate: Not(IsExperimental),
			matches:   []string{"linux-amd64-manager", "win-amd64-manager", "win-amd64-lcow", "linux-arm", "linux-newversion"},
		},

		{
			name:      "windows, non lcow",
			predicate: And(SupportsPlatform(Platform{OS: "windows", Arch: "amd64"}), Not(SupportsOS("lcow"))),
			matches:   []string{"win-amd64-manager"},
		},
		{
			name:      "experimental and v5,  or v7",
			predicate: Or(SupportsAPIVersion(APIVersion{7, 0}), And(SupportsAPIVersion(APIVersion{5, 0}), IsExperimental)),
			matches:   []string{"linux-newversion", "linux-experimental"},
		},
	}

	for _, c := range cases {
		nodes := cluster.FindNodes(c.predicate)
		nodeNames := make([]string, len(nodes))
		for ix, n := range nodes {
			nodeNames[ix] = n.Name
		}
		if !slicesContainsSame(nodeNames, c.matches) {
			t.Errorf("case %s: expected: %v, got %v", c.name, c.matches, nodeNames)
		}
	}
}
