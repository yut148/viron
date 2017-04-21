package controller

import (
	"encoding/json"

	// for codegen.ParseDSL
	_ "github.com/cam-inc/dmc/example-go/design"

	"strings"

	"github.com/cam-inc/dmc/example-go/bridge"
	"github.com/cam-inc/dmc/example-go/common"
	"github.com/cam-inc/dmc/example-go/gen/app"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/goagen/codegen"
	"github.com/goadesign/goa/goagen/gen_swagger"
)

var swaggerAll *genswagger.Swagger

func init() {
	codegen.ParseDSL()
	sw, err := genswagger.New(design.Design)
	if err != nil {
		panic(err)
	}
	swaggerAll = sw
}

func filter(s genswagger.Swagger, roles map[string][]string) genswagger.Swagger {
	if s.Paths != nil {
		for uri, p := range s.Paths {
			path, ok := p.(*genswagger.Path)
			if ok != true {
				continue
			}

			resource := strings.Split(uri, "/")[1]
			raw, _ := path.MarshalJSON()
			mt := map[string]interface{}{}
			json.Unmarshal(raw, &mt)

			for method := range mt {
				if roles[method] == nil || (common.InStringArray("*", roles[method]) < 0 && common.InStringArray(resource, roles[method]) < 0) {
					delete(mt, method)
				}
			}

			newRaw, _ := json.Marshal(&mt)
			var newPath genswagger.Path
			json.Unmarshal(newRaw, &newPath)

			if len(mt) <= 0 {
				delete(s.Paths, uri)
			} else {
				s.Paths[uri] = newPath
			}
		}
	}
	return s
}

// SwaggerController implements the swagger resource.
type SwaggerController struct {
	*goa.Controller
}

// NewSwaggerController creates a swagger controller.
func NewSwaggerController(service *goa.Service) *SwaggerController {
	return &SwaggerController{Controller: service.NewController("SwaggerController")}
}

// Show runs the show action.
func (c *SwaggerController) Show(ctx *app.ShowSwaggerContext) error {
	// DmcController_Show: start_implement

	// Put your logic here
	var sw genswagger.Swagger

	cl := ctx.Context.Value(bridge.JwtClaims)
	if cl == nil {
		// swagger.json自体に認証をかけていないときは全部返す
		sw = *swaggerAll
	} else {
		// JWTclaimsからRoleを取り出す
		var roles map[string][]string
		claims := cl.(jwtgo.MapClaims)
		json.Unmarshal([]byte(claims["roles"].(string)), &roles)
		sw = filter(*swaggerAll, roles)
	}
	res, err := json.Marshal(&sw)
	if err != nil {
		panic(err)
	}

	// DmcController_Show: end_implement
	return ctx.OK(res)
}