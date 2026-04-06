package pathplanner

import (
	"fmt"
	"strings"

	"github.com/LikeMyGames/FRC-Like-Robot/state/utils/curves"
)

var (
	paths map[string]*Path = make(map[string]*Path)
	autos map[string]*Auto = make(map[string]*Auto)
)

func init() {
	loadAllPathPlannerPaths()
	loadAllPathPlannerAutos()
}

func PathPlannerPath(name string) *Path {
	return paths[name]
}

func PathPlannerAuto(name string) *Auto {
	return autos[name]
}

func RefreshPathPlannerPathCurve(name string) {
	paths[name].RefreshCurve()
}

func (p *Path) RefreshCurve() {
	for i := range len(p.Waypoints) - 2 {
		start := p.Waypoints[i]
		end := p.Waypoints[i+1]
		curve := curves.NewCubicBezier(
			start.Anchor,
			start.NextControl,
			end.PrevControl,
			end.Anchor,
		)
		p.bezier_curves = append(p.bezier_curves, curve)
	}
}

func (p Constraints) String() string {
	s := new(strings.Builder)

	fmt.Fprintf(s, "Velocity: %v\n", p.Velocity)
	fmt.Fprintf(s, "Acceleration: %v\n", p.Acceleration)
	fmt.Fprintf(s, "Angular Velocity: %v\n", p.AngularVelocity)
	fmt.Fprintf(s, "Angular Acceleration: %v\n", p.AngularAcceleration)
	fmt.Fprintf(s, "Acceleration: %v\n", p.Acceleration)
	fmt.Fprintf(s, "Nominal Voltage: %v\n", p.NominalVoltage)
	fmt.Fprintf(s, "Unlimited: %v", p.Unlimited)

	return s.String()
}

func (c *Command) String() string {
	s := new(strings.Builder)

	fmt.Fprintf(s, "Type: %s", c.Type)

	return s.String()
}

func (p *Path) String() string {
	s := new(strings.Builder)

	fmt.Fprintf(s, "Path (%s): %s\n", p.Version, p.Name)
	fmt.Fprintf(s, "\tWaypoints: %v\n", len(p.Waypoints))
	fmt.Fprintf(s, "\tRotation Targets: %v\n", len(p.RotationTargets))
	fmt.Fprintf(s, "\tEvent Markers: %v\n", len(p.EventMarkers))
	fmt.Fprintf(s, "\tPoint Towards Zones: %v\n", len(p.PointTowardsZones))
	fmt.Fprintf(s, "\tConstraint Zones: %v\n", len(p.ConstraintZones))
	fmt.Fprintf(s, "\tGlobal Constraints: \n%s\n", p.GlobalConstraints)

	return s.String()
}

func (a *Auto) String() string {
	s := new(strings.Builder)
	fmt.Fprintf(s, "Auto (%s): %s\n", a.Version, a.Name)
	fmt.Fprintln(s, a.Command)
	fmt.Fprintf(s, "\tReset Odometry: %v\n", a.ResetOdometry)
	fmt.Fprintf(s, "\tChoreo Auto: %v", a.ChoreoAuto)

	return s.String()
}
