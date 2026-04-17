package nodes

import "OnlineGame/config"

type ArmyID uint64

type Army struct {
	ID            ArmyID     `json:"id"`
	Pos           uint       `json:"pos"`
	NodeEdgeID    NodeEdgeID `json:"node_edge_id"`
	HeadingFromID NodeID     `json:"heading_from_id"`
	HeadingToID   NodeID     `json:"heading_to_id"`
	OwnerPlayerID uint       `json:"owner_player_id"`
	Value         uint       `json:"value"`
}

func (a *Army) Tick(world WorldInterface) {
	nodeEdge := world.GetNodeEdgeByID(a.NodeEdgeID)
	if nodeEdge == nil {
		panic("node edge not found")
	}
	headingFrom := world.GetNodeByID(a.HeadingFromID)
	if headingFrom == nil {
		panic("heading from id not found")
	}
	headingTo := world.GetNodeByID(a.HeadingToID)
	if headingTo == nil {
		panic("heading to id not found")
	}

	if a.Value == 0 {
		return
	}
	a.Pos += uint(float64(config.Game().ArmySpeed) * world.GetDelta())

	// Collision armies
	for _, otherArmy := range nodeEdge.Armies {
		if a.OwnerPlayerID == otherArmy.OwnerPlayerID {
			continue // Skip army collisions for 1 player armies
		}
		otherArmyHeadingFrom := world.GetNodeByID(otherArmy.HeadingFromID)
		if headingTo.ID == otherArmyHeadingFrom.ID && a.ID != otherArmy.ID {
			distance := (nodeEdge.Length - otherArmy.Pos) - a.Pos
			if distance < 0 && otherArmy.Value != 0 {
				if a.Value >= otherArmy.Value {
					a.Value -= otherArmy.Value
					otherArmy.Value = 0
				} else {
					otherArmy.Value -= a.Value
					a.Value = 0
				}
			}
		}
	}

	// Army reached HeadingToID node
	isReachedGoal := a.Pos >= nodeEdge.Length

	if isReachedGoal && a.Value != 0 {
		if a.OwnerPlayerID != headingTo.OwnerID { // If army reached other owner node
			// Shield logic
			if headingTo.Shield > a.Value {
				headingTo.Shield -= a.Value
			} else {
				a.Value -= headingTo.Shield
				headingTo.Shield = 0
			}
			// Node value logic
			if a.Value <= headingTo.Value {
				headingTo.Value -= a.Value
				a.Value = 0
			} else {
				remainingValue := a.Value - headingTo.Value
				headingTo.Value = remainingValue
				headingTo.OwnerID = headingFrom.OwnerID
				headingTo.ResetTicks()
				headingTo.IsAlwaysSendArmy = false
				a.Value = 0
			}
		} else { // If army reached same owner node
			headingTo.Value += a.Value
			a.Value = 0
		}
	}
}
