package pathplanner

import (
	"github.com/LikeMyGames/FRC-Like-Robot/state/utils/mathutils"
)

type (
	Auto struct {
		Name          string
		Version       string  `json:"version"`
		Command       Command `json:"command"`
		ResetOdometry bool    `json:"resetOdom"`
		ChoreoAuto    bool    `json:"choreoAuto"`
	}

	Path struct {
		Name                  string
		Version               string             `json:"version"`
		Waypoints             []Waypoint         `json:"waypoints"`
		RotationTargets       []RotationTarget   `json:"rotationTargets"`
		ConstraintZones       []ConstraintZone   `json:"constraintZones"`
		PointTowardsZones     []PointTowardsZone `json:"pointTowardsZones"`
		EventMarkers          []EventMarker      `json:"eventMarkers"`
		GlobalConstraints     Constraints        `json:"globalConstraints"`
		GoalEndState          State              `json:"goalEndState"`
		IdealStartingState    State              `json:"idealStartingState"`
		Reversed              bool               `json:"reversed"`
		UseDefaultConstraints bool               `json:"useDefaultConstraints"`
	}

	Waypoint struct {
		Anchor      *mathutils.Vector2D `json:"anchor"`
		PrevControl *mathutils.Vector2D `json:"prevControl"`
		NextControl *mathutils.Vector2D `json:"nextControl"`
	}

	RotationTarget struct {
		RelativePos     float64 `json:"waypointRelativePos"`
		RotationDegrees float64 `json:"rotationDegrees"` // in degrees
	}

	ConstraintZone struct {
		Name           string  `json:"name"`
		MinRelativePos float64 `json:"minWaypointRelativePos"`
		MaxRelativePos float64 `json:"maxWaypointRelativePos"`
		Constraints    Constraints
	}

	Constraints struct {
		Velocity            float64 `json:"maxVelocity"`            // in meters per second
		Acceleration        float64 `json:"maxAcceleration"`        // in meters per second per second
		AngularVelocity     float64 `json:"maxAngularVelocity"`     // in degrees per second
		AngularAcceleration float64 `json:"maxAngularAcceleration"` // in degrees per second per second
		NominalVoltage      float64 `json:"nominalVoltage"`         // in Volts
		Unlimited           bool
	}

	PointTowardsZone struct {
	}

	EventMarker struct {
		Name           string  `json:"name"`
		RelativePos    float64 `json:"waypointRelativePos"`
		EndRelativePos float64 `json:"endWaypointRelativePos"`
		Command        Command `json:"command"`
	}

	Command struct {
		Type     string     `json:"type"`     // command type
		Children []*Command `json:"commands"` // used in any groups (sequential or all parallel types) (type == "sequential" || "parallel")
		PathName *string    `json:"pathName"` // used in path following commands (type == "path")
		WaitTime *float64   `json:"waitTime"` // used in wait commands (type == "wait")
		Name     *string    `json:"name"`     // used in Named Commands (type == "named")
	}

	State struct {
		Velocity float64 `json:"velocity"` // in meters per second
		Rotation float64 `json:"rotation"` // in degrees
	}

	NamedCommand struct {
		Argument any
		Function func(any)
	}
)
