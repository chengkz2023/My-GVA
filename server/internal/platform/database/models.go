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

// ========== Operation record types ==========

type SysOperationRecord struct {
	GVA_MODEL
	Ip           string        `json:"ip" form:"ip" gorm:"column:ip;comment:请求ip"`
	Method       string        `json:"method" form:"method" gorm:"column:method;comment:请求方法"`
	Path         string        `json:"path" form:"path" gorm:"column:path;comment:请求路径"`
	Status       int           `json:"status" form:"status" gorm:"column:status;comment:请求状态"`
	Latency      time.Duration `json:"latency" form:"latency" gorm:"column:latency;comment:延迟" swaggertype:"string"`
	Agent        string        `json:"agent" form:"agent" gorm:"type:text;column:agent;comment:代理"`
	ErrorMessage string        `json:"error_message" form:"error_message" gorm:"column:error_message;comment:错误信息"`
	Body         string        `json:"body" form:"body" gorm:"type:text;column:body;comment:请求Body"`
	Resp         string        `json:"resp" form:"resp" gorm:"type:text;column:resp;comment:响应Body"`
	UserID       int           `json:"user_id" form:"user_id" gorm:"column:user_id;comment:用户id"`
	User         SysUser       `json:"user"`
}

func (SysOperationRecord) TableName() string { return "sys_operation_records" }
