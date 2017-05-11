package cmd

import (
	"sort"

	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/ssh"
)

func init() {
	com.Register("cmd", &Component{})
}

type Component struct{}

type Preprocessor interface {
	PreprocessOrder() uint
	PreprocessSession(sess ssh.Session) (cont bool, msg string)
}

func Preprocessors() []Preprocessor {
	var processors []Preprocessor
	for _, com := range com.Enabled(new(Preprocessor), nil) {
		processors = append(processors, com.(Preprocessor))
	}
	sort.Slice(processors, func(i, j int) bool {
		return processors[i].PreprocessOrder() <= processors[j].PreprocessOrder()
	})
	return processors
}
