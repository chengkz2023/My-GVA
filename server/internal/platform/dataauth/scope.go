// Package dataauth 提供行级数据权限的通用构造器。
//
// 模式：角色通过数据权限映射（sys_data_authority_id，见角色模块 SetDataAuthority）持有
// 若干「可见部门/组织 ID」。业务查询以这些 ID 为范围过滤数据行。本包只负责把 ID 集合
// 转成 SQL 过滤条件；「从当前登录角色解析 ID 集合」由业务模块自行注入提供者完成，
// 使业务模块不必反向依赖角色模块。
package dataauth

import "gorm.io/gorm"

// Scope 返回应用了行级过滤的查询：WHERE <column> IN (<ids>)。
//
// 约定：
//   - ids 为空表示「无可见范围」，返回的查询将匹配不到任何行（1 = 0）；
//     若业务语义是「空 = 全部可见」，请在调用前自行跳过。
//   - column 必须是可信的表列名（代码内常量），严禁拼接用户输入。
func Scope(db *gorm.DB, column string, ids []uint) *gorm.DB {
	if len(ids) == 0 {
		return db.Where("1 = 0")
	}
	values := make([]any, len(ids))
	for i, id := range ids {
		values[i] = id
	}
	return db.Where(column+" IN ?", values)
}
