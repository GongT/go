package CSI

import "strings"

// 多个CSI序列的组合
type Seq struct {
	seq      []St
	disabled bool
}

func Combine(seq ...St) Seq {
	return Seq{
		seq:      seq,
		disabled: false,
	}
}

func (s *Seq) String() string {
	if s.disabled {
		return ""
	}

	ss := strings.Builder{}
	for _, v := range s.seq {
		if ss.Len() > 0 {
			ss.WriteString(";")
		}
		s := v.Sequence()
		if s != "" {
			ss.WriteString(s)
		}
	}

	return "\x1b[" + ss.String() + "m"
}

func (s *Seq) Add(seq ...St) *Seq {
	s.seq = append(s.seq, seq...)
	return s
}

func (s *Seq) Delete(seq ...St) *Seq {
	for _, del := range seq {
		for i, v := range s.seq {
			if v == del {
				s.seq = append(s.seq[:i], s.seq[i+1:]...)
				break
			}
		}
	}
	return s
}

func (s *Seq) Disable() *Seq {
	s.disabled = true
	return s
}

func (s *Seq) Enable() *Seq {
	s.disabled = false
	return s
}

func (e *Seq) ToReset() *Seq {
	var bits St = 0
	for _, v := range e.seq {
		bits |= v
	}

	reset := bits.noColorBits().ToReset()

	// 原序列中包含前景色，则添加前景色重置
	// FIXME: 这里可以用一个位操作一次性解决两者，不用俩if
	if bits&Fore != 0 {
		reset |= Fore
	}
	if bits&Back != 0 {
		reset |= Back
	}

	e.seq = []St{reset}

	return e
}
