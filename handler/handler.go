package handler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/micro/go-micro/errors"
	proto2 "github.com/micro/message-srv/proto/message"
	proto "github.com/micro/micro/api/proto"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
)

type Message struct{}

var (
	Client proto2.MessageClient
)

func extractValue(pair *proto.Pair) string {
	if pair == nil {
		return ""
	}
	if len(pair.Values) == 0 {
		return ""
	}
	return pair.Values[0]
}

func (m *Message) Create(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	var cr *proto2.CreateRequest
	var err error

	if d := extractValue(req.Post["event"]); len(d) > 0 {
		var ev *proto2.Event
		err = json.Unmarshal([]byte(d), &ev)
		if err == nil {
			cr = &proto2.CreateRequest{Event: ev}
		}
	} else {
		err = json.Unmarshal([]byte(req.Body), &cr)
	}

	if err != nil {
		return errors.BadRequest("go.micro.api.message", "invalid event")
	}

	if len(cr.Event.Channel) == 0 {
		return errors.BadRequest("go.micro.api.message", "invalid channel")
	}

	if len(cr.Event.Text) == 0 {
		return errors.BadRequest("go.micro.api.message", "invalid text")
	}

	cr.Event.Namespace = "default"
	cr.Event.Id = uuid.NewUUID().String()
	cr.Event.Created = time.Now().UnixNano()
	cr.Event.Updated = 0

	if _, err := Client.Create(ctx, cr); err != nil {
		return errors.InternalServerError("go.micro.api.message", err.Error())
	}

	rsp.StatusCode = 200
	rsp.Body = `{}`
	return nil
}

func (m *Message) Update(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	var cr *proto2.CreateRequest
	var err error

	if d := extractValue(req.Post["event"]); len(d) > 0 {
		var ev *proto2.Event
		err = json.Unmarshal([]byte(d), &ev)
		if err == nil {
			cr = &proto2.CreateRequest{Event: ev}
		}
	} else {
		err = json.Unmarshal([]byte(req.Body), &cr)
	}

	if err != nil {
		return errors.BadRequest("go.micro.api.message", "invalid event")
	}

	if len(cr.Event.Id) == 0 {
		return errors.BadRequest("go.micro.api.message", "invalid id")
	}

	if len(cr.Event.Channel) == 0 {
		return errors.BadRequest("go.micro.api.message", "invalid channel")
	}

	if len(cr.Event.Text) == 0 {
		return errors.BadRequest("go.micro.api.message", "invalid text")
	}

	if cr.Event.Created == 0 {
		return errors.BadRequest("go.micro.api.message", "invalid created")
	}

	cr.Event.Namespace = "default"
	cr.Event.Updated = time.Now().UnixNano()

	if _, err := Client.Create(ctx, cr); err != nil {
		return errors.InternalServerError("go.micro.api.message", err.Error())
	}

	rsp.StatusCode = 200
	rsp.Body = `{}`
	return nil
}

func (m *Message) Delete(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	var dr *proto2.DeleteRequest
	var err error

	channel := extractValue(req.Post["channel"])
	id := extractValue(req.Post["id"])

	if len(channel) > 0 && len(id) > 0 {
		dr = &proto2.DeleteRequest{
			Channel: channel,
			Id:      id,
		}

	} else {
		err = json.Unmarshal([]byte(req.Body), &dr)
	}

	if err != nil {
		return errors.BadRequest("go.micro.api.message", "invalid request")
	}

	if len(dr.Channel) == 0 {
		return errors.BadRequest("go.micro.api.message", "invalid channel")
	}

	if len(dr.Id) == 0 {
		return errors.BadRequest("go.micro.api.message", "invalid id")
	}

	dr.Namespace = "default"

	if _, err := Client.Delete(ctx, dr); err != nil {
		return errors.InternalServerError("go.micro.api.message", err.Error())
	}

	rsp.StatusCode = 200
	rsp.Body = `{}`
	return nil
}

func (m *Message) Search(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	var sr *proto2.SearchRequest
	var err error

	query := extractValue(req.Post["query"])
	channel := extractValue(req.Post["channel"])
	limit, _ := strconv.ParseInt(extractValue(req.Post["limit"]), 10, 64)
	offset, _ := strconv.ParseInt(extractValue(req.Post["offset"]), 10, 64)

	if len(req.Body) > 0 {
		err = json.Unmarshal([]byte(req.Body), &sr)
	} else {
		sr = &proto2.SearchRequest{
			Query:   query,
			Channel: channel,
			Limit:   limit,
			Offset:  offset,
			Reverse: true,
		}
	}
	if err != nil {
		return errors.BadRequest("go.micro.api.message", "invalid request")
	}

	if len(sr.Channel) == 0 {
		return errors.BadRequest("go.micro.api.message", "invalid channel")
	}

	sr.Namespace = "default"

	srsp, err := Client.Search(ctx, sr)
	if err != nil {
		return errors.InternalServerError("go.micro.api.message", err.Error())
	}

	b, _ := json.Marshal(srsp)

	rsp.StatusCode = 200
	rsp.Body = string(b)
	return nil
}

func (m *Message) Read(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	var rr *proto2.ReadRequest
	var err error

	channel := extractValue(req.Post["channel"])
	id := extractValue(req.Post["id"])

	if len(channel) > 0 && len(id) > 0 {
		rr = &proto2.ReadRequest{
			Channel: channel,
			Id:      id,
		}

	} else {
		err = json.Unmarshal([]byte(req.Body), &rr)
	}

	fmt.Println(req.Body)

	if err != nil {
		return errors.BadRequest("go.micro.api.message", "invalid request")
	}

	if len(rr.Channel) == 0 {
		return errors.BadRequest("go.micro.api.message", "invalid channel")
	}

	if len(rr.Id) == 0 {
		return errors.BadRequest("go.micro.api.message", "invalid id")
	}

	rr.Namespace = "default"

	rrsp, err := Client.Read(ctx, rr)
	if err != nil {
		return errors.InternalServerError("go.micro.api.message", err.Error())
	}

	b, _ := json.Marshal(rrsp)

	rsp.StatusCode = 200
	rsp.Body = string(b)
	return nil
}
