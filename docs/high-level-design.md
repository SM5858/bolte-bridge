# Bolte Bridge — High-Level Design

A two-way bridge between an established email **mailing list** and a **Matrix
room**. Messages posted to the list appear in the room; messages sent in the
room are delivered to the list. The goal is *seamless* communication: a
participant on either side should be able to follow and join the conversation
without knowing the other medium exists.

---

## 1. Goals

### Goals

- **Bidirectional message relay** between one mailing list and one Matrix room.
- **Sender fidelity** — a message from `alice@example.org` should look like it
  comes from Alice in Matrix (own display name / avatar), and a Matrix user's
  message should look like it comes from a stable, attributable address on the
  list.
- **Threading fidelity** — email reply chains (`In-Reply-To`/`References`) map
  to Matrix replies/threads and vice versa.
- **No message loops** — the bridge never re-bridges its own traffic.

---

## 2. Architecture Overview

```text
        ┌─────────────────────┐                 ┌──────────────────────┐
        │   Mailing List      │                 │   Matrix Homeserver  │
        │  (SMTP / IMAP)      │                 │  (Client-Server +    │
        │                     │                 │   Application Svc)   │
        └───────┬─────▲───────┘                 └────────▲─────┬───────┘
                │     │                                  │     │
        inbound │     │ outbound                 outbound│     │ inbound
        email   │     │ email                    to room │     │ events
                │     │                                  │     │
        ┌───────▼─────┴──────────────────────────────────┴─────▼────────┐
        │                        Bolte Bridge                           │
        │                                                               │
        │  ┌──────────────┐   ┌───────────────┐   ┌──────────────────┐  │
        │  │ Email Ingest │   │   Core Relay  │   │  Matrix Adapter  │  │
        │  │ (IMAP/LMTP)  │──▶  - normalize  │◀──│ (appservice/     │  │
        │  │              │   │  - identity map│  │  mautrix-go)     │  │
        │  │ Email Egress │◀──  - threading  │──▶│                  │  │
        │  │ (SMTP submit)│   │  - loop guard │   │                  │  │
        │  └──────────────┘   └───────┬───────┘   └──────────────────┘  │
        │                              │                                │
        │                     ┌────────▼────────┐                       │
        │                     │   State Store   │  (message↔event map,  │
        │                     │   (SQLite)      │   identities, threads)│
        │                     └─────────────────┘                       │
        └───────────────────────────────────────────────────────────────┘
```

Three internal components around a shared **state store**:

1. **Email Adapter** — receives list mail and submits outbound mail.
2. **Matrix Adapter** — receives room events and posts on behalf of users.
3. **Core Relay** — the medium-agnostic translation layer: normalization,
   identity mapping, threading, dedup/loop prevention.

The adapters speak in a small internal `Message` model so the core never deals
with raw MIME or raw Matrix events.

### Execution model: one-shot CLI on a systemd timer

The bridge is a CLI that runs to completion on each invocation and exits;
a **systemd timer** triggers it on a fixed interval (e.g. every 1–15 minutes).
Each tick performs one full relay cycle:

1. **Read Matrix** — one-shot `GET /sync?since=<token>&timeout=0` returns
   everything since the stored token immediately; record the new `next_batch`.
2. **Relay room → list** — new room events become outbound email (SMTP submit).
3. **Read email** — IMAP fetch of pending messages since the last UID mark.
4. **Relay list → room** — new mail is posted as ghost users via the appservice
   token.
5. **Persist cursors** — sync token, IMAP UID, and relay statuses committed
   transactionally, then exit.

Both intake paths are **pull-based**. Our latency will be bounded by the
SystemD timer interval and will be acceptable for our asynchronous
mailing list. This design gives us operation simplicity: no long-running
server, and trivial restart and fault-recovery semantics.

Runs must not overlap. A `Type=oneshot` service unit is not double-started by
systemd while a prior run is active, and an advisory `flock` guards the case
where a tick outruns its interval.
