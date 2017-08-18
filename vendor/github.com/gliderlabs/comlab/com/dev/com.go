package dev

import "github.com/gliderlabs/comlab/pkg/com"

func init() {
	com.Register("dev", &Component{})
}

type Component struct{}
