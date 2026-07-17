package game

import "time"

func (g *Game) updateRound(deltaTime time.Duration) {
	if g.phase == PhaseVictory || g.phase == PhaseDefeat {
		return
	}

	g.phaseTimeLeft -= deltaTime

	if g.phaseTimeLeft > 0 {
		return
	}

	g.advancePhase()
}

func (g *Game) advancePhase() {
	switch g.phase {
	case PhaseBuild:
		if !g.grid.IsCastleEnclosed() {
			g.startPhase(PhaseDefeat)
			return
		}

		g.startPhase(PhaseBattle)

	case PhaseBattle:
		g.roundNumber++
		g.startPhase(PhaseRebuild)

	case PhaseRebuild:
		if !g.grid.IsCastleEnclosed() {
			g.startPhase(PhaseDefeat)
			return
		}

		g.startPhase(PhaseBattle)
	}
}

func (g *Game) startPhase(phase Phase) {
	g.phase = phase
	g.phaseTimeLeft = phaseDuration(phase)
}

func (g *Game) canBuild() bool {
	return g.phase == PhaseBuild || g.phase == PhaseRebuild
}
