package sim

// LightCycle models a fixed-duration traffic light FSM.
type LightCycle struct {
	state    string // "green", "yellow", "red"
	elapsed  int    // seconds elapsed in current state
}

const (
	greenDuration  = 30
	yellowDuration = 5
	redDuration    = 25
)

func NewLightCycle() *LightCycle { return &LightCycle{state: "green", elapsed: 0} }

// Tick advances the timer by one second and transitions state when needed.
func (l *LightCycle) Tick() {
	switch l.state {
	case "green":
		l.elapsed++
		if l.elapsed >= greenDuration {
			l.state = "yellow"
			l.elapsed = 0
		}
	case "yellow":
		l.elapsed++
		if l.elapsed >= yellowDuration {
			l.state = "red"
			l.elapsed = 0
		}
	case "red":
		l.elapsed++
		if l.elapsed >= redDuration {
			l.state = "green"
			l.elapsed = 0
		}
	}
}

func (l *LightCycle) State() string { return l.state }
func (l *LightCycle) Elapsed() int  { return l.elapsed }
