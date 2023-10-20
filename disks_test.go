package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListDisks(t *testing.T) {
	disks, err := ListDisks()
	assert.Nil(t, err, "ListDisks should never return an error")
	assert.NotEmpty(t, disks, "ListDisks should return at least one disk")
}

func TestListDiskPartitions(t *testing.T) {
	partitions, err := ListDiskPartitions()
	assert.Nil(t, err, "ListDiskPartitions should never return an error")
	assert.NotEmpty(t, partitions, "ListDiskPartitions should return at least one disk")
}

func TestGetDisk(t *testing.T) {
	disks, err := ListDisks()
	assert.Nil(t, err, "ListDisks should never return an error")
	assert.NotEmpty(t, disks, "ListDisks should return at least one disk")

	disk, err := GetDisk(disks[0].Device)
	assert.Nil(t, err, "GetDisk should never return an error")
	assert.NotNil(t, disk, "GetDisk should return a disk")
	assert.Equal(t, disks[0].Device, disk.Device, "GetDisk should return the correct disk")
}

func TestGetDiskUsage(t *testing.T) {
	disks, err := ListDisks()
	assert.Nil(t, err, "ListDisks should never return an error")
	assert.NotEmpty(t, disks, "ListDisks should return at least one disk")

	for _, disk := range disks {
		assert.NotNil(t, disk.Usage, "ListDisks result should return disk usage")
	}

	diskUsage, err := GetDiskUsage(disks[0].Device)
	assert.Nil(t, err, "GetDiskUsage should never return an error")
	assert.NotNil(t, diskUsage, "GetDiskUsage should return disk usage")
}
