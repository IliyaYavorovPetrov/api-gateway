package routing

type ReqRoutingInfo struct {
	SourceURL      string `redis:"sourceURL"`
	DestinationURL string `redis:"destinationURL"`
	MethodHTTP     string `redis:"methodHTTP"`
	IsAuthNeeded   bool   `redis:"isAuthNeeded"`
}

func NewReqRoutingInfo(sourceURL string, destinationURL string, methodHTTP string, isAuthNeeded bool) (ReqRoutingInfo, error) {
	return ReqRoutingInfo{sourceURL, destinationURL, methodHTTP, isAuthNeeded}, nil
}
