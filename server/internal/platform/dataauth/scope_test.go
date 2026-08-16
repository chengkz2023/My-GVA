package dataauth

import (
	"strings"
	"testing"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// dummyDialector 仅用于 DryRun 构建 SQL，不建立真实连接。
type dummyDialector struct{}

func (dummyDialector) Name() string { return "dummy" }
func (dummyDialector) Initialize(db *gorm.DB) error {
	callbacks.RegisterDefaultCallbacks(db, &callbacks.Config{})
	return nil
}
func (dummyDialector) Migrator(db *gorm.DB) gorm.Migrator { return nil }
func (dummyDialector) DataTypeOf(field *schema.Field) string {
	return ""
}
func (dummyDialector) DefaultValueOf(field *schema.Field) clause.Expression {
	return clause.Expr{SQL: "NULL"}
}
func (dummyDialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v interface{}) {
	writer.WriteByte('?')
}
func (dummyDialector) QuoteTo(writer clause.Writer, str string) {
	writer.WriteString(str)
}
func (dummyDialector) Explain(sql string, vars ...interface{}) string { return sql }

func dryRunDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(dummyDialector{}, &gorm.Config{DryRun: true})
	if err != nil {
		t.Fatalf("open dry-run db: %v", err)
	}
	return db
}

func TestScopeInClause(t *testing.T) {
	db := dryRunDB(t)
	tx := Scope(db, "dept_id", []uint{3, 5}).Find(&[]struct{}{})
	if tx.Error != nil {
		t.Fatalf("Find() error = %v", tx.Error)
	}
	sql := tx.Statement.SQL.String()
	if !strings.Contains(sql, "dept_id IN") {
		t.Fatalf("sql = %q, want dept_id IN clause", sql)
	}
	if len(tx.Statement.Vars) != 2 {
		t.Fatalf("vars = %v, want two bind values", tx.Statement.Vars)
	}
}

func TestScopeEmptyDeniesAll(t *testing.T) {
	db := dryRunDB(t)
	tx := Scope(db, "dept_id", nil).Find(&[]struct{}{})
	if tx.Error != nil {
		t.Fatalf("Find() error = %v", tx.Error)
	}
	sql := tx.Statement.SQL.String()
	if !strings.Contains(sql, "1 = 0") {
		t.Fatalf("sql = %q, want deny-all condition", sql)
	}
}
