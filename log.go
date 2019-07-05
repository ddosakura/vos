package vos

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
	if v.DebugMode {
		v.Log("D", Green, vs...)
	}
}

// Info Log
func (v *OS) Info(vs ...interface{}) {
	v.Log("I", Blue, vs...)
}

// Warn Log
func (v *OS) Warn(vs ...interface{}) {
	v.Log("W", Yellow, vs...)
}

// Error Log
func (v *OS) Error(e error) {
	v.Log("E", Red, e)
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
	if v.DebugMode {
		v.Log("D", Green, vs...)
	}
}

// Info Log
func (v *Session) Info(vs ...interface{}) {
	v.Log("I", Blue, vs...)
}

// Warn Log
func (v *Session) Warn(vs ...interface{}) {
	v.Log("W", Yellow, vs...)
}

// Error Log
func (v *Session) Error(e error) {
	v.Log("E", Red, e)
	panic(e)
}
