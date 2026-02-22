# Dialogix Customizations (WhatsMeow Fork)

This changelog tracks changes applied on top of the upstream `tulir/whatsmeow` for Dialogix.

## 2026-02-21

### Added
- `mentions.go`
  - `SendGhostMentionAll` helper to mention all group participants without showing @mentions in text.
- `extensions/interactive/interactive.go`
  - Builders: `NewListMessage`, `NewButtonsMessage`, `NewNativeFlowMessage`
  - Send helpers: `SendListMessage`, `SendButtonsMessage`, `SendNativeFlowMessage`
  - Validation: `Validate*`, `Validate*WithLimits`
  - Configurable limits via `ValidationLimits` and `DefaultLimits`
- `extensions/interactive/README.md` with examples and limits

### Notes
- NativeFlow JSON validation is enforced for `MessageParamsJSON` and `ButtonParamsJSON`.
- Limits are configurable via `Validate*WithLimits`.
