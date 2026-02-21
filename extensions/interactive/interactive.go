package interactive

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

type ValidationLimits struct {
	MaxButtons              int
	MaxListSections         int
	MaxRowsPerSection       int
	MaxTotalListRows        int
	MaxTitleLength          int
	MaxSubtitleLength       int
	MaxBodyLength           int
	MaxFooterLength         int
	MaxButtonTextLength     int
	MaxRowTitleLength       int
	MaxRowDescriptionLength int
	MaxJSONLength           int
}

var DefaultLimits = ValidationLimits{
	MaxButtons:              3,
	MaxListSections:         10,
	MaxRowsPerSection:       10,
	MaxTotalListRows:        100,
	MaxTitleLength:          60,
	MaxSubtitleLength:       60,
	MaxBodyLength:           1024,
	MaxFooterLength:         60,
	MaxButtonTextLength:     20,
	MaxRowTitleLength:       24,
	MaxRowDescriptionLength: 72,
	MaxJSONLength:           2048,
}

// Button represents a quick-reply button.
type Button struct {
	ID   string
	Text string
}

// ListRow represents a row in a list section.
type ListRow struct {
	ID          string
	Title       string
	Description string
}

// ListSection groups rows in a list message.
type ListSection struct {
	Title string
	Rows  []ListRow
}

// ListMessageParams defines the list message payload.
type ListMessageParams struct {
	Title       string
	Description string
	ButtonText  string
	FooterText  string
	Sections    []ListSection
}

// ButtonsMessageParams defines the buttons message payload.
type ButtonsMessageParams struct {
	ContentText string
	FooterText  string
	Buttons     []Button
}

// NativeFlowButton represents a native flow button payload.
type NativeFlowButton struct {
	Name             string
	ButtonParamsJSON string
}

// NativeFlowMessageParams defines the native flow interactive message payload.
type NativeFlowMessageParams struct {
	HeaderTitle       string
	HeaderSubtitle    string
	BodyText          string
	FooterText        string
	MessageParamsJSON string
	MessageVersion    int32
	Buttons           []NativeFlowButton
}

// NewListMessage builds a list message payload.
func NewListMessage(params ListMessageParams) *waE2E.Message {
	sections := make([]*waE2E.ListMessage_Section, 0, len(params.Sections))
	for _, section := range params.Sections {
		rows := make([]*waE2E.ListMessage_Row, 0, len(section.Rows))
		for _, row := range section.Rows {
			rows = append(rows, &waE2E.ListMessage_Row{
				RowID:       strPtr(row.ID),
				Title:       strPtr(row.Title),
				Description: strPtr(row.Description),
			})
		}
		sections = append(sections, &waE2E.ListMessage_Section{
			Title: strPtr(section.Title),
			Rows:  rows,
		})
	}

	listMessage := &waE2E.ListMessage{
		Title:       strPtr(params.Title),
		Description: strPtr(params.Description),
		ButtonText:  strPtr(params.ButtonText),
		FooterText:  strPtr(params.FooterText),
		ListType:    waE2E.ListMessage_SINGLE_SELECT.Enum(),
		Sections:    sections,
	}

	return &waE2E.Message{ListMessage: listMessage}
}

// NewButtonsMessage builds a quick-reply buttons message payload.
func NewButtonsMessage(params ButtonsMessageParams) *waE2E.Message {
	buttons := make([]*waE2E.ButtonsMessage_Button, 0, len(params.Buttons))
	for _, btn := range params.Buttons {
		buttons = append(buttons, &waE2E.ButtonsMessage_Button{
			ButtonID:   strPtr(btn.ID),
			ButtonText: &waE2E.ButtonsMessage_Button_ButtonText{DisplayText: strPtr(btn.Text)},
			Type:       waE2E.ButtonsMessage_Button_RESPONSE.Enum(),
		})
	}

	buttonsMessage := &waE2E.ButtonsMessage{
		ContentText: strPtr(params.ContentText),
		FooterText:  strPtr(params.FooterText),
		Buttons:     buttons,
		HeaderType:  waE2E.ButtonsMessage_EMPTY.Enum(),
	}

	return &waE2E.Message{ButtonsMessage: buttonsMessage}
}

// NewNativeFlowMessage builds a native flow interactive message payload.
func NewNativeFlowMessage(params NativeFlowMessageParams) *waE2E.Message {
	buttons := make([]*waE2E.InteractiveMessage_NativeFlowMessage_NativeFlowButton, 0, len(params.Buttons))
	for _, btn := range params.Buttons {
		buttons = append(buttons, &waE2E.InteractiveMessage_NativeFlowMessage_NativeFlowButton{
			Name:             strPtr(btn.Name),
			ButtonParamsJSON: strPtr(btn.ButtonParamsJSON),
		})
	}

	version := params.MessageVersion
	if version == 0 {
		version = 1
	}

	nativeFlow := &waE2E.InteractiveMessage_NativeFlowMessage{
		Buttons:           buttons,
		MessageParamsJSON: strPtr(params.MessageParamsJSON),
		MessageVersion:    proto.Int32(version),
	}

	interactive := &waE2E.InteractiveMessage{
		InteractiveMessage: &waE2E.InteractiveMessage_NativeFlowMessage_{
			NativeFlowMessage: nativeFlow,
		},
		Header: &waE2E.InteractiveMessage_Header{
			Title:              strPtr(params.HeaderTitle),
			Subtitle:           strPtr(params.HeaderSubtitle),
			HasMediaAttachment: proto.Bool(false),
		},
		Body: &waE2E.InteractiveMessage_Body{
			Text: strPtr(params.BodyText),
		},
		Footer: &waE2E.InteractiveMessage_Footer{
			Text: strPtr(params.FooterText),
		},
	}

	return &waE2E.Message{
		MessageContextInfo: &waE2E.MessageContextInfo{
			DeviceListMetadataVersion: proto.Int32(3),
		},
		InteractiveMessage: interactive,
	}
}

// SendNativeFlowMessage sends a native flow interactive message using the provided client.
func SendNativeFlowMessage(ctx context.Context, client *whatsmeow.Client, jid types.JID, params NativeFlowMessageParams) (whatsmeow.SendResponse, error) {
	if err := ValidateNativeFlowMessage(params); err != nil {
		return whatsmeow.SendResponse{}, err
	}
	return client.SendMessage(ctx, jid, NewNativeFlowMessage(params))
}

// SendListMessage sends a list message using the provided client.
func SendListMessage(ctx context.Context, client *whatsmeow.Client, jid types.JID, params ListMessageParams) (whatsmeow.SendResponse, error) {
	if err := ValidateListMessage(params); err != nil {
		return whatsmeow.SendResponse{}, err
	}
	return client.SendMessage(ctx, jid, NewListMessage(params))
}

// SendButtonsMessage sends a buttons message using the provided client.
func SendButtonsMessage(ctx context.Context, client *whatsmeow.Client, jid types.JID, params ButtonsMessageParams) (whatsmeow.SendResponse, error) {
	if err := ValidateButtonsMessage(params); err != nil {
		return whatsmeow.SendResponse{}, err
	}
	return client.SendMessage(ctx, jid, NewButtonsMessage(params))
}

// ValidateListMessage validates list message params with DefaultLimits.
func ValidateListMessage(params ListMessageParams) error {
	return ValidateListMessageWithLimits(params, DefaultLimits)
}

// ValidateListMessageWithLimits validates list message params with custom limits.
func ValidateListMessageWithLimits(params ListMessageParams, limits ValidationLimits) error {
	l := normalizeLimits(limits)
	if strings.TrimSpace(params.Title) == "" {
		return errors.New("list title is required")
	}
	if len(params.Title) > l.MaxTitleLength {
		return fmt.Errorf("list title exceeds %d characters", l.MaxTitleLength)
	}
	if strings.TrimSpace(params.ButtonText) == "" {
		return errors.New("list button text is required")
	}
	if len(params.ButtonText) > l.MaxButtonTextLength {
		return fmt.Errorf("list button text exceeds %d characters", l.MaxButtonTextLength)
	}
	if len(params.Description) > l.MaxBodyLength {
		return fmt.Errorf("list description exceeds %d characters", l.MaxBodyLength)
	}
	if len(params.FooterText) > l.MaxFooterLength {
		return fmt.Errorf("list footer exceeds %d characters", l.MaxFooterLength)
	}
	if len(params.Sections) == 0 {
		return errors.New("list requires at least one section")
	}
	if len(params.Sections) > l.MaxListSections {
		return fmt.Errorf("list supports up to %d sections", l.MaxListSections)
	}
	rowCount := 0
	for i, section := range params.Sections {
		if strings.TrimSpace(section.Title) == "" {
			return fmt.Errorf("section %d title is required", i+1)
		}
		if len(section.Title) > l.MaxTitleLength {
			return fmt.Errorf("section %d title exceeds %d characters", i+1, l.MaxTitleLength)
		}
		if len(section.Rows) == 0 {
			return fmt.Errorf("section %d requires at least one row", i+1)
		}
		if len(section.Rows) > l.MaxRowsPerSection {
			return fmt.Errorf("section %d supports up to %d rows", i+1, l.MaxRowsPerSection)
		}
		for j, row := range section.Rows {
			rowCount++
			if strings.TrimSpace(row.ID) == "" {
				return fmt.Errorf("section %d row %d id is required", i+1, j+1)
			}
			if strings.TrimSpace(row.Title) == "" {
				return fmt.Errorf("section %d row %d title is required", i+1, j+1)
			}
			if len(row.Title) > l.MaxRowTitleLength {
				return fmt.Errorf("section %d row %d title exceeds %d characters", i+1, j+1, l.MaxRowTitleLength)
			}
			if len(row.Description) > l.MaxRowDescriptionLength {
				return fmt.Errorf("section %d row %d description exceeds %d characters", i+1, j+1, l.MaxRowDescriptionLength)
			}
		}
	}
	if rowCount > l.MaxTotalListRows {
		return fmt.Errorf("list supports up to %d total rows", l.MaxTotalListRows)
	}
	return nil
}

// ValidateButtonsMessage validates buttons message params with DefaultLimits.
func ValidateButtonsMessage(params ButtonsMessageParams) error {
	return ValidateButtonsMessageWithLimits(params, DefaultLimits)
}

// ValidateButtonsMessageWithLimits validates buttons message params with custom limits.
func ValidateButtonsMessageWithLimits(params ButtonsMessageParams, limits ValidationLimits) error {
	l := normalizeLimits(limits)
	if strings.TrimSpace(params.ContentText) == "" {
		return errors.New("buttons content text is required")
	}
	if len(params.ContentText) > l.MaxBodyLength {
		return fmt.Errorf("buttons content text exceeds %d characters", l.MaxBodyLength)
	}
	if len(params.FooterText) > l.MaxFooterLength {
		return fmt.Errorf("buttons footer text exceeds %d characters", l.MaxFooterLength)
	}
	if len(params.Buttons) == 0 {
		return errors.New("buttons message requires at least one button")
	}
	if len(params.Buttons) > l.MaxButtons {
		return fmt.Errorf("buttons message supports up to %d buttons", l.MaxButtons)
	}
	for i, btn := range params.Buttons {
		if strings.TrimSpace(btn.ID) == "" {
			return fmt.Errorf("button %d id is required", i+1)
		}
		if strings.TrimSpace(btn.Text) == "" {
			return fmt.Errorf("button %d text is required", i+1)
		}
		if len(btn.Text) > l.MaxButtonTextLength {
			return fmt.Errorf("button %d text exceeds %d characters", i+1, l.MaxButtonTextLength)
		}
	}
	return nil
}

// ValidateNativeFlowMessage validates native flow params with DefaultLimits.
func ValidateNativeFlowMessage(params NativeFlowMessageParams) error {
	return ValidateNativeFlowMessageWithLimits(params, DefaultLimits)
}

// ValidateNativeFlowMessageWithLimits validates native flow params with custom limits.
func ValidateNativeFlowMessageWithLimits(params NativeFlowMessageParams, limits ValidationLimits) error {
	l := normalizeLimits(limits)
	if strings.TrimSpace(params.BodyText) == "" {
		return errors.New("native flow body text is required")
	}
	if len(params.HeaderTitle) > l.MaxTitleLength {
		return fmt.Errorf("native flow header title exceeds %d characters", l.MaxTitleLength)
	}
	if len(params.HeaderSubtitle) > l.MaxSubtitleLength {
		return fmt.Errorf("native flow header subtitle exceeds %d characters", l.MaxSubtitleLength)
	}
	if len(params.BodyText) > l.MaxBodyLength {
		return fmt.Errorf("native flow body text exceeds %d characters", l.MaxBodyLength)
	}
	if len(params.FooterText) > l.MaxFooterLength {
		return fmt.Errorf("native flow footer text exceeds %d characters", l.MaxFooterLength)
	}
	if len(params.MessageParamsJSON) > l.MaxJSONLength {
		return fmt.Errorf("native flow message params json exceeds %d characters", l.MaxJSONLength)
	}
	if params.MessageParamsJSON != "" && !json.Valid([]byte(params.MessageParamsJSON)) {
		return errors.New("native flow message params json is invalid")
	}
	if len(params.Buttons) == 0 {
		return errors.New("native flow requires at least one button")
	}
	if len(params.Buttons) > l.MaxButtons {
		return fmt.Errorf("native flow supports up to %d buttons", l.MaxButtons)
	}
	for i, btn := range params.Buttons {
		if strings.TrimSpace(btn.Name) == "" {
			return fmt.Errorf("native flow button %d name is required", i+1)
		}
		if strings.TrimSpace(btn.ButtonParamsJSON) == "" {
			return fmt.Errorf("native flow button %d params json is required", i+1)
		}
		if len(btn.ButtonParamsJSON) > l.MaxJSONLength {
			return fmt.Errorf("native flow button %d params json exceeds %d characters", i+1, l.MaxJSONLength)
		}
		if !json.Valid([]byte(btn.ButtonParamsJSON)) {
			return fmt.Errorf("native flow button %d params json is invalid", i+1)
		}
	}
	return nil
}

func normalizeLimits(limits ValidationLimits) ValidationLimits {
	if limits.MaxButtons == 0 {
		limits.MaxButtons = DefaultLimits.MaxButtons
	}
	if limits.MaxListSections == 0 {
		limits.MaxListSections = DefaultLimits.MaxListSections
	}
	if limits.MaxRowsPerSection == 0 {
		limits.MaxRowsPerSection = DefaultLimits.MaxRowsPerSection
	}
	if limits.MaxTotalListRows == 0 {
		limits.MaxTotalListRows = DefaultLimits.MaxTotalListRows
	}
	if limits.MaxTitleLength == 0 {
		limits.MaxTitleLength = DefaultLimits.MaxTitleLength
	}
	if limits.MaxSubtitleLength == 0 {
		limits.MaxSubtitleLength = DefaultLimits.MaxSubtitleLength
	}
	if limits.MaxBodyLength == 0 {
		limits.MaxBodyLength = DefaultLimits.MaxBodyLength
	}
	if limits.MaxFooterLength == 0 {
		limits.MaxFooterLength = DefaultLimits.MaxFooterLength
	}
	if limits.MaxButtonTextLength == 0 {
		limits.MaxButtonTextLength = DefaultLimits.MaxButtonTextLength
	}
	if limits.MaxRowTitleLength == 0 {
		limits.MaxRowTitleLength = DefaultLimits.MaxRowTitleLength
	}
	if limits.MaxRowDescriptionLength == 0 {
		limits.MaxRowDescriptionLength = DefaultLimits.MaxRowDescriptionLength
	}
	if limits.MaxJSONLength == 0 {
		limits.MaxJSONLength = DefaultLimits.MaxJSONLength
	}
	return limits
}

func strPtr(value string) *string {
	if value == "" {
		return nil
	}
	return proto.String(value)
}
