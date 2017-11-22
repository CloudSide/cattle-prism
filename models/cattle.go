package models

import (

)

type CattlePagination struct {
    First               string                                  `json:"first,omitempty"`
    Previous            string                                  `json:"previous,omitempty"`
    Next                string                                  `json:"next,omitempty"`
    Last                string                                  `json:"last,omitempty"`
    Limit               int                                     `json:"limit,omitempty"`
    Total               int                                     `json:"total,omitempty"`
    Partial             bool                                    `json:"partial"`
}

type CattleFilterAttribute struct {
    Modifier            string                                  `json:"modifier,omitempty"`
    Value               interface{}                             `json:"value,omitempty"`
}

type CattleSort struct {
    Name                string                                  `json:"name,omitempty"`
    Order               string                                  `json:"order,omitempty"`
    Reverse             string                                  `json:"reverse,omitempty"`
}

type CattleResource struct {
    Id                  string                                  `json:"id,omitempty"`
	Type                string                                  `json:"type,omitempty"`
	ResourceType        string                                  `json:"resourceType,omitempty"`
	Links			    map[string]string                       `json:"links,omitempty"`
    SortLinks           map[string]string                       `json:"sortLinks,omitempty"`
    Actions			    map[string]string                       `json:"actions,omitempty"`
    CreateTypes         map[string]string                       `json:"createTypes,omitempty"`
    Data                []interface{}                           `json:"data,omitempty"`
    Pagination          CattlePagination                        `json:"pagination,omitempty"`
    Sort                CattleSort                              `json:"sort,omitempty"`
    Filters             map[string][]CattleFilterAttribute      `json:"filters,omitempty"`
}

type CattleError struct {
    Id                  string                                  `json:"id,omitempty"`
    Type                string                                  `json:"type,omitempty"`
    Links			    map[string]string                       `json:"links,omitempty"`
    Actions			    map[string]string                       `json:"actions,omitempty"`
    Status              int                                     `json:"status,omitempty"`
    Code                string                                  `json:"code,omitempty"`
    Message             string                                  `json:"message,omitempty"`
    Detail              string                                  `json:"detail,omitempty"`
    BaseType            string                                  `json:"baseType,omitempty"`
}
