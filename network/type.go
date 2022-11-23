package network

type LibAddReq struct {
	Url string `json:"url"`
}
type BaseResp struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Result  []string `jaon:"result"`
}
