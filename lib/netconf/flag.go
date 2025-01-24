package netconf

import "github.com/spf13/pflag"

// BindFlag binds the network identifier flag.
func BindFlag(flags *pflag.FlagSet, network *ID) {
	flags.StringVar((*string)(network), "network", string(*network), "Story network to participate in: story, odyssey, homer or local")
}
