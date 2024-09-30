package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/twreporter/go-api/internal/graphql"
)

// todo: remove this method before release
func TestGraphQL(c *gin.Context) (int, gin.H, error) {
	validOrder := "twreporter-167524234768464095200"
	reqOrder := c.Query("order")
	if reqOrder != validOrder {
		return http.StatusBadRequest, nil, errors.New("invalid params")
	}
	req := graphql.NewRequest(`
		query Query($where: PrimeDonationWhereInput!) {
  		primeDonations(where: $where) {
    		user {
      		email
    		}
  		}
		}
	`)
	type Where struct {
		Order struct {
			Equals string `json:"equals"`
		} `json:"order_number"`
	}
	whereVar := Where{}
	whereVar.Order.Equals = validOrder

	req.Var("where", whereVar)
	res, err := graphql.Query(req)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	return http.StatusOK, gin.H{
		"res": res,
	}, nil
}
