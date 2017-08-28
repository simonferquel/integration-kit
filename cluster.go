package integrationkit

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// APIVersion represents a docker api version
type APIVersion struct {
	Major, Minor int
}

// LowerOrEquals returns true if this version is lower or equal to the one passed as parameter
func (v *APIVersion) LowerOrEquals(other APIVersion) bool {
	if v.Major > other.Major {
		return false
	}
	return v.Major < other.Major || v.Minor <= other.Minor
}

// GreaterOrEquals returns true if this version is greater or equal to the one passed as parameter
func (v *APIVersion) GreaterOrEquals(other APIVersion) bool {
	if v.Major < other.Major {
		return false
	}
	return v.Major > other.Major || v.Minor >= other.Minor
}

// MarshalJSON marshals the version in a json string
func (v *APIVersion) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v.%v"`, v.Major, v.Minor)), nil
}

// UnmarshalJSON unmarshals the version from a json string
func (v *APIVersion) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parts := strings.Split(s, ".")
	if len(parts) != 2 {
		return fmt.Errorf("Malformed api version %s. Should be <major>.<minor>", s)
	}
	major, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return err
	}
	minor, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return err
	}
	v.Major = int(major)
	v.Minor = int(minor)
	return nil
}

// ParseAPIVersion parses a version from a non-quoted string
func ParseAPIVersion(s string) (APIVersion, error) {
	var v APIVersion
	err := v.UnmarshalJSON([]byte(`"` + s + `"`))
	return v, err
}

// Node describes a machine running a docker daemon
type Node struct {
	Name                     string            `json:"name"`
	DockerEnv                map[string]string `json:"env"`
	HostPlatform             Platform          `json:"hostPlatform"`
	SupportedPlatforms       []Platform        `json:"supportedPlatforms,omitempty"`
	IsSwarmManager           bool              `json:"swarmManager"`
	IsSwarmClassicController bool              `json:"swarmClassicController"`
	MinAPIVersion            APIVersion        `json:"minApiVersion"`
	MaxAPIVersion            APIVersion        `json:"maxApiVersion"`
	Experimental             bool              `json:"experimental"`
}

// Cluster describes a cluster of Nodes
type Cluster struct {
	Nodes           []*Node `json:"nodes"`
	HasSwarm        bool    `json:"hasSwarm"`
	HasSwarmClassic bool    `json:"hasSwarmClassic"`
}

// NodePredicate is a predicate used to filter nodes
type NodePredicate func(*Node) bool

// SupportsOS returns a predicate for nodes supporting the given OS (with any arch)
func SupportsOS(os OS) NodePredicate {
	return func(n *Node) bool {
		if n.HostPlatform.OS == os {
			return true
		}
		for _, p := range n.SupportedPlatforms {
			if p.OS == os {
				return true
			}
		}
		return false
	}
}

// SupportsPlatform returns a predicate for nodes supporting the given OS/arch
func SupportsPlatform(platform Platform) NodePredicate {
	return func(n *Node) bool {
		if n.HostPlatform == platform {
			return true
		}
		for _, p := range n.SupportedPlatforms {
			if p == platform {
				return true
			}
		}
		return false
	}
}

// IsOS returns true if the node machine is running the specified OS
func IsOS(os OS) NodePredicate {
	return func(n *Node) bool {
		return n.HostPlatform.OS == os
	}
}

// IsPlatform returns true if the node machine is running the specified OS/Arch
func IsPlatform(p Platform) NodePredicate {
	return func(n *Node) bool {
		return n.HostPlatform == p
	}
}

// IsSwarmManager is a predicate matching only swarm managers
var IsSwarmManager NodePredicate = func(n *Node) bool {
	return n.IsSwarmManager
}

// IsSwarmClassicController is a predivate matching only classic swarm controllers
var IsSwarmClassicController NodePredicate = func(n *Node) bool {
	return n.IsSwarmClassicController
}

// SupportsAPIVersion is a predicate matching nodes supporting a given api version
func SupportsAPIVersion(v APIVersion) NodePredicate {
	return func(n *Node) bool {
		return n.MinAPIVersion.LowerOrEquals(v) && n.MaxAPIVersion.GreaterOrEquals(v)
	}
}

// IsExperimental is a predicate matching nodes running docker with experimental features enabled
var IsExperimental NodePredicate = func(n *Node) bool {
	return n.Experimental
}

// And combines multiple predicates in an AND logical operation
func And(predicates ...NodePredicate) NodePredicate {
	return func(n *Node) bool {
		for _, p := range predicates {
			if !p(n) {
				return false
			}
		}
		return true
	}
}

// Or combines multiple predicates in an OR logical operation
func Or(predicates ...NodePredicate) NodePredicate {
	return func(n *Node) bool {
		for _, p := range predicates {
			if p(n) {
				return true
			}
		}
		return false
	}
}

// Not negates an existing predicate
func Not(predicate NodePredicate) NodePredicate {
	return func(n *Node) bool {
		return !predicate(n)
	}
}

// FindNodes returns nodes matching the supplied predicate
func (c *Cluster) FindNodes(predicate NodePredicate) []*Node {
	var result []*Node
	for _, n := range c.Nodes {
		if predicate(n) {
			result = append(result, n)
		}
	}
	return result
}

// ClusterFromNodes inialize a cluster description from a node collection
func ClusterFromNodes(nodes []*Node) *Cluster {
	c := Cluster{Nodes: nodes}
	for _, n := range nodes {
		if n.IsSwarmClassicController {
			c.HasSwarmClassic = true
		}
		if n.IsSwarmManager {
			c.HasSwarm = true
		}
	}
	return &c
}
