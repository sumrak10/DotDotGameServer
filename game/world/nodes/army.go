package nodes

type ArmyID uint64

type Army struct {
	ID          ArmyID    `json:"id"`
	Pos         float32   `json:"pos"`
	NodeEdge    *NodeEdge `json:"node_edge"`
	HeadingFrom *Node     `json:"heading_from"`
	HeadingTo   *Node     `json:"heading_to"`
	Value       uint      `json:"value"`
}

func (a *Army) Tick() {
	if a.Value == 0 {
		return
	}
	a.Pos += 0.01

	isReachedGoal := a.Pos >= a.NodeEdge.Length
	isSameOwner := a.HeadingTo.OwnerID == a.HeadingFrom.OwnerID

	// Collision armies
	if !isSameOwner {
		for _, otherArmy := range a.NodeEdge.Armies {
			if a.HeadingTo.ID == otherArmy.HeadingFrom.ID && a.ID != otherArmy.ID {
				distance := (a.NodeEdge.Length - otherArmy.Pos) - a.Pos
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
	}

	// Army reached HeadingTo node
	if isReachedGoal && a.Value != 0 {
		if !isSameOwner {
			// Shield logic
			if a.HeadingTo.Shield > a.Value {
				a.HeadingTo.Shield -= a.Value
			} else {
				a.Value -= a.HeadingTo.Shield
				a.HeadingTo.Shield = 0
			}
			// Node value logic
			if a.Value <= a.HeadingTo.Value {
				a.HeadingTo.Value -= a.Value
				a.Value = 0
			} else {
				remainingValue := a.Value - a.HeadingTo.Value
				a.HeadingTo.Value = remainingValue
				a.HeadingTo.OwnerID = a.HeadingFrom.OwnerID
				a.Value = 0
			}
		} else {
			a.HeadingTo.Value += a.Value
			a.Value = 0
		}
	}
}
