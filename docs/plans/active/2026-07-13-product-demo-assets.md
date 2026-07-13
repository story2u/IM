# Product Demo Assets

> Status: active
> Branch: `features/product-demo-assets`
> Baseline: `release/v2.0.0` at `9e43250`

## Goal

Deliver a truthful unauthenticated product home, deterministic Web demo, repeatable screenshots,
Playwright recording, Remotion presentation, FFmpeg derivatives, mobile screenshot harnesses, and
release documentation without exposing production data or credentials.

## Acceptance

- Unauthenticated `/` is a responsive product home; authenticated `/` remains the real dashboard.
- `/demo/*` is available only when `DEMO_MODE=true` and uses isolated fictional data.
- Playwright produces the documented Web screenshots and a deterministic raw WebM recording.
- Remotion and FFmpeg produce validated MP4, WebM, GIF and cover outputs without narration.
- iOS/Android Debug-only demo routes and screenshot scripts fail clearly without platform tools.
- README and release docs distinguish implemented, Beta, and externally unverified capabilities.

## Work

1. Audit routes, auth, data, mobile projects, workflows and host media tools.
2. Build isolated demo data/routes and unauthenticated home.
3. Add screenshot and recording automation with stable selectors.
4. Add Remotion timeline and FFmpeg post-processing.
5. Add Debug-only mobile demo routes and capture scripts.
6. Update README, docs, workflow and checks; generate and inspect actual Web assets.

## Safety Decisions

- Demo data uses fixed IDs, `example.com` addresses, fictional groups and a fixed clock.
- Demo routes never initialize OAuth, RevenueCat, Telegram or mutation APIs.
- Production builds fail closed because both server and client demo gates default to false.
- Generated binary video is a build artifact; only lightweight screenshots/GIF/cover are candidates
  for normal Git history. Full MP4/WebM are prepared for GitHub Release.

## Verification Log

- Pending.
