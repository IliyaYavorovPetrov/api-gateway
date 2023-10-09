package models

import (
	"encoding/json"

	"fmt"
)

type ReqRoutingInfo struct {
	SourceURL      string `redis:"sourceURL" json:"sourceURL"`
	DestinationURL string `redis:"destinationURL" json:"destinationURL"`
	MethodHTTP     string `redis:"methodHTTP" json:"methodHTTP"`
	IsAuthNeeded   bool   `redis:"isAuthNeeded" json:"isAuthNeeded"`
}

func (rri ReqRoutingInfo) Equals(other interface{}) bool {
	if session, ok := other.(ReqRoutingInfo); ok {
		return rri == session
	}

	return false
}

func (rri ReqRoutingInfo) ToString() string {
	data, err := json.Marshal(rri)
	if err != nil {
		fmt.Println("error encoding to json:", err)
		return ""
	}

	return string(data)
}
