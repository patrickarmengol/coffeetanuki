package model

import "slices"

type PermissionCodes []string

func (pcs PermissionCodes) Contains(c string) bool {
	return slices.Contains(pcs, c)
}
