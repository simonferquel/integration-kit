package exec

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/simonferquel/integration-kit"
)

func TestParallelism(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := &integrationkit.Cluster{
		Nodes: []*integrationkit.Node{
			&integrationkit.Node{
				Name: "n1",
			},
			&integrationkit.Node{
				Name: "n2",
			},
			&integrationkit.Node{
				Name: "n3",
			},
		},
	}
	schedulingResults := make(map[string]int)

	d := NewDispatcher(ctx, c)

	wg := sync.WaitGroup{}
	mut := sync.Mutex{}
	wg.Add(3)

	tStart := time.Now()
	job := func(ctx context.Context, n *integrationkit.Node) error {
		time.Sleep(100 * time.Millisecond)
		mut.Lock()
		defer mut.Unlock()
		schedulingResults[n.Name]++
		wg.Done()
		return nil
	}
	go func() {
		d.Run(ctx, func(*integrationkit.Node) bool { return true }, false, job)
	}()
	go func() {
		d.Run(ctx, func(*integrationkit.Node) bool { return true }, false, job)
	}()
	go func() {
		d.Run(ctx, func(*integrationkit.Node) bool { return true }, false, job)
	}()

	wg.Wait()

	if time.Since(tStart) > 200*time.Millisecond {
		t.Error("should have taken less than 200 ms")
	}
	if schedulingResults["n1"] != 1 {
		t.Errorf("work not correctly distributed. n1 had %v tasks, expected 1", schedulingResults["n1"])
	}
	if schedulingResults["n2"] != 1 {
		t.Errorf("work not correctly distributed. n2 had %v tasks, expected 1", schedulingResults["n2"])
	}
	if schedulingResults["n3"] != 1 {
		t.Errorf("work not correctly distributed. n3 had %v tasks, expected 1", schedulingResults["n3"])
	}
}

func TestWholeClusterLock(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := &integrationkit.Cluster{
		Nodes: []*integrationkit.Node{
			&integrationkit.Node{
				Name: "n1",
			},
			&integrationkit.Node{
				Name: "n2",
			},
			&integrationkit.Node{
				Name: "n3",
			},
		},
	}

	wg := sync.WaitGroup{}
	wg.Add(3)
	tStart := time.Now()
	d := NewDispatcher(ctx, c)
	job := func(ctx context.Context, n *integrationkit.Node) error {
		time.Sleep(100 * time.Millisecond)
		wg.Done()
		return nil
	}

	go func() {
		d.Run(ctx, func(*integrationkit.Node) bool { return true }, true, job)
	}()
	go func() {
		d.Run(ctx, func(*integrationkit.Node) bool { return true }, true, job)
	}()
	go func() {
		d.Run(ctx, func(*integrationkit.Node) bool { return true }, false, job)
	}()

	wg.Wait()

	if time.Since(tStart) < 300*time.Millisecond {
		t.Error("should have taken more than 300 ms")
	}
}

func TestPredicates(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := &integrationkit.Cluster{
		Nodes: []*integrationkit.Node{
			&integrationkit.Node{
				Name:         "l1",
				HostPlatform: integrationkit.Platform{OS: "linux", Arch: "x86_64"},
			},
			&integrationkit.Node{
				Name:         "l2",
				HostPlatform: integrationkit.Platform{OS: "linux", Arch: "x86_64"},
			},
			&integrationkit.Node{
				Name:         "w1",
				HostPlatform: integrationkit.Platform{OS: "windows", Arch: "x86_64"},
			},
			&integrationkit.Node{
				Name:         "w2",
				HostPlatform: integrationkit.Platform{OS: "windows", Arch: "x86_64"},
			},
		},
	}

	schedulingResults := make(map[string]int)

	d := NewDispatcher(ctx, c)

	wg := sync.WaitGroup{}
	mut := sync.Mutex{}
	wg.Add(2)

	job := func(ctx context.Context, n *integrationkit.Node) error {
		mut.Lock()
		defer mut.Unlock()
		schedulingResults[n.Name]++
		wg.Done()
		return nil
	}
	go func() {
		d.Run(ctx, integrationkit.IsOS("linux"), false, job)
	}()
	go func() {
		d.Run(ctx, integrationkit.IsOS("windows"), false, job)
	}()

	wg.Wait()

	if schedulingResults["l1"]+schedulingResults["l2"] != 1 {
		t.Error("work not correctly distributed. invalid count of linux jobs")
	}
	if schedulingResults["w1"]+schedulingResults["w2"] != 1 {
		t.Error("work not correctly distributed. invalid count of windows jobs")
	}
}
