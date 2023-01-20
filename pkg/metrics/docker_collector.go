// Copyright 2022 The envd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/driver"
)

type dockerCollector struct {
	Metrics
	client        driver.Client
	running       bool
	metricsStream chan Metrics
	done          chan bool
	lastCpu       float64
	lastSysCpu    float64
}

func NewDockerCollector(client driver.Client) Collector {
	return &dockerCollector{
		client: client,
	}
}

func (c *dockerCollector) Stop() error {
	c.done <- true
	return nil
}

func (c *dockerCollector) Watch(ctx context.Context, cid string) chan Metrics {
	if c.running {
		return c.metricsStream
	}
	c.metricsStream = make(chan Metrics)
	c.running = true
	c.done = make(chan bool)
	stats := make(chan *driver.Stats)
	go func() {
		defer close(stats)
		err := c.client.Stats(ctx, cid, stats, c.done)
		if err != nil {
			logrus.Error(err)
		}
		c.running = false
	}()
	go func() {
		defer close(c.metricsStream)
		for s := range stats {
			c.ReadCPU(s)
			c.ReadMem(s)
			c.ReadNet(s)
			c.ReadIO(s)
			c.metricsStream <- c.Metrics
		}
	}()

	return c.metricsStream
}

func (c *dockerCollector) ReadCPU(stats *driver.Stats) {
	ncpus := uint8(stats.CPUStats.OnlineCPUs)
	if ncpus == 0 {
		ncpus = uint8(len(stats.CPUStats.CPUUsage.PercpuUsage))
	}
	total := float64(stats.CPUStats.CPUUsage.TotalUsage)
	system := float64(stats.CPUStats.SystemCPUUsage)

	cpudiff := total - c.lastCpu
	syscpudiff := system - c.lastSysCpu

	c.NCpus = ncpus
	c.CPUUtil = percent(cpudiff, syscpudiff)
	c.lastCpu = total
	c.lastSysCpu = system
	c.Pids = int(stats.PidsStats.Current)
}

func (c *dockerCollector) ReadMem(stats *driver.Stats) {
	c.MemUsage = int64(stats.MemoryStats.Usage - stats.MemoryStats.Stats.Cache)
	c.MemLimit = int64(stats.MemoryStats.Limit)
	c.MemPercent = percent(float64(c.MemUsage), float64(c.MemLimit))
}

func (c *dockerCollector) ReadNet(stats *driver.Stats) {
	var rx, tx int64
	for _, network := range stats.Networks {
		rx += int64(network.RxBytes)
		tx += int64(network.TxBytes)
	}
	c.NetRx, c.NetTx = rx, tx
}

func (c *dockerCollector) ReadIO(stats *driver.Stats) {
	var read, write int64
	for _, blk := range stats.BlkioStats.IOServiceBytesRecursive {
		if blk.Op == "Read" {
			read += int64(blk.Value)
		}
		if blk.Op == "Write" {
			write += int64(blk.Value)
		}
	}
	c.IOBytesRead, c.IOBytesWrite = read, write
}
