package schemas

import "time"

type Response struct {
	Ok bool `json:"ok"`
	Error interface{} `json:"error"`
	D interface{} `json:"d"` // Data, if OK
	
	Time time.Time `json:"time"`
}

func NewOkResponse(d interface{}) Response {
	return Response{
		Ok: true,
		Error: nil,
		D: d,

		Time: time.Now(),
	}
}

func NewErrResponse(errMsg string) Response {
	return Response{
		Ok: false,
		Error: errMsg,
		D: nil,
		
		Time: time.Now(),
	}
}