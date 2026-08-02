package validator

import "github.com/chengkz2023/My-GVA/server/utils"

type Rules = utils.Rules

func Verify(value any, rules Rules) error {
	return utils.Verify(value, rules)
}
