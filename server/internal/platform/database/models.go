package database

import (
	"time"

	"github.com/google/uuid"
)

// ========== User types ==========

type Login interface {
	GetUsername() string
	GetNickname() string
	GetUUID() uuid.UUID
	GetUserId() uint
	GetAuthorityId() uint
	GetUserInfo() any
}

var _ Login = new(SysUser)

type SysUser struct {
	GVA_MODEL
	UUID          uuid.UUID      `json:"uuid" gorm:"index;comment:用户UUID"`
	Username      string         `json:"userName" gorm:"index;comment:用户登录名"`
	Password      string         `json:"-"  gorm:"comment:用户登录密码"`
	NickName      string         `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`
	HeaderImg     string         `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"`
	AuthorityId   uint           `json:"authorityId" gorm:"default:888;comment:用户角色ID"`
	Authority     SysAuthority   `json:"authority" gorm:"foreignKey:AuthorityId;references:AuthorityId;comment:用户角色"`
	Authorities   []SysAuthority `json:"authorities" gorm:"many2many:sys_user_authority;"`
	Phone         string         `json:"phone"  gorm:"comment:用户手机号"`
	Email         string         `json:"email"  gorm:"comment:用户邮箱"`
	Enable        int            `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"`
	OriginSetting JSONMap `json:"originSetting" form:"originSetting" gorm:"type:text;default:null;column:origin_setting;comment:配置;"`
}

func (SysUser) TableName() string {
	return "sys_users"
}

func (s *SysUser) GetUsername() string {
	return s.Username
}

func (s *SysUser) GetNickname() string {
	return s.NickName
}

func (s *SysUser) GetUUID() uuid.UUID {
	return s.UUID
}

func (s *SysUser) GetUserId() uint {
	return s.ID
}

func (s *SysUser) GetAuthorityId() uint {
	return s.AuthorityId
}

func (s *SysUser) GetUserInfo() any {
	return *s
}

// SysUserAuthority is the join table for sysUser and sysAuthority
type SysUserAuthority struct {
	SysUserId               uint `gorm:"column:sys_user_id"`
	SysAuthorityAuthorityId uint `gorm:"column:sys_authority_authority_id"`
}

func (s *SysUserAuthority) TableName() string {
	return "sys_user_authority"
}

// ========== Role types ==========

type SysAuthority struct {
	CreatedAt       time.Time       // 创建时间
	UpdatedAt       time.Time       // 更新时间
	DeletedAt       *time.Time      `sql:"index"`
	AuthorityId     uint            `json:"authorityId" gorm:"not null;unique;primary_key;comment:角色ID;size:90"`
	AuthorityName   string          `json:"authorityName" gorm:"comment:角色名"`
	ParentId        *uint           `json:"parentId" gorm:"comment:父角色ID"`
	DataAuthorityId []*SysAuthority `json:"dataAuthorityId" gorm:"many2many:sys_data_authority_id;"`
	Children        []SysAuthority  `json:"children" gorm:"-"`
	SysBaseMenus    []SysBaseMenu   `json:"menus" gorm:"many2many:sys_authority_menus;"`
	Users           []SysUser       `json:"-" gorm:"many2many:sys_user_authority;"`
	DefaultRouter   string          `json:"defaultRouter" gorm:"comment:默认菜单;default:dashboard"`
}

func (SysAuthority) TableName() string {
	return "sys_authorities"
}

type SysAuthorityBtn struct {
	AuthorityId      uint           `gorm:"comment:角色ID"`
	SysMenuID        uint           `gorm:"comment:菜单ID"`
	SysBaseMenuBtnID uint           `gorm:"comment:菜单按钮ID"`
	SysBaseMenuBtn   SysBaseMenuBtn ` gorm:"comment:按钮详情"`
}

type SysAuthorityMenu struct {
	MenuId      string `json:"menuId" gorm:"comment:菜单ID;column:sys_base_menu_id"`
	AuthorityId string `json:"-" gorm:"comment:角色ID;column:sys_authority_authority_id"`
}

func (s SysAuthorityMenu) TableName() string {
	return "sys_authority_menus"
}

// ========== Menu types ==========

type SysMenu struct {
	SysBaseMenu
	MenuId      uint                   `json:"menuId" gorm:"comment:菜单ID"`
	AuthorityId uint                   `json:"-" gorm:"comment:角色ID"`
	Children    []SysMenu              `json:"children" gorm:"-"`
	Parameters  []SysBaseMenuParameter `json:"parameters" gorm:"foreignKey:SysBaseMenuID;references:MenuId"`
	Btns        map[string]uint        `json:"btns" gorm:"-"`
}

type SysBaseMenu struct {
	GVA_MODEL
	MenuLevel     uint                   `json:"-"`
	ParentId      uint                   `json:"parentId" gorm:"comment:父菜单ID"`
	Path          string                 `json:"path" gorm:"comment:路由path"`
	Name          string                 `json:"name" gorm:"comment:路由name"`
	Hidden        bool                   `json:"hidden" gorm:"comment:是否在列表隐藏"`
	Component     string                 `json:"component" gorm:"comment:对应前端文件路径"`
	Sort          int                    `json:"sort" gorm:"comment:排序标记"`
	Meta          `json:"meta" gorm:"embedded"`
	SysAuthoritys []SysAuthority         `json:"authoritys" gorm:"many2many:sys_authority_menus;"`
	Children      []SysBaseMenu          `json:"children" gorm:"-"`
	Parameters    []SysBaseMenuParameter `json:"parameters"`
	MenuBtn       []SysBaseMenuBtn       `json:"menuBtn"`
}

type Meta struct {
	ActiveName     string `json:"activeName" gorm:"comment:高亮菜单"`
	KeepAlive      bool   `json:"keepAlive" gorm:"comment:是否缓存"`
	DefaultMenu    bool   `json:"defaultMenu" gorm:"comment:是否是基础路由（开发中）"`
	Title          string `json:"title" gorm:"comment:菜单名"`
	Icon           string `json:"icon" gorm:"comment:菜单图标"`
	CloseTab       bool   `json:"closeTab" gorm:"comment:自动关闭tab"`
	TransitionType string `json:"transitionType" gorm:"comment:路由切换动画"`
}

type SysBaseMenuParameter struct {
	GVA_MODEL
	SysBaseMenuID uint
	Type          string `json:"type" gorm:"comment:地址栏携带参数为params还是query"`
	Key           string `json:"key" gorm:"comment:地址栏携带参数的key"`
	Value         string `json:"value" gorm:"comment:地址栏携带参数的值"`
}

func (SysBaseMenu) TableName() string {
	return "sys_base_menus"
}

type SysBaseMenuBtn struct {
	GVA_MODEL
	Name          string `json:"name" gorm:"comment:按钮关键key"`
	Desc          string `json:"desc" gorm:"按钮备注"`
	SysBaseMenuID uint   `json:"sysBaseMenuID" gorm:"comment:菜单ID"`
}

// ========== Other types (AutoMigrate-only) ==========

type SysDictionary struct {
	GVA_MODEL
	Name                 string                `json:"name" form:"name" gorm:"column:name;comment:字典名（中）"`
	Type                 string                `json:"type" form:"type" gorm:"column:type;comment:字典名（英）"`
	Status               *bool                 `json:"status" form:"status" gorm:"column:status;comment:状态"`
	Desc                 string                `json:"desc" form:"desc" gorm:"column:desc;comment:描述"`
	ParentID             *uint                 `json:"parentID" form:"parentID" gorm:"column:parent_id;comment:父级字典ID"`
	Children             []SysDictionary       `json:"children" gorm:"foreignKey:ParentID"`
	SysDictionaryDetails []SysDictionaryDetail `json:"sysDictionaryDetails" form:"sysDictionaryDetails"`
}

func (SysDictionary) TableName() string {
	return "sys_dictionaries"
}

type SysDictionaryDetail struct {
	GVA_MODEL
	Label           string                `json:"label" form:"label" gorm:"column:label;comment:展示值"`
	Value           string                `json:"value" form:"value" gorm:"column:value;comment:字典值"`
	Extend          string                `json:"extend" form:"extend" gorm:"column:extend;comment:扩展值"`
	Status          *bool                 `json:"status" form:"status" gorm:"column:status;comment:启用状态"`
	Sort            int                   `json:"sort" form:"sort" gorm:"column:sort;comment:排序标记"`
	SysDictionaryID int                   `json:"sysDictionaryID" form:"sysDictionaryID" gorm:"column:sys_dictionary_id;comment:关联标记"`
	ParentID        *uint                 `json:"parentID" form:"parentID" gorm:"column:parent_id;comment:父级字典详情ID"`
	Children        []SysDictionaryDetail `json:"children" gorm:"foreignKey:ParentID"`
	Level           int                   `json:"level" form:"level" gorm:"column:level;comment:层级深度"`
	Path            string                `json:"path" form:"path" gorm:"column:path;comment:层级路径"`
	Disabled        bool                  `json:"disabled" gorm:"-"`
}

func (SysDictionaryDetail) TableName() string {
	return "sys_dictionary_details"
}

type SysError struct {
	GVA_MODEL
	Form     *string `json:"form" form:"form" gorm:"comment:错误来源;column:form;type:text;" binding:"required"`
	Info     *string `json:"info" form:"info" gorm:"comment:错误内容;column:info;type:text;"`
	Level    string  `json:"level" form:"level" gorm:"comment:日志等级;column:level;"`
	Solution *string `json:"solution" form:"solution" gorm:"comment:解决方案;column:solution;type:text"`
	Status   string  `json:"status" form:"status" gorm:"comment:处理状态;column:status;type:varchar(20);default:未处理;"`
}

func (SysError) TableName() string {
	return "sys_error"
}

type JwtBlacklist struct {
	GVA_MODEL
	Jwt string `gorm:"type:text;comment:jwt"`
}

type SysParams struct {
	GVA_MODEL
	Name  string `json:"name" form:"name" gorm:"column:name;comment:参数名称;" binding:"required"`
	Key   string `json:"key" form:"key" gorm:"column:key;comment:参数键;" binding:"required"`
	Value string `json:"value" form:"value" gorm:"column:value;comment:参数值;" binding:"required"`
	Desc  string `json:"desc" form:"desc" gorm:"column:desc;comment:参数说明;"`
}

func (SysParams) TableName() string {
	return "sys_params"
}

type SysVersion struct {
	GVA_MODEL
	VersionName *string `json:"versionName" form:"versionName" gorm:"comment:版本名称;column:version_name;size:255;" binding:"required"`
	VersionCode *string `json:"versionCode" form:"versionCode" gorm:"comment:版本号;column:version_code;size:100;" binding:"required"`
	Description *string `json:"description" form:"description" gorm:"comment:版本描述;column:description;size:500;"`
	VersionData *string `json:"versionData" form:"versionData" gorm:"comment:版本数据JSON;column:version_data;type:text;"`
}

func (SysVersion) TableName() string {
	return "sys_versions"
}
