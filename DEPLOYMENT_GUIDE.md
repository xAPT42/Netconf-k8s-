# Deployment Guide - netconf-k8s

This guide will walk you through the complete deployment process for the netconf-k8s project.

## Phase 1: Initial Setup and Clone

First, clone the existing repository and set up your local environment.

```bash
cd /home/kali/projets_plan
git clone https://github.com/xAPT42/netconf-k8s.git
cd netconf-k8s
```

## Phase 2: Configure Git Identity

**CRITICAL**: Ensure Git is configured with YOUR credentials before making any commits.

```bash
git config user.name "Your Name"
git config user.email "your.email@example.com"
```

Verify the configuration:

```bash
git config user.name
git config user.email
```

## Phase 3: Create Feature Branch

Create a new branch for the initial project setup:

```bash
git checkout -b feat/initial-setup
```

## Phase 4: Copy Generated Files

All the project files have been generated in `/home/kali/projets_plan/netconf-k8s/`.

Verify the structure:

```bash
tree -L 3 -a
```

Expected output:
```
.
├── .github/
│   └── workflows/
│       └── ci-cd.yml
├── cmd/
│   └── main.go
├── docs/
│   └── architecture.md
├── k8s/
│   ├── checker-cronjob.yaml
│   └── router-deployment.yaml
├── Dockerfile
├── go.mod
├── go.sum
├── README.md
└── DEPLOYMENT_GUIDE.md
```

## Phase 5: Commit and Push

### Step 1: Stage all files

```bash
git add .
```

### Step 2: Verify what will be committed

```bash
git status
```

### Step 3: Create the commit

```bash
git commit -m "feat: add initial project structure, code and CI/CD pipeline

This commit includes:
- Complete README with architecture diagrams and setup instructions
- Detailed architecture documentation with placeholders for screenshots
- Kubernetes manifests for router deployment and compliance checker CronJob
- Multi-stage Dockerfile for optimized container builds
- Go application implementing NETCONF client for compliance checking
- GitHub Actions workflow for automated CI/CD to GCP
- Full project documentation and deployment guide

The project demonstrates cloud-native network automation using:
- NETCONF protocol for network device communication
- Kubernetes for container orchestration
- GCP (GKE + Artifact Registry) for cloud deployment
- GitHub Actions for continuous integration and deployment"
```

### Step 4: Push to GitHub

```bash
git push origin feat/initial-setup
```

## Phase 6: Create Pull Request

1. Navigate to your GitHub repository: `https://github.com/xAPT42/netconf-k8s`

2. You should see a yellow banner suggesting to create a Pull Request for `feat/initial-setup`

3. Click **"Compare & pull request"**

4. Fill in the Pull Request details:

   **Title:**
   ```
   Initial Project Setup - Complete NETCONF K8s Implementation
   ```

   **Description:**
   ```markdown
   ## Summary

   This PR establishes the complete foundation for the netconf-k8s project, a cloud-native network compliance monitoring system.

   ## What's Included

   - ✅ Complete project documentation (README + architecture docs)
   - ✅ Kubernetes manifests (Deployment, Service, CronJob)
   - ✅ Go application for NETCONF-based compliance checking
   - ✅ Multi-stage Dockerfile for optimized builds
   - ✅ GitHub Actions CI/CD pipeline for GCP deployment
   - ✅ Full project structure and dependencies

   ## Architecture Highlights

   - **Network Automation**: NETCONF protocol client in Go
   - **Container Orchestration**: Kubernetes CronJob for scheduled checks
   - **Cloud Deployment**: GKE cluster with Artifact Registry
   - **CI/CD**: Automated build, push, and deploy on every commit

   ## Next Steps After Merge

   1. Configure GCP secrets in GitHub repository settings
   2. Upload project logo to `assets/logo.png`
   3. Deploy to GCP and capture screenshots for documentation
   4. Update `docs/architecture.md` with live deployment evidence

   ## Testing Plan

   - [ ] Verify all files are present and correctly structured
   - [ ] Review Go code for NETCONF implementation
   - [ ] Validate Kubernetes manifests syntax
   - [ ] Check GitHub Actions workflow configuration
   - [ ] Ensure documentation is clear and complete
   ```

5. Click **"Create pull request"**

6. Review the changes one final time

7. Click **"Merge pull request"** → **"Confirm merge"**

8. Optionally, delete the `feat/initial-setup` branch after merge

## Phase 7: Post-Merge Configuration

### Step 1: Update local main branch

```bash
git checkout main
git pull origin main
```

### Step 2: Configure GitHub Secrets

Navigate to: `Settings` → `Secrets and variables` → `Actions` → `New repository secret`

Add the following secrets:

**Secret 1: GCP_PROJECT_ID**
- Name: `GCP_PROJECT_ID`
- Value: Your GCP project ID (e.g., `my-netconf-project`)

**Secret 2: GCP_SA_KEY**
- Name: `GCP_SA_KEY`
- Value: Contents of your service account JSON key file

To get the service account key:
```bash
gcloud iam service-accounts keys create key.json \
  --iam-account=github-actions-sa@YOUR_PROJECT_ID.iam.gserviceaccount.com

cat key.json
```

Copy the entire JSON content and paste it as the secret value.

### Step 3: Create Assets Directory and Add Logo

```bash
mkdir -p assets
```

Place your logo file in `assets/logo.png`. Recommended size: 300x300px, transparent PNG.

Commit and push:
```bash
git add assets/logo.png
git commit -m "docs: add project logo"
git push origin main
```

### Step 4: Update Placeholders in Files

You'll need to replace the following placeholders:

**In README.md:**
- Replace `xAPT42` with your actual GitHub username (3 occurrences)

**In go.mod:**
- Replace `xAPT42` with your actual GitHub username

**In k8s/checker-cronjob.yaml:**
- The image path will be automatically updated by the CI/CD pipeline, but verify the GCP project ID

Use this command to check for placeholders:
```bash
grep -r "xAPT42" .
grep -r "GCP_PROJECT_ID" .
```

Replace them:
```bash
# For macOS
find . -type f -name "*.md" -exec sed -i '' 's/xAPT42/your-github-username/g' {} +
find . -type f -name "go.mod" -exec sed -i '' 's/xAPT42/your-github-username/g' {} +

# For Linux
find . -type f -name "*.md" -exec sed -i 's/xAPT42/your-github-username/g' {} +
find . -type f -name "go.mod" -exec sed -i 's/xAPT42/your-github-username/g' {} +
```

Commit the changes:
```bash
git add .
git commit -m "docs: update repository URLs with actual username"
git push origin main
```

### Step 5: Verify CI/CD Pipeline

After pushing to `main`, the GitHub Actions workflow should trigger automatically.

Monitor the workflow:
1. Go to the `Actions` tab in your GitHub repository
2. Click on the latest workflow run
3. Watch the build and deploy jobs execute

Expected timeline:
- Build and Push: ~3-5 minutes
- Deploy: ~2-3 minutes

## Phase 8: Capture Screenshots

Once the deployment is successful, capture the following screenshots:

### GCP Screenshots

1. **GKE Cluster Dashboard**
   - Navigate to: GCP Console → Kubernetes Engine → Clusters
   - Screenshot: Cluster overview showing nodes and resources

2. **Workloads**
   - Navigate to: GCP Console → Kubernetes Engine → Workloads
   - Screenshot: List of deployments and cronjobs

3. **Artifact Registry**
   - Navigate to: GCP Console → Artifact Registry → Repositories
   - Screenshot: Docker images with tags

### Kubernetes Screenshots

4. **CronJob Details**
   - Click on the CronJob → Screenshot the execution history

5. **Pod Logs - Success**
   - Click on a successful job pod → View logs
   - Screenshot showing PASS status

6. **Pod Logs - Failure** (optional)
   - Manually create a failure scenario
   - Screenshot showing FAIL status

### GitHub Screenshots

7. **GitHub Actions Success**
   - Screenshot of successful workflow run

8. **GitHub Actions Logs**
   - Screenshot of build or deploy job logs

### Step 6: Update Architecture Documentation

Edit `docs/architecture.md` and replace the placeholder text with actual screenshots:

```bash
vim docs/architecture.md
# or
code docs/architecture.md
```

For each section, replace:
```markdown
**Placeholder:** `[Screenshot: description]`
```

With:
```markdown
![Description](../assets/screenshots/screenshot-name.png)
```

Place your screenshots in `assets/screenshots/` directory:

```bash
mkdir -p assets/screenshots
# Copy your screenshots here
```

Commit the documentation updates:
```bash
git add docs/architecture.md assets/screenshots/
git commit -m "docs: add deployment screenshots and evidence"
git push origin main
```

## Phase 9: Final Verification

### Checklist

- [ ] Repository has all files correctly structured
- [ ] README displays correctly with logo and badges
- [ ] GitHub Actions workflow runs successfully
- [ ] GKE cluster is running with all workloads
- [ ] CronJob executes every 5 minutes
- [ ] Compliance checker logs show expected output
- [ ] Architecture documentation includes screenshots
- [ ] All placeholders have been replaced
- [ ] Project is ready to be pinned on GitHub profile

### Testing the Deployment

SSH into a checker pod to test manually:

```bash
kubectl get pods -l app=netconf-checker
kubectl logs <pod-name>
```

Trigger a manual job:

```bash
kubectl create job manual-check-1 --from=cronjob/netconf-checker-cronjob
kubectl logs job/manual-check-1
```

View all resources:

```bash
kubectl get all
```

## Phase 10: Pin Repository on GitHub

1. Go to your GitHub profile
2. Click **"Customize your pins"**
3. Select `netconf-k8s`
4. Save

Your impressive DevOps/Network Automation project is now showcased on your profile!

## Troubleshooting

### Issue: GitHub Actions fails with authentication error

**Solution:**
- Verify `GCP_SA_KEY` secret is correctly set
- Ensure service account has proper IAM roles
- Check service account key is not expired

### Issue: CronJob pods fail to connect to router

**Solution:**
```bash
kubectl logs <checker-pod-name>
kubectl get svc netconf-router-service
kubectl get endpoints netconf-router-service
```

Verify the service is routing to healthy pods.

### Issue: Router pods not ready

**Solution:**
```bash
kubectl describe pod <router-pod-name>
kubectl logs <router-pod-name>
```

Check for image pull errors or resource constraints.

### Issue: Deployment not updating with new image

**Solution:**
```bash
kubectl rollout restart deployment/netconf-router-deployment
kubectl delete cronjob/netconf-checker-cronjob
kubectl apply -f k8s/checker-cronjob.yaml
```

## Maintenance Commands

### View Recent Job Executions
```bash
kubectl get jobs --sort-by=.metadata.creationTimestamp
```

### View CronJob Schedule
```bash
kubectl get cronjob netconf-checker-cronjob
```

### Manually Trigger Compliance Check
```bash
kubectl create job manual-check-$(date +%s) --from=cronjob/netconf-checker-cronjob
```

### View All Logs
```bash
kubectl logs -l app=netconf-router --tail=50
kubectl logs -l app=netconf-checker --tail=50
```

### Scale Router Deployment
```bash
kubectl scale deployment/netconf-router-deployment --replicas=3
```

### Delete All Resources
```bash
kubectl delete -f k8s/
```

---

**Congratulations!** You've successfully deployed a complete cloud-native network automation project. This demonstrates mastery of:

- Network automation (NETCONF)
- Container orchestration (Kubernetes)
- Cloud platforms (GCP)
- CI/CD pipelines (GitHub Actions)
- Infrastructure as Code
- DevOps best practices

This project showcases real-world skills that are highly valued in modern infrastructure and network engineering roles.
