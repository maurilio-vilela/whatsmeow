// Copyright (c) 2021 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package whatsmeow

import (
	"context"
	"errors"

	"google.golang.org/protobuf/proto"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

// SendGhostMentionAll sends a text message to a group mentioning all participants
// without including any @mentions in the visible message text.
func (cli *Client) SendGhostMentionAll(ctx context.Context, groupJID types.JID, text string, extra ...SendRequestExtra) (SendResponse, error) {
	if groupJID.Server != types.GroupServer {
		return SendResponse{}, errors.New("ghost mention all is only supported for groups")
	}

	cachedData, err := cli.getCachedGroupData(ctx, groupJID)
	if err != nil {
		return SendResponse{}, err
	}

	mentioned := make([]string, 0, len(cachedData.Members))
	for _, member := range cachedData.Members {
		mentioned = append(mentioned, member.ToNonAD().String())
	}

	message := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: &waE2E.ContextInfo{
				MentionedJID: mentioned,
			},
		},
	}

	return cli.SendMessage(ctx, groupJID, message, extra...)
}
