# Interactive helpers (buttons, list, native flow)

This package provides small helpers for sending WhatsApp interactive messages with WhatsMeow.

## What was added

Files:
- `extensions/interactive/interactive.go`
- `extensions/interactive/README.md`

Exports:
- **Builders**
  - `NewListMessage(params ListMessageParams) *waE2E.Message`
  - `NewButtonsMessage(params ButtonsMessageParams) *waE2E.Message`
  - `NewNativeFlowMessage(params NativeFlowMessageParams) *waE2E.Message`
- **Send helpers**
  - `SendListMessage(ctx, client, jid, params)`
  - `SendButtonsMessage(ctx, client, jid, params)`
  - `SendNativeFlowMessage(ctx, client, jid, params)`
- **Validation**
  - `ValidateListMessage(params)`
  - `ValidateButtonsMessage(params)`
  - `ValidateNativeFlowMessage(params)`
  - `ValidateListMessageWithLimits(params, limits)`
  - `ValidateButtonsMessageWithLimits(params, limits)`
  - `ValidateNativeFlowMessageWithLimits(params, limits)`
- **Limits**
  - `ValidationLimits` struct
  - `DefaultLimits` values

## Types

```go
// Buttons

type Button struct {
  ID   string
  Text string
}

type ButtonsMessageParams struct {
  ContentText string
  FooterText  string
  Buttons     []Button
}

// List

type ListRow struct {
  ID          string
  Title       string
  Description string
}

type ListSection struct {
  Title string
  Rows  []ListRow
}

type ListMessageParams struct {
  Title       string
  Description string
  ButtonText  string
  FooterText  string
  Sections    []ListSection
}

// Native flow

type NativeFlowButton struct {
  Name             string
  ButtonParamsJSON string
}

type NativeFlowMessageParams struct {
  HeaderTitle       string
  HeaderSubtitle    string
  BodyText          string
  FooterText        string
  MessageParamsJSON string
  MessageVersion    int32
  Buttons           []NativeFlowButton
}
```

## Usage examples

### List message

```go
params := interactive.ListMessageParams{
  Title:      "Escolha uma opção",
  ButtonText: "Abrir lista",
  Sections: []interactive.ListSection{
    {
      Title: "Suporte",
      Rows: []interactive.ListRow{
        {ID: "1", Title: "Financeiro"},
        {ID: "2", Title: "Suporte Técnico"},
      },
    },
  },
}

_, err := interactive.SendListMessage(ctx, client, jid, params)
```

### Buttons message

```go
params := interactive.ButtonsMessageParams{
  ContentText: "Quer continuar?",
  FooterText:  "Dialogix",
  Buttons: []interactive.Button{
    {ID: "yes", Text: "Sim"},
    {ID: "no", Text: "Não"},
  },
}

_, err := interactive.SendButtonsMessage(ctx, client, jid, params)
```

### Native flow (interactive message)

```go
params := interactive.NativeFlowMessageParams{
  HeaderTitle: "Atendimento",
  BodyText:    "Escolha uma opção:",
  FooterText:  "Dialogix",
  MessageParamsJSON: `{"from":"api","templateId":"x"}`,
  Buttons: []interactive.NativeFlowButton{
    {Name: "quick_reply", ButtonParamsJSON: `{"display_text":"Botão 1","id":"btn1"}`},
    {Name: "quick_reply", ButtonParamsJSON: `{"display_text":"Botão 2","id":"btn2"}`},
  },
}

_, err := interactive.SendNativeFlowMessage(ctx, client, jid, params)
```

## Validation limits

Default limits (see `DefaultLimits`):
- Buttons: 3
- Sections: 10
- Rows per section: 10
- Total rows: 100
- Title length: 60
- Subtitle length: 60
- Body length: 1024
- Footer length: 60
- Button text length: 20
- Row title length: 24
- Row description length: 72
- JSON length: 2048

Custom limits example:

```go
limits := interactive.ValidationLimits{
  MaxButtons:    5,
  MaxBodyLength: 2048,
  MaxJSONLength: 4096,
}

err := interactive.ValidateNativeFlowMessageWithLimits(params, limits)
```

## Notes

- Native flow requires **valid JSON** for `MessageParamsJSON` and each `ButtonParamsJSON`.
- For media headers, extend `NewNativeFlowMessage` with `InteractiveMessage_Header` media fields as needed.
