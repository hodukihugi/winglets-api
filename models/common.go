package models

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type ArrStringFilterType string

func (t *ArrStringFilterType) Values() []string {
	if t == nil {
		return []string{}
	}
	return strings.Split(string(*t), ",")
}

func NewArrStringFilterType(args ...string) *ArrStringFilterType {
	res := ArrStringFilterType(strings.Join(args, ","))
	return &res
}

type Pagination struct {
	CurrentPage int `json:"current_page" form:"current_page"`
	PerPage     int `json:"per_page" form:"per_page"`
}

var DefaultPagination = Pagination{
	CurrentPage: 1,
	PerPage:     10,
}

// ParsePagination parses pagination from query string
func ParsePagination(ctx *gin.Context) (*Pagination, error) {
	var pagination Pagination
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return nil, errors.New(fmt.Sprintf("fail to bind pagination query [%s]", err.Error()))
	}
	if pagination.CurrentPage == 0 && pagination.PerPage == 0 {
		return &DefaultPagination, nil
	}
	if pagination.PerPage < 1 || pagination.CurrentPage < 1 {
		ctx.JSON(http.StatusBadRequest, nil)
		return nil, errors.New("invalid pagination")
	}
	return &pagination, nil
}

type PaginationResp struct {
	Pagination
	Count int64 `json:"count"`
}

type PaginationReq Pagination

type TaxRate string

const (
	TaxRateBASExcluded     TaxRate = "bas_excluded"
	TaxRateGSTOnIncome     TaxRate = "gst_on_income"
	TaxRateGSTOnExpenses   TaxRate = "gst_on_expenses"
	TaxRateGSTOnImports    TaxRate = "gst_on_imports"
	TaxRateGSTFreeExpenses TaxRate = "gst_free_expenses"
	TaxRateGSTFreeIncome   TaxRate = "gst_free_income"
)

func (tr TaxRate) Validate() bool {
	return tr == TaxRateBASExcluded ||
		tr == TaxRateGSTOnIncome ||
		tr == TaxRateGSTOnExpenses ||
		tr == TaxRateGSTOnImports ||
		tr == TaxRateGSTFreeExpenses ||
		tr == TaxRateGSTFreeIncome
}

func (tr TaxRate) TaxFree() bool {
	return tr == TaxRateGSTFreeExpenses ||
		tr == TaxRateGSTFreeIncome ||
		tr == TaxRateBASExcluded
}

type Currency string

const (
	CurrencyAUD Currency = "AUD"
)

func (c Currency) Validate() bool {
	return c == CurrencyAUD
}
