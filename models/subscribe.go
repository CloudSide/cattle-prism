package models

import (

)

type SubscribeDataResource struct {
    Id                  string                              `json:id,omitempty`
	Type                string                              `json:type,omitempty`
    AccountId           string                              `json:accountId,omitempty`
    CompanyId           string                              `json:companyId,omitempty`
}

type SubscribeData struct {
    Options             map[string]interface{}              `json:options,omitempty`
    Resources           []interface{}                       `json:resources,omitempty`
    Resource            SubscribeDataResource               `json:resource,omitempty`
}

type SubscribeResource struct {
    Id                  string                              `json:id,omitempty`
    Name                string                              `json:name,omitempty`
    ResourceId          string                              `json:resourceId,omitempty`
    ResourceType        string                              `json:resourceType,omitempty`
    Time                int64                               `json:time,omitempty`
    Data                SubscribeData                       `json:data,omitempty"`
}
