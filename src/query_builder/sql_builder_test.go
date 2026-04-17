package querybuilder

import (
	"fmt"
	"testing"

	"github.com/bobllor/assert"
	"github.com/bobllor/cloud-project/src/session"
	"github.com/bobllor/cloud-project/src/user"
)

func TestSelectBuild(t *testing.T) {
	sb := NewSqlBuilder(user.TableName)

	query := sb.Select().
		Columns(user.ColumnAccountID, user.ColumnActive, user.ColumnCreatedOn).
		Build()

	fmt.Println(query)
}

func TestSelectWhereBuild(t *testing.T) {
	t.Run("Basic Chaining", func(t *testing.T) {
		sb := NewSqlBuilder(user.TableName)
		str := sb.Select().Columns(user.ColumnUsername).Where().
			Equal(user.ColumnAccountID, "12345").And().
			Equal(user.ColumnActive, true).Or().
			Equal(user.ColumnCreatedOn, 12345).Build()

		baseStr := fmt.Sprintf(
			"SELECT %s FROM %s WHERE %s = (?) AND %s = (?) OR %s = (?)",
			user.ColumnUsername,
			user.TableName,
			user.ColumnAccountID,
			user.ColumnActive,
			user.ColumnCreatedOn,
		)

		assert.Equal(t, str, baseStr)
	})

	t.Run("In Operator", func(t *testing.T) {
		sb := NewSqlBuilder(user.TableName)
		str := sb.Select().Columns(user.ColumnAccountID).Where().
			In(user.ColumnAccountID, "12345", "5678", "1111").Build()

		baseStr := fmt.Sprintf(
			"SELECT %s FROM %s WHERE %s IN (?,?,?)",
			user.ColumnAccountID,
			user.TableName,
			user.ColumnAccountID,
		)

		assert.Equal(t, str, baseStr)
	})

	t.Run("Exists Operator", func(t *testing.T) {
		mainSb := NewSqlBuilder(user.TableName)
		subSb := NewSqlBuilder(session.TableName)

		subStr := subSb.Select().Columns(session.ColumnCreatedOn).Where().
			Equal(session.ColumnAccountID, "15555").Build()
		mainStr := mainSb.Select().Columns(user.ColumnAccountID).Where().
			Exists(subStr, subSb.Args()...).Build()

		fmt.Println(mainStr, mainSb.Args())
	})
}
