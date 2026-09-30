# Rudder

A mobile-first developer cockpit: a lightweight companion daemon (`rudder`) written in Go running on a remote machine or laptop, paired with a React Native (Expo) mobile app to steer, supervise, and review uncommitted code changes from AI agents (Aider-first).

## Architecture

- `daemon/`: Go backend daemon (HTTP REST, WebSocket, Git engine, Aider supervisor, SQLite persistence).
- `app/`: React Native (Expo) mobile client (swipe diff review, live steering, lock-screen notifications).

## Quick Start

### Build & Run Daemon

```bash
make build
make run
```

### Generate TypeScript Types

```bash
make gen-types
```
