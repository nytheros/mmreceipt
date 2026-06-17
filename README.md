# Mattermost Read Receipts Plugin

WhatsApp-style read receipts for Mattermost 10.x Direct Messages and Group Messages.

## Features

- Sent (`✓`), delivered (`✓✓` gray), and read (`✓✓` blue) indicators.
- Read detection requires the message to be at least 70% visible for 2 continuous seconds using `IntersectionObserver`.
- Delivery detection is based on client acknowledgements over authenticated plugin REST endpoints, not online/offline status.
- Receipt records are persisted in the Mattermost plugin KV store.
- WebSocket updates refresh sender UI without page reloads.
- Admin settings for enabling the plugin and selecting DM, GM, or both.

## Build

```bash
make dist
```

Upload `dist/readreceipt.tar.gz` in **System Console -> Plugin Management -> Upload Plugin**.

## Configuration

The plugin is disabled by default. Enable **Enable Read Receipts** and choose **Enable For** from `DM`, `GM`, or `DM_AND_GM`.

## Security

The REST API validates authenticated users, restricts support to DM/GM channels, requires recipients to be channel members before writing receipts, and only allows senders to retrieve receipt status for their own messages.
