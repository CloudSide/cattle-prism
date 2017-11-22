package models

import (

)

type TokenResource struct {
    CattleResource
    Data                []TokenDataItem                     `json:"data,omitempty"`
}

type TokenDataItem struct {
    Id                  string                              `json:"id,omitempty"`
    Type                string                              `json:"type,omitempty"`
    Links               map[string]string                   `json:"links,omitempty"`
    Actions             map[string]string                   `json:"actions,omitempty"`
    BaseType            string                              `json:"baseType,omitempty"`
    AccountId           string                              `json:"accountId,omitempty"`
    AuthProvider        string                              `json:"authProvider,omitempty"`
    Code                string                              `json:"code,omitempty"`
    Enabled             bool                                `json:"enabled"`
    Jwt                 string                              `json:"jwt,omitempty"`
    RedirectUrl         string                              `json:"redirectUrl,omitempty"`
    Security            bool                                `json:"security"`
    User                string                              `json:"user,omitempty"`
    UserType            string                              `json:"userType,omitempty"`
    UserIdentity        TokenDataItemUserIdentity           `json:"userIdentity,omitempty"`
}

type TokenDataItemUserIdentity struct {
    Id                  string                              `json:"id,omitempty"`
    ExternalId          string                              `json:"externalId,omitempty"`
    ProfilePicture      string                              `json:"profilePicture,omitempty"`
    Name                string                              `json:"name,omitempty"`
    ExternalIdType      string                              `json:"externalIdType,omitempty"`
    ProfileUrl          string                              `json:"profileUrl,omitempty"`
    Login               string                              `json:"login,omitempty"`
    User                bool                                `json:"user"`
    CompanyId           string                              `json:"companyId,omitempty"`
    Cached              bool                                `json:"cached"`
}
