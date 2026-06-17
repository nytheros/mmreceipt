# Mattermost Read Receipts Plugin

WhatsApp-style delivery and read indicators for Mattermost Direct Message (DM) and Group Message (GM) channels.

The plugin adds sender-side message status indicators:

- `✓` — message was sent.
- `✓✓` in gray — at least one recipient delivered the message to an active Mattermost web client.
- `✓✓` in blue — at least one recipient kept the message visibly open long enough to count as read.

## Compatibility

| Requirement | Version / Notes |
| --- | --- |
| Mattermost Server | `10.0.0` or newer |
| Plugin OS/architecture | Linux `amd64` bundle is produced by the default Makefile |
| Conversations | Direct Messages, Group Messages, or both, based on plugin configuration |
| Client support | Mattermost web app with plugin support enabled |

## What to install before using this plugin

### For Mattermost administrators

If you only want to install the finished plugin bundle, you need:

1. A running Mattermost `10.x` server.
2. System Console access with permission to upload and enable plugins.
3. Plugin uploads enabled in Mattermost.
4. The built plugin archive: `dist/readreceipt.tar.gz`.

You do **not** need Go, Node.js, npm, or Make on the Mattermost server if someone else already built the archive for you.

### For developers or anyone building the plugin

Install these tools before running `make dist`:

1. **Go 1.22 or newer** — builds the server-side plugin executable.
2. **Node.js 18 or newer** — runs the webapp build toolchain.
3. **npm** — installs the webapp dependencies declared in `webapp/package.json`.
4. **GNU Make** — runs the build targets in `Makefile`.
5. **tar** — packages the plugin bundle.

Optional but useful:

- A local Mattermost development server for manual testing.
- Git for version control and release tagging.

## Features

- Adds sent, delivered, and read indicators to DM and GM posts.
- Detects reads only after a post is at least 70% visible for 2 continuous seconds.
- Uses `IntersectionObserver` in the webapp instead of marking every rendered message as read immediately.
- Records delivery/read acknowledgements through authenticated plugin REST endpoints.
- Persists receipt state in the Mattermost plugin KV store.
- Sends WebSocket events so sender UIs refresh without a page reload.
- Provides admin settings to enable receipts for DMs, GMs, or both.
- Restricts receipt status lookup to the original sender of a message.

## How it works

1. The webapp component observes supported DM/GM posts as they appear in the message list.
2. Recipient clients acknowledge delivery when a supported post is rendered.
3. Recipient clients acknowledge read status only after the visibility threshold and dwell time are met.
4. The server validates the authenticated user, channel type, channel membership, and sender permissions.
5. Receipt state is stored in the plugin KV store and pushed to senders through plugin WebSocket events.

## Build

From the repository root, run:

```bash
make dist
```

The command performs the full release build:

1. Builds the Go server plugin into `server/dist/plugin-linux-amd64`.
2. Installs webapp npm dependencies.
3. Builds the Mattermost webapp bundle into `webapp/dist/main.js`.
4. Packages the plugin into `dist/readreceipt.tar.gz`.

## Install in Mattermost

1. Build the archive with `make dist` or obtain a release archive.
2. In Mattermost, go to **System Console → Plugin Management → Upload Plugin**.
3. Upload `dist/readreceipt.tar.gz`.
4. Enable the plugin.
5. Go to the plugin settings and turn on **Enable Read Receipts**.
6. Choose which conversation types are supported with **Enable For**:
   - `DM`
   - `GM`
   - `DM_AND_GM`

## Configuration

The plugin is disabled by default.

| Setting | Default | Description |
| --- | --- | --- |
| Enable Read Receipts | `false` | Master switch for delivery/read receipt collection and display. |
| Enable For | `DM_AND_GM` | Selects whether receipts apply to DMs, GMs, or both. |

## Development commands

```bash
# Build the complete plugin archive
make dist

# Build only the server executable
make server

# Build only the webapp bundle
make webapp

# Remove generated build artifacts
make clean
```

## Project structure

```text
.
├── assets/                 # Plugin icon and static assets
├── server/                 # Go plugin server, REST API, storage, and WebSocket logic
├── webapp/                 # TypeScript/React Mattermost webapp plugin
├── Makefile                # Build and packaging targets
├── plugin.json             # Mattermost plugin manifest and settings schema
└── README.md
```

## Security model

The server-side API validates every write and read request. It checks that:

- The caller is authenticated.
- The target channel is a supported DM or GM channel.
- The acknowledging user is a member of the target channel.
- Receipt status is only returned to the sender of the post.
- Receipt updates are stored server-side instead of trusting sender-controlled UI state.

## Known limitations

- The default build targets Linux `amd64` only. Adjust `GOOS`, `GOARCH`, and the manifest if you need another platform.
- Delivery means a recipient's web client acknowledged the post; it does not necessarily mean the person saw it.
- Read receipts currently focus on Mattermost web clients. Other clients may not emit the same webapp acknowledgements.

## Troubleshooting

- **Upload fails:** confirm the server version is Mattermost `10.0.0` or newer and plugin uploads are enabled.
- **No indicators appear:** confirm the plugin is enabled and **Enable Read Receipts** is turned on.
- **Receipts do not update:** check that the conversation type is enabled in **Enable For** and that users are in a DM or GM channel.
- **Build fails during webapp build:** run `cd webapp && npm install` and verify Node.js/npm versions.
- **Build fails during server build:** verify Go `1.22+` is installed and available on `PATH`.

## License

Add your project license here before publishing a release.
