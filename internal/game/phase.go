package game

import "time"

type Phase int

const (
	PhaseBuild Phase = iota
	PhaseBattle
	PhaseRebuild
	PhaseVictory
	PhaseDefeat
)

func (p Phase) String() string {
	switch p {
	case PhaseBuild:
		return "BUILD"
	case PhaseBattle:
		return "BATTLE"
	case PhaseRebuild:
		return "REBUILD"
	case PhaseVictory:
		return "VICTORY"
	case PhaseDefeat:
		return "DEFEAT"
	default:
		return "UNKNOWN"
	}
}

func phaseDuration(phase Phase) time.Duration {
	switch phase {
	case PhaseBuild:
		return 10 * time.Second

	case PhaseBattle:
		return 10 * time.Second

	case PhaseRebuild:
		return 5 * time.Second

	default:
		return 0
	}
}
