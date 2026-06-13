package application

type TreeResponse struct {
	Menus []MenuResponse `json:"menus"`
}

type MenuResponse struct {
	ID        uint            `json:"ID"`
	ParentID  uint            `json:"parentId"`
	Path      string          `json:"path"`
	Name      string          `json:"name"`
	Hidden    bool            `json:"hidden"`
	Component string          `json:"component"`
	Sort      int             `json:"sort"`
	Meta      MenuMeta        `json:"meta"`
	Children  []MenuResponse  `json:"children"`
	Btns      []any           `json:"btns"`
	MenuBtn   []any           `json:"menuBtn"`
	Parent    string          `json:"parent,omitempty"`
}

type MenuMeta struct {
	Title      string `json:"title"`
	Icon       string `json:"icon"`
	KeepAlive  bool   `json:"keepAlive"`
	ActiveName string `json:"activeName"`
	DefaultMenu bool  `json:"defaultMenu"`
	CloseTab   bool   `json:"closeTab"`
	Hidden     bool   `json:"hidden"`
}

type SaveMenuRequest struct {
	ID             uint                 `json:"id"`
	ParentID       uint                 `json:"parentId"`
	Path           string               `json:"path"`
	Name           string               `json:"name"`
	Hidden         bool                 `json:"hidden"`
	Component      string               `json:"component"`
	Sort           int                  `json:"sort"`
	Title          string               `json:"title"`
	Icon           string               `json:"icon"`
	ActiveName     string               `json:"activeName"`
	KeepAlive      bool                 `json:"keepAlive"`
	DefaultMenu    bool                 `json:"defaultMenu"`
	CloseTab       bool                 `json:"closeTab"`
	TransitionType string               `json:"transitionType"`
	Parameters     []MenuParameterInput `json:"parameters"`
	Buttons        []MenuButtonInput    `json:"menuBtn"`
}

type MenuParameterInput struct {
	Type  string `json:"type"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MenuButtonInput struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type MenuDetailResponse struct {
	Menu        MenuResponse          `json:"menu"`
	Meta        MenuMetaResponse      `json:"meta"`
	Parameters  []MenuParameterResponse `json:"parameters"`
	Buttons     []MenuButtonResponse  `json:"menuBtn"`
}

type MenuMetaResponse struct {
	ActiveName     string `json:"activeName"`
	KeepAlive      bool   `json:"keepAlive"`
	DefaultMenu    bool   `json:"defaultMenu"`
	CloseTab       bool   `json:"closeTab"`
	TransitionType string `json:"transitionType"`
}

type MenuParameterResponse struct {
	Type  string `json:"type"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MenuButtonResponse struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type SaveMenuResponse struct {
	ID uint `json:"id"`
}
