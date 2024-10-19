// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.1-0.20240304154729-dd082985a9b6 DO NOT EDIT.
package generated

import (
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Error defines model for Error.
type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// Brand defines model for Brand.
type Brand = string

// CategoriesFilter defines model for CategoriesFilter.
type CategoriesFilter = []string

// Limit defines model for Limit.
type Limit = int64

// Page defines model for Page.
type Page = int64

// Year defines model for Year.
type Year = int64

// UnexpectedError defines model for UnexpectedError.
type UnexpectedError = Error

// CreateAdminConfigFormdataBody defines parameters for CreateAdminConfig.
type CreateAdminConfigFormdataBody struct {
	Key   string `form:"key" json:"key"`
	Value string `form:"value" json:"value"`
}

// UpdateAdminConfigFormdataBody defines parameters for UpdateAdminConfig.
type UpdateAdminConfigFormdataBody struct {
	Value string `form:"value" json:"value"`
}

// CreateMakerPieAdminMultipartBody defines parameters for CreateMakerPieAdmin.
type CreateMakerPieAdminMultipartBody struct {
	Displayname      string             `json:"displayname"`
	Fresh            bool               `json:"fresh"`
	Id               string             `json:"id"`
	ImageFile        openapi_types.File `json:"imageFile"`
	Labels           string             `json:"labels"`
	Makerid          string             `json:"makerid"`
	PackCount        int32              `json:"packCount"`
	PackPriceInPence int32              `json:"packPriceInPence"`
	Validated        bool               `json:"validated"`
	WebLink          string             `json:"webLink"`
	Year             int32              `json:"year"`
}

// UpdateMakerPieMultipartBody defines parameters for UpdateMakerPie.
type UpdateMakerPieMultipartBody struct {
	Categories       []string            `json:"categories"`
	Displayname      string              `json:"displayname"`
	Fresh            bool                `json:"fresh"`
	ImageFile        *openapi_types.File `json:"imageFile,omitempty"`
	Labels           string              `json:"labels"`
	PackCount        int32               `json:"packCount"`
	PackPriceInPence int32               `json:"packPriceInPence"`
	Validated        bool                `json:"validated"`
	WebLink          string              `json:"webLink"`
}

// YearPersonalRankingParams defines parameters for YearPersonalRanking.
type YearPersonalRankingParams struct {
	Page  *Page  `form:"page,omitempty" json:"page,omitempty"`
	Limit *Limit `form:"limit,omitempty" json:"limit,omitempty"`

	// Categories Categories to filter pies on
	Categories *CategoriesFilter `form:"categories,omitempty" json:"categories,omitempty"`
}

// YearAllPiesParams defines parameters for YearAllPies.
type YearAllPiesParams struct {
	Page  *Page  `form:"page,omitempty" json:"page,omitempty"`
	Limit *Limit `form:"limit,omitempty" json:"limit,omitempty"`

	// Categories Categories to filter pies on
	Categories *CategoriesFilter `form:"categories,omitempty" json:"categories,omitempty"`
}

// YearBrandPiesParams defines parameters for YearBrandPies.
type YearBrandPiesParams struct {
	Page  *Page  `form:"page,omitempty" json:"page,omitempty"`
	Limit *Limit `form:"limit,omitempty" json:"limit,omitempty"`

	// Categories Categories to filter pies on
	Categories *CategoriesFilter `form:"categories,omitempty" json:"categories,omitempty"`
}

// YearCategoryPiesParams defines parameters for YearCategoryPies.
type YearCategoryPiesParams struct {
	Page  *Page  `form:"page,omitempty" json:"page,omitempty"`
	Limit *Limit `form:"limit,omitempty" json:"limit,omitempty"`

	// Categories Categories to filter pies on
	Categories *CategoriesFilter `form:"categories,omitempty" json:"categories,omitempty"`
}

// CreateAdminConfigFormdataRequestBody defines body for CreateAdminConfig for application/x-www-form-urlencoded ContentType.
type CreateAdminConfigFormdataRequestBody CreateAdminConfigFormdataBody

// UpdateAdminConfigFormdataRequestBody defines body for UpdateAdminConfig for application/x-www-form-urlencoded ContentType.
type UpdateAdminConfigFormdataRequestBody UpdateAdminConfigFormdataBody

// CreateMakerPieAdminMultipartRequestBody defines body for CreateMakerPieAdmin for multipart/form-data ContentType.
type CreateMakerPieAdminMultipartRequestBody CreateMakerPieAdminMultipartBody

// UpdateMakerPieMultipartRequestBody defines body for UpdateMakerPie for multipart/form-data ContentType.
type UpdateMakerPieMultipartRequestBody UpdateMakerPieMultipartBody
