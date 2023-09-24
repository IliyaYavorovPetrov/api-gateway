package routing

type ReqRoutingInfo struct {
	SourceURL      string `redis:"sourceURL" json:"sourceURL"`
	DestinationURL string `redis:"destinationURL" json:"destinationURL"`
	MethodHTTP     string `redis:"methodHTTP" json:"methodHTTP"`
	IsAuthNeeded   bool   `redis:"isAuthNeeded" json:"isAuthNeeded"`
}

func NewReqRoutingInfo(sourceURL string, destinationURL string, methodHTTP string, isAuthNeeded bool) (ReqRoutingInfo, error) {
	return ReqRoutingInfo{sourceURL, destinationURL, methodHTTP, isAuthNeeded}, nil
}
