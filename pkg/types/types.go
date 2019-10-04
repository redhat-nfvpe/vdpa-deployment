// SPDX-License-Identifier: Apache-2.0
// Copyright(c) 2019 Red Hat, Inc.

//
// Package types implements the GO types for the project.
//

package types

//
// Exported Types
//
const (
	GRPCEndpoint = "/var/run/vdpa/vdpa.sock"
)

type VFPCIMapping struct {
	PciAddress string  `json:"pciAddress"`
	Socketpath string  `json:"socketpath"`
}

