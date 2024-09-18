package version

import (
	"ecommerce-api/internal/consts"
	"ecommerce-api/utils"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

// RenderHandler
// the handler method should always check version_method exists or not
// if that exists, it will execute it, instead the given method
func RenderHandler(ctx *gin.Context, object interface{}, method string, args ...interface{}) {
	inputs := make([]reflect.Value, len(args))

	// passing the context to the methods
	// first argument should be the ctx
	inputs = append(inputs, reflect.ValueOf(ctx))

	for _, v := range args {
		inputs = append(inputs, reflect.ValueOf(v))
	}

	// read the data from context
	systemAcceptedVs, _ := utils.GetContext[[]string](ctx, consts.ContextSystemAcceptedVersions)
	headerVersionIndex, _ := utils.GetContext[int](ctx, consts.ContextAcceptedVersionIndex)

	// loop thorugh
	for i := len(systemAcceptedVs[0:headerVersionIndex]); i >= 0; i-- {
		versionMethod := fmt.Sprintf("%s_%s", strings.ToUpper(systemAcceptedVs[i]), method)

		// check object implement the method
		// like if the method is GetUsers, and version is v1 ; it will check v1_GetUsers
		callableMethod := reflect.ValueOf(object).MethodByName(versionMethod)
		if callableMethod.IsValid() {
			// callableMethod.Call(inputs)[0].Interface()
			callableMethod.Call(inputs)
			return
		}

	}

	// check objConv implement the method
	callableMethod := reflect.ValueOf(object).MethodByName(method)
	if callableMethod.IsValid() {
		// callableMethod.Call(inputs)[0].Interface()
		callableMethod.Call(inputs)
		return
	} else {
		panic(fmt.Sprintf("unable to locate the method %v", method))
	}

}
