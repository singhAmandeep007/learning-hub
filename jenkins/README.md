# Jenkins local setup (macOS + Homebrew)

This document captures the current local Jenkins setup used to run the E2E pipeline in this repository.

## Current setup

- OS: macOS
- Jenkins: `jenkins-lts` via Homebrew services
- Jenkins URL: `http://localhost:8080`
- Docker runtime: Docker Desktop
- Git source for pipeline: local repository path (`file:///Users/<user>/Developer/Projects/learning-hub/v1`) - Pipeline script from SCM -> Git -> Repository URL
- Pipeline file: `jenkins/Jenkinsfile.e2e` -> Script path
- E2E compose file: `docker-compose.e2e.yml`

> Note: host port `8080` is used by Jenkins, so Firebase emulator is mapped as `8081:8080` in `docker-compose.e2e.yml`.

## 1) Install and start prerequisites

```bash
brew install jenkins-lts
brew services start jenkins-lts
```

Install and start Docker Desktop, then verify:

```bash
docker version
docker compose version
docker info
```

## 2) Allow local Git checkout in Jenkins

When Jenkins uses `file://` Git URLs, enable local checkout policy:

```bash
launchctl setenv JAVA_TOOL_OPTIONS "-Dhudson.plugins.git.GitSCM.ALLOW_LOCAL_CHECKOUT=true"
brew services restart jenkins-lts
```

You can verify in Jenkins Script Console:

```groovy
println System.getProperty("hudson.plugins.git.GitSCM.ALLOW_LOCAL_CHECKOUT")
```

Expected output: `true`

## 3) Create the Jenkins Pipeline job

In Jenkins UI:

1. **New Item** → **Pipeline**
2. Configure:
	 - Definition: `Pipeline script from SCM`
	 - SCM: `Git`
	 - Repository URL: `file:///Users/<user>/Developer/Projects/learning-hub/v1`
	 - Branch Specifier: `*/main`
	 - Script Path: `jenkins/Jenkinsfile.e2e`
3. Save and click **Build Now**

## 4) What the pipeline does

- Checks out source from local Git URL
- Runs E2E stack:

```bash
docker compose -f docker-compose.e2e.yml up --build --abort-on-container-exit --exit-code-from e2e
```

- Copies Playwright artifacts from the `e2e` service:
	- `playwright-report`
	- `test-results`
- Archives artifacts as `e2e-artifacts/**` in Jenkins build
- Tears down compose stack

## 5) Where to see reports

In Jenkins build page:

- Open **Artifacts**
- Navigate to:
	- `e2e-artifacts/playwright-report/index.html`
	- `e2e-artifacts/test-results/`

## 6) Common issues and fixes

### A) Local Git checkout blocked

Error mentions:

`GitSCM.ALLOW_LOCAL_CHECKOUT`

Fix: run step (2) above and restart Jenkins service.

### B) `docker: command not found`

`jenkins/Jenkinsfile.e2e` exports a PATH that includes Docker Desktop and common bin locations.

If still failing, confirm Docker Desktop is installed and running, then restart Jenkins:

```bash
brew services restart jenkins-lts
```

### C) Port conflict on `8080`

Jenkins uses `localhost:8080`, so Firebase emulator is exposed on `8081` in E2E compose.

### D) Container name conflict

If you see conflicts from earlier runs, clean once:

```bash
docker compose -f docker-compose.e2e.yml down -v --remove-orphans
```

### E) `no space left on device`

Prune Docker data:

```bash
docker builder prune -af
docker image prune -af
docker volume prune -f
```

Also increase Docker Desktop disk image size if needed.

## 7) Current file references

- Pipeline: `jenkins/Jenkinsfile.e2e`
- Compose stack: `docker-compose.e2e.yml`
- E2E config: `e2e/playwright.config.ts`
