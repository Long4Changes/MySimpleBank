package api

import (
	"github.com/Long4Changes/MySimpleBank/db/util"
	"github.com/go-playground/validator/v10"
)

// 问题引入：
// Currency string `json:"currency" binding:"required,oneof=USD EUR CAD"`
// 如果我们以后要引入更多货币种类怎么办，难道全部写在 oneof=后面吗？
// 而且这一行代码可能在很多 API 都会用到，具有重复性

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	// fieldLevel.Field() returns reflect.Value 反射对象
	// 通过 .Interface() 把 reflect.Value 还原成 any 类型
	// 通过 .(string) 进行类型断言，把 any 转成 string
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// check currency is supported
		return util.IsSupportedCurrency(currency)
	}
	return false
}