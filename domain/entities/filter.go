package entities

import "time"

// GeneralFilter - Used to receive the filters to list something
type GeneralFilter struct {
	//Column - Column of struct, like id, isActive, name... to know witch column will be ordinate
	Column string

	//OrdinationAsc - true if ordination is asc else ordination is desc
	OrdinationAsc bool

	//Limit - limit of registers that will show on screen
	Limit int64

	//Page - Current page showing
	Page int64

	//Search - Name of item to filter
	Search string

	//Type - Type of some item
	Type int64

	//IsInactive - true if inactive else false
	IsInactive bool

	//Status situation number
	Status int64
}

// FilterUserAccess - Used to receive the filters to modal of access
type FilterUserAccess struct {
	//GroupBy - contains for what will be grouped
	GroupBy string

	//StartDate - contains the start date to filter
	StartDate time.Time

	//EndDate - contains the end date to filter
	EndDate time.Time
}

type Ratings struct {
	Index      int64
	IndexEmpty bool
}
