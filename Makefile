.PHONY: harness-check backend-sync backend-check pi-agent-sync pi-agent-check frontend-check ios-check android-check demo-web demo-screenshots demo-record demo-video demo-assets demo-ios-screenshots demo-android-screenshots check

PYTHON ?= python3
UV ?= uv
PNPM ?= corepack pnpm@10.25.0

harness-check:
	$(PYTHON) scripts/harness_check.py
	$(PYTHON) -m unittest discover -s scripts/tests -p 'test_*.py'

backend-sync:
	cd backend && $(UV) sync --locked --dev

backend-check: backend-sync
	cd backend && $(UV) run --locked python -m compileall -q app tests scripts alembic
	cd backend && $(UV) run --locked ruff check app tests scripts alembic --select E,F,ASYNC --ignore E501
	cd backend && $(UV) run --locked pytest -q

pi-agent-sync:
	cd backend/pi-agent-runtime && npm ci --ignore-scripts

pi-agent-check: pi-agent-sync
	cd backend/pi-agent-runtime && npm run check
	cd backend/pi-agent-runtime && npm test

frontend-check:
	cd frontend && $(PNPM) lint
	cd frontend && $(PNPM) typecheck
	cd frontend && $(PNPM) test
	cd frontend && $(PNPM) build

# 需要 macOS + Xcode + xcodegen（brew install xcodegen）；CI 用 macOS runner。
ios-check:
	cd mobile/ios && xcodegen generate
	cd mobile/ios && xcodebuild test -project OpportunityRadar.xcodeproj -scheme OpportunityRadar \
		-destination 'platform=iOS Simulator,OS=latest,name=iPhone 16' \
		-derivedDataPath .build/DerivedData CODE_SIGNING_ALLOWED=NO

# 需要 JDK 17 + Android SDK；首次运行 `cd mobile/android && gradle wrapper` 生成 wrapper。
android-check:
	cd mobile/android && ./gradlew --no-daemon lintDebug testDebugUnitTest assembleDebug

demo-web:
	cd frontend && DEMO_MODE=true NEXT_PUBLIC_DEMO_MODE=true $(PNPM) dev --hostname 127.0.0.1 --port 3100

demo-screenshots:
	cd frontend && $(PNPM) demo:screenshots

demo-record:
	cd frontend && $(PNPM) demo:record

demo-video:
	bash scripts/demo/render-demo.sh
	bash scripts/demo/render-webm.sh
	bash scripts/demo/render-gif.sh
	bash scripts/demo/generate-cover.sh

demo-assets: demo-screenshots demo-record demo-video
	bash scripts/demo/check-assets.sh

demo-ios-screenshots:
	bash scripts/demo/capture-ios.sh

demo-android-screenshots:
	bash scripts/demo/capture-android.sh

check: harness-check backend-check pi-agent-check frontend-check
