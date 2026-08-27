package fpath

import "fmt"

func (ip *IPath) Format(stat fmt.State, verb rune) {
	do_fmt(stat, verb, "ImmutablePath", ip.value.Raw())
}

func (p *Path) Format(stat fmt.State, verb rune) {
	do_fmt(stat, verb, "Path", p.Raw())
}

func do_fmt(stat fmt.State, verb rune, name string, value string) {
	full := stat.Flag('+')

	if full {
		fmt.Fprintf(stat, "%s<%s>", name, value)
	} else {
		relative, err := ToCwdRelative(value)
		if err != nil {
			fmt.Fprintf(stat, "%s<%s>", name, value)
		} else {
			fmt.Fprintf(stat, "%s<%s>", name, relative)
		}
	}
}

func (ip *IPath) GoString() string {
	return fmt.Sprintf("&fpath.IPath{%q}", ip.value.Raw())
}

func (p *Path) GoString() string {
	return fmt.Sprintf("&fpath.Path{%q}", p.Raw())
}
