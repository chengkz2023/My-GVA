package validator

import "github.com/flipped-aurora/gin-vue-admin/server/utils"

type Rules = utils.Rules

func Verify(value any, rules Rules) error {
	return utils.Verify(value, rules)
}
