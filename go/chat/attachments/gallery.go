package attachments

import (
	"github.com/keybase/client/go/chat/globals"
	"github.com/keybase/client/go/chat/utils"
	"github.com/keybase/client/go/protocol/chat1"
	"github.com/keybase/client/go/protocol/gregor1"
	"golang.org/x/net/context"
)

type NextMessageOptions struct {
	BackInTime bool
	ImagesOnly bool
}

type Gallery struct {
	globals.Contextified
	utils.DebugLabeler

	PrevStride, NextStride int
}

func NewGallery(g *globals.Context) *Gallery {
	return &Gallery{
		Contextified: globals.NewContextified(g),
		DebugLabeler: utils.NewDebugLabeler(g.GetLog(), "Attachments.Gallery", false),
		NextStride:   10,
		PrevStride:   50,
	}
}

func (g *Gallery) NextMessage(ctx context.Context, uid gregor1.UID,
	convID chat1.ConversationID, msgID chat1.MessageID, opts NextMessageOptions) (res *chat1.MessageUnboxed, err error) {
	eligible := func(msg chat1.MessageUnboxed) bool {
		if !msg.IsValid() {
			return false
		}
		body := msg.Valid().MessageBody
		typ, err := body.MessageType()
		if err != nil || typ != chat1.MessageType_ATTACHMENT {
			return false
		}
		if !opts.ImagesOnly {
			return true
		}
		md := body.Attachment().Object.Metadata
		atyp, err := md.AssetType()
		if err != nil {
			return false
		}
		return atyp == chat1.AssetMetadataType_IMAGE || atyp == chat1.AssetMetadataType_VIDEO
	}
	if opts.BackInTime {
		pagination := utils.XlateMessageIDControlToPagination(&chat1.MessageIDControl{
			Pivot: &msgID,
			Num:   g.NextStride,
		})
		for {
			g.Debug(ctx, "NextMessage: backintime starting scan: next: %x", pagination.Next)
			tv, err := g.G().ConvSource.Pull(ctx, convID, uid, chat1.GetThreadReason_GENERAL,
				&chat1.GetThreadQuery{
					MessageTypes: []chat1.MessageType{chat1.MessageType_ATTACHMENT},
				}, pagination)
			if err != nil {
				return res, err
			}
			for _, m := range tv.Messages {
				if !eligible(m) {
					continue
				}
				return &m, nil
			}
			g.Debug(ctx, "NextMessage: missed all messages: len: %d", len(tv.Messages))
			if tv.Pagination.Last {
				break
			}
			pagination = tv.Pagination
			pagination.Num = g.NextStride
			pagination.Previous = nil
		}
	} else {
		pivot := msgID
		for {
			g.Debug(ctx, "NextMessage: forwardintime starting scan: pivot: %v", pivot)
			// Move forward in the thread looking for attachments, 50 messages at a time
			pagination := utils.XlateMessageIDControlToPagination(&chat1.MessageIDControl{
				Pivot:  &pivot,
				Num:    g.PrevStride,
				Recent: true,
			})
			tv, err := g.G().ConvSource.Pull(ctx, convID, uid, chat1.GetThreadReason_GENERAL,
				&chat1.GetThreadQuery{
					MessageTypes: []chat1.MessageType{chat1.MessageType_ATTACHMENT},
				}, pagination)
			if err != nil {
				return res, err
			}
			if len(tv.Messages) > 0 && tv.Messages[0].GetMessageID() != msgID {
				g.Debug(ctx, "NextMessage: hit: len: %d top: %d bottom: %d",
					len(tv.Messages), tv.Messages[0].GetMessageID(),
					tv.Messages[len(tv.Messages)-1].GetMessageID())
				for i := len(tv.Messages) - 1; i >= 0; i-- {
					if !eligible(tv.Messages[i]) {
						continue
					}
					return &tv.Messages[i], nil
				}
			}
			if tv.Pagination.Last {
				break
			}
			pivot += chat1.MessageID(g.PrevStride)
		}
	}
	return res, nil
}
