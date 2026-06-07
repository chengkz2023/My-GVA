package domain

type Menu struct {
	ID        uint
	ParentID  uint
	Path      string
	Name      string
	Hidden    bool
	Component string
	Sort      int
	Title     string
	Icon      string
	Children  []Menu
}
