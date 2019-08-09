package journaltypes

// trancactional contains all the journal types that influence credits in some way

type BaseTransaction struct {
	Base

	// General buys
	Cost      int64
	TotalCost int64

	// General sales
	TotalSale int64

	// Exploration data
	BaseValue, Bonus int64

	// Missions & Community Goals
	Reward   int64
	Donation int64

	// Moving stuff around
	TransferCost int64

	// Modules etc.
	BuyPrice  int64
	SellPrice int64

	// Fines / Vouchers / PowerplaySalary - can be both positive and negative. Thanks, Frontier :(
	Amount int64

	// Ship buying & selling
	ShipPrice int64 // Can be both buy & sale
}

// ComputeDelta computes the transactional delta for a BaseTransaction
func (t BaseTransaction) ComputeDelta() int64 {
	var total int64

	total = total + t.TotalSale
	total = total + t.BaseValue
	total = total + t.Bonus
	total = total + t.Reward
	total = total + t.SellPrice

	total = total - t.Cost
	total = total - t.TotalCost
	total = total - t.Donation
	total = total - t.TransferCost
	total = total - t.BuyPrice

	// Handle ShipPrice
	shipPriceSign := int64(1)
	if t.Event == "ShipyeardBuy" {
		shipPriceSign = -1
	}
	total = total + (shipPriceSign * t.ShipPrice)

	// handle Amount - thanks again Frontier...
	amountSign := int64(1)
	switch t.Event {
	case "PayFines":
		amountSign = -1
		break
	case "PayLegacyFines":
		amountSign = -1
		break
	}
	total = total + (amountSign * t.Amount)

	return total
}

type BuyExplorationData struct {
	BaseTransaction
}

type SellExplorationData struct {
	BaseTransaction
}

type BuyTradeData struct {
	BaseTransaction
}

type MarketBuy struct {
	BaseTransaction
}

type MarketSale struct {
	BaseTransaction
}

type BuyAmmo struct {
	BaseTransaction
}

type BuyDrones struct {
	BaseTransaction
}

type CommunityGoalReward struct {
	BaseTransaction
}

type CrewHire struct {
	BaseTransaction
}

type FetchRemoteModule struct {
	BaseTransaction
}

type MissionCompleted struct {
	BaseTransaction
}

type ModuleBuy struct {
	BaseTransaction
}

type ModuleSell struct {
	BaseTransaction
}

type ModuleSellRemote struct {
	BaseTransaction
}

type PayFines struct {
	BaseTransaction
}

type PayLegacyFines struct {
	BaseTransaction
}

type RedeemVoucher struct {
	BaseTransaction
}

type RefuelAll struct {
	BaseTransaction
}

type RefuelPartial struct {
	BaseTransaction
}

type Repair struct {
	BaseTransaction
}

type RepairAll struct {
	BaseTransaction
}

type RestockVehicle struct {
	BaseTransaction
}

type SellDrones struct {
	BaseTransaction
}

type ShipyardBuy struct {
	BaseTransaction

	ShipType string
}

type ShipyardSell struct {
	BaseTransaction
}

type ShipyardTransfer struct {
	BaseTransaction
}

type PowerplayFastTrack struct {
	BaseTransaction
}

type PowerplaySalary struct {
	BaseTransaction
}
