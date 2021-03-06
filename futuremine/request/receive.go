package request

import (
	"github.com/Futuremine-chain/futuremine/futuremine/types"
	request2 "github.com/Futuremine-chain/futuremine/service/request"
	"github.com/Futuremine-chain/futuremine/tools/rlp"
)

const (
	maxSyncCount = 1000
	minSyncCount = 1
)

func (r *RequestHandler) respLastHeight(req *ReqStream) (*Response, error) {
	var message string
	var body []byte
	code := Success
	height := r.chain.LastHeight()
	body, err := rlp.EncodeToBytes(height)
	if err != nil {
		code = Failed
		message = err.Error()
	}
	response := NewResponse(code, message, body)
	return response, nil
}

func (r *RequestHandler) respSendMsg(req *ReqStream) (*Response, error) {
	defer func() {
		req.stream.Reset()
		req.stream.Close()
	}()

	var message string
	var body []byte
	code := Success
	msg, err := types.DecodeMessage(req.request.Body)
	if err != nil {
		code = Failed
		message = err.Error()
	} else {
		r.receiveMessage(msg.ToMessage())
	}
	response := NewResponse(code, message, body)
	return response, nil
}

func (r *RequestHandler) respSendBlock(req *ReqStream) (*Response, error) {
	var message string
	var body []byte
	code := Success
	rlpBlock, err := types.DecodeRlpBlock(req.request.Body)
	if err != nil {
		code = Failed
		message = err.Error()
	}
	r.receiveBlock(rlpBlock.ToBlock())
	response := NewResponse(code, message, body)
	return response, nil
}

func (r *RequestHandler) respGetBlocks(req *ReqStream) (*Response, error) {
	var message string
	var body []byte
	var params []uint64
	var height, count uint64
	var index uint64
	code := Success
	rlpBlocks := make([]*types.RlpBlock, 0)
	lastHeight := r.chain.LastHeight()
	err := rlp.DecodeBytes(req.request.Body, &params)
	if err != nil {
		code = Failed
		message = err.Error()
	} else if len(params) == 2 {
		height = params[0]
		count = params[1]
		if count < minSyncCount {
			count = minSyncCount
		} else if count > maxSyncCount {
			count = maxSyncCount
		}
		if lastHeight >= height {
			for lastHeight >= height && index < count {
				block, err := r.chain.GetRlpBlockHeight(height)
				if err != nil {
					code = Failed
					message = err.Error()
					response := NewResponse(code, message, body)
					return response, nil
				} else {
					rlpBlocks = append(rlpBlocks, block.(*types.RlpBlock))
					height++
					index++
				}
			}
			body, _ = types.EncodeRlpBlocks(rlpBlocks)
		} else {
			code = Failed
			message = request2.Err_BlockNotFound.Error()
		}
	} else {
		code = Failed
		message = "wrong params"
	}

	response := NewResponse(code, message, body)
	return response, nil
}

func (r *RequestHandler) respIsEqual(req *ReqStream) (*Response, error) {
	var message string
	var body []byte
	code := Success
	header, err := types.DecodeHeader(req.request.Body)
	if err != nil {
		code = Failed
		return NewResponse(code, message, body), nil
	}
	localHeader, err := r.chain.GetHeaderHeight(header.Height)
	if err != nil {
		code = Failed
		return NewResponse(code, message, body), nil
	}
	isEqual := localHeader.GetHash().IsEqual(header.Hash)
	body, _ = rlp.EncodeToBytes(isEqual)
	return NewResponse(code, message, body), nil
}

func (r *RequestHandler) respLocalInfo(req *ReqStream) (*Response, error) {
	var message string
	var body []byte
	code := Success

	if r.getLocal != nil {
		local := r.getLocal()
		body, _ = rlp.EncodeToBytes(local)
	} else {
		code = Failed
		message = "no local info"
	}
	return NewResponse(code, message, body), nil
}
