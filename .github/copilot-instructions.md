# Copilot / AI Agent Guidance for this workspace

Summary
- This workspace is a multi-project monorepo: many PHP web apps live under `htdocs/` (Yii1 and Yii2 based), plus a Go skeleton in `mygo/`.
- Primary local dev flow uses XAMPP (Apache+PHP) and Composer; some subprojects use Docker for tests.

High-level architecture notes
- PHP monorepo: multiple independent web apps in `htdocs/` (examples: [htdocs/dkds](htdocs/dkds), [htdocs/webusagi](htdocs/webusagi), [htdocs/webdeka](htdocs/webdeka)). Many follow Yii conventions (see `protected/` or `vendor/yiisoft/yii2`).
- Go project: `mygo/` contains package layout (`mygo/internal`, `mygo/pkg`) but `mygo/cmd` is currently empty — treat this as a skeleton.

Developer workflows (concrete, discoverable)
- Local web: use XAMPP to host apps under `htdocs/` (workspace contains `xampp/`), point Apache document root to the app folder and run PHP as usual.
- Dependency install: run `composer install` in the PHP app root (see [htdocs/webdeka/composer.json](htdocs/webdeka/composer.json)).
- Docker/test: some vendor packages provide `Makefile` / docker targets (example: [htdocs/webusagi/vendor/yiisoft/yii2-gii/Makefile](htdocs/webusagi/vendor/yiisoft/yii2-gii/Makefile)). Inspect project subfolders for `docker-compose.yml` (example: [htdocs/dkds/docker-compose.yml](htdocs/dkds/docker-compose.yml)).
- Database: SQL fixtures / dumps are in the repo (example: [htdocs/dbdatainv.sql](htdocs/dbdatainv.sql)). Check app `protected/config` or `config/` files for DB DSNs when wiring local dev.

Project-specific patterns to follow
- Framework detection: prefer Yii conventions when present — `protected/` (Yii1) or `vendor/yiisoft/yii2` (Yii2). Use these clues to determine routing and bootstrap.
- Controllers and views: most PHP apps follow `controllers/` + `views/` layout (see [htdocs/dkds/controllers](htdocs/dkds/controllers)). When changing controllers, update related views and asset bundles under `web/` or `assets/`.
- Avoid editing `vendor/` or third-party plugin code directly; prefer patches, composer overrides, or forking.

Integration points & external dependencies
- Composer-managed PHP packages (vendor/). Always run `composer install` after pulling changes.
- Many apps expect a local MySQL/Postgres instance (XAMPP or Docker). Look for `config` files in each app root to confirm connection strings.
- Static assets and JS live under each app's `web/` or `web/js` directories — editing may require clearing (or rebuilding) caches.

Editing / PR guidance for agents
- Make minimal, focused edits in the app's folder. Reference app-specific config and routes before changing behavior.
- Add/update tests next to code if the app already contains tests; otherwise, include a short manual test note in the PR description (how to run manually via browser/XAMPP).
- Always run `composer install` and, when present, the project's test command (or vendor Makefile targets) before suggesting changes.

Where to look first (quick links)
- App examples: [htdocs/dkds](htdocs/dkds) — controllers, config, docker-compose.
- Composer example: [htdocs/webdeka/composer.json](htdocs/webdeka/composer.json).
- Vendor/test Makefile example: [htdocs/webusagi/vendor/yiisoft/yii2-gii/Makefile](htdocs/webusagi/vendor/yiisoft/yii2-gii/Makefile).
- DB fixtures: [htdocs/dbdatainv.sql](htdocs/dbdatainv.sql).
- Go skeleton: [mygo/](mygo/) and [mygo/cmd](mygo/cmd) (cmd is currently empty).

If anything is ambiguous
- Ask: which app (path under `htdocs/`) is the target, whether to run changes under XAMPP or Docker, and whether to update shared assets or a single app.

End
