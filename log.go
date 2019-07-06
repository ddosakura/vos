package vos

// LogMode for vos/session
type LogMode uint8

// LogMode
const (
	LmN LogMode = iota // No Log
	LmE
	LmW
	LmI
	LmD
)

// --- VOS ---

// Log Action
func (v *OS) Log(prompt string, color func(a ...interface{}) string, vs ...interface{}) {
	vss := make([]interface{}, 0, len(vs)+1)
	vss = append(vss, color("["+prompt+"]"))
	vss = append(vss, vs...)
	v.Logger.Println(vss...)
}

// Debug Log
func (v *OS) Debug(vs ...interface{}) {
	if v.LogMode == LmD {
		v.Log("D", Green, vs...)
	}
}

// Info Log
func (v *OS) Info(vs ...interface{}) {
	if v.LogMode >= LmI {
		v.Log("I", Blue, vs...)
	}
}

// Warn Log
func (v *OS) Warn(vs ...interface{}) {
	if v.LogMode >= LmW {
		v.Log("W", Yellow, vs...)
	}
}

// Error Log
func (v *OS) Error(e error) {
	if v.LogMode >= LmE {
		v.Log("E", Red, e)
	}
	panic(e)
}

// --- Session ---

// Log Action
func (v *Session) Log(prompt string, color func(a ...interface{}) string, vs ...interface{}) {
	vss := make([]interface{}, 0, len(vs)+1)
	vss = append(vss, color("["+prompt+"]"))
	vss = append(vss, vs...)
	v.Logger.Println(vss...)
}

// Debug Log
func (v *Session) Debug(vs ...interface{}) {
	if v.LogMode == LmD {
		v.Log("D", Green, vs...)
	}
}

// Info Log
func (v *Session) Info(vs ...interface{}) {
	if v.LogMode >= LmI {
		v.Log("I", Blue, vs...)
	}
}

// Warn Log
func (v *Session) Warn(vs ...interface{}) {
	if v.LogMode >= LmW {
		v.Log("W", Yellow, vs...)
	}
}

// Error Log
func (v *Session) Error(e error) {
	if v.LogMode >= LmE {
		v.Log("E", Red, e)
	}
	panic(e)
}
