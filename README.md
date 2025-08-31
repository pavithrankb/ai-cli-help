
# 🚀 CLI-Help: AI-Powered Command Line Assistant

**Author:** Pavithran KB
**Version:** 0.5.0  

`cli-help` is an AI-powered CLI assistant that translates **natural language** into real shell commands.  
It integrates with **Amazon Bedrock (Claude models)** using a Bearer token.  

---

## ✨ Features

- ✅ **Safe Mode (default):** Blocks dangerous commands like `rm`  
- ⚠️ **Unsafe Mode (`--unsafe`):** Allows *all* commands (use carefully!)  
- 🤖 **AI Mode (`--ai`):** Translate natural language → shell commands via Bedrock Claude  
- 🔍 **Verbose Mode (`--verbose`):** Show raw JSON responses from Bedrock  
- 🔄 **Custom Models:** Select model via `BEDROCK_MODEL_ID` (defaults to Claude 3 Haiku)  

---

## 🔧 Installation

```bash
# Clone the repo
git clone https://github.com/pavithrankb/ai-cli-help
cd ai-cli-help

# Build
go build -o cli-help main.go

# Make it accessible
sudo mv cli-help /usr/local/bin/
```

---

## 🔑 Setup

Export your AWS Bedrock Bearer token:

```bash
export AWS_BEARER_TOKEN_BEDROCK=your_bearer_token_here
```

Optional: choose which Claude model to use:

```bash
# Default is Claude 3 Haiku
export BEDROCK_MODEL_ID=anthropic.claude-3-sonnet-20240229-v1:0
```

---

## 🖥️ Usage

```bash
cli-help [options]
```

### Options:
- `--help`       Show help message  
- `--version`    Show version info  
- `--unsafe`     Enable **UNSAFE MODE** (dangerous commands allowed)  
- `--ai`         Enable **AI Mode** (natural language → commands)  
- `--verbose`    Print raw AI response JSON  

---

## ⚙️ Modes Explained

| Mode         | Description                                                                 |
|--------------|-----------------------------------------------------------------------------|
| **Safe**     | Default. Blocks `rm` and destructive commands.                              |
| **Unsafe**   | Allows all commands, including `rm`. Use only if you trust your prompt.     |
| **AI**       | Turns natural language into commands using Claude via Bedrock API.          |
| **Verbose**  | Prints Claude's raw JSON output (for debugging).                            |

---

## 💡 Example Use Cases


### 1️⃣ Install packages

```bash
./cli-help --ai
🤖 AI Mode is ON: Type natural language and Ill translate to commands.
> install nginx on ubuntu
🤖 Suggested command:
sudo apt-get update && sudo apt-get install -y nginx
```

---

### 2️⃣ Delete files (requires `--unsafe`)

```bash
 ./cli-help --ai --unsafe     
🤖 AI mode enabled! Using Bedrock Anthropic Claude (Bearer token).
⚠️ UNSAFE MODE enabled! 'rm' commands are now allowed.
Welcome to CLI Helper 🚀
⚠️ Running in UNSAFE MODE. Dangerous commands are allowed.
🤖 AI Mode is ON: Type natural language and Ill translate to commands.
Type your command (or natural language request). Type 'exit' to quit.
> delete all log files older than 7 days
🤖 Suggested command:
find /var/log -type f -name "*.log" -mtime +7 -delete
```
---

### 3️⃣ Create 10 test files with timestamps (for dev testing)

```bash
./cli-help --ai
🤖 AI Mode is ON: Type natural language and Ill translate to commands.
> create 10 files with names file1..file10 and add timestamp inside
🤖 Suggested command:
for i in {1..10}; do echo $(date) > "file${i}_$(date +%s).txt"; done
Run this? (y/n): y
```

👉 Result: `file1_<timestamp>.txt ... file10_<timestamp>.txt`

---

### 4️⃣ Install Kubernetes, deploy Nginx, and delete it

```bash
./cli-help --ai
🤖 AI mode enabled! Using Bedrock Anthropic Claude (Bearer token).
Welcome to CLI Helper 🚀
✅ Running in SAFE MODE. 'rm' commands are blocked.
🤖 AI Mode is ON: Type natural language and Ill translate to commands.
Type your command (or natural language request). Type 'exit' to quit.
> Install K3s with default addons disabled
🤖 Suggested command:
curl -sfL https://get.k3s.io | sh -s - --disable-default-addon=coredns --disable-default-addon=servicelb --disable-default-addon=traefik
Run this? (y/n): y
[INFO]  Finding release for channel stable
[INFO]  Using v1.33.4+k3s1 as release
[INFO]  Downloading hash https://github.com/k3s-io/k3s/releases/download/v1.33.4+k3s1/sha256sum-amd64.txt
.
.
.
[INFO]  systemd: Starting k3s

---

> Set KUBECONFIG environment variable
🤖 Suggested command:
export KUBECONFIG=/path/to/kubeconfig
Run this? (y/n): y

> Deploy Nginx using K3s
🤖 Suggested command:
sudo curl -sfL https://get.k3s.io | sh -
kubectl create deployment nginx --image=nginx
kubectl expose deployment nginx --port=80 --type=NodePort
Run this? (y/n): y
[INFO]  Finding release for channel stable
[INFO]  Using v1.33.4+k3s1 as release
[INFO]  Downloading hash https://github.com/k3s-io/k3s/releases/download/v1.33.4+k3s1/sha256sum-amd64.txt
[INFO]  Skipping binary downloaded, installed k3s matches hash
.
.
.
deployment.apps/nginx created
service/nginx exposed

---

> List all pods status
🤖 Suggested command:
kubectl get pods --all-namespaces
Run this? (y/n): y
NAMESPACE     NAME                                      READY   STATUS             RESTARTS      AGE
default       nginx-5869d7778c-hsb92                    1/1     Running            0             84s
kube-system   coredns-64fd4b4794-tb5qc                  1/1     Running            0             84s
kube-system   helm-install-traefik-9m6bl                1/1     Running            0             84s
kube-system   helm-install-traefik-crd-ntqqs            0/1     Completed          0             85s
kube-system   local-path-provisioner-774c6665dc-92tgc   1/1     Running            0             84s
kube-system   metrics-server-7bfffcd44-dvwbd            1/1     Running            0             84s

---

> Stop and uninstall k3s
🤖 Suggested command:
sudo systemctl stop k3s
sudo systemctl disable k3s
sudo /usr/local/bin/k3s-uninstall.sh
Run this? (y/n): y
Removed "/etc/systemd/system/multi-user.target.wants/k3s.service".
++ id -u
+ '[' 0 -eq 0 ']'
+ K3S_DATA_DIR=/var/lib/rancher/k3s
+ /usr/local/bin/k3s-killall.sh
+ for service in /etc/systemd/system/k3s*.service
+ '[' -s /etc/systemd/system/k3s.service ']'
.
.
.

Removed:
  container-selinux-3:2.233.0-1.amzn2023.noarch                 k3s-selinux-1.6-1.el8.noarch                

Complete!
+ rm -f /etc/yum.repos.d/rancher-k3s-common.repo
+ remove_uninstall
+ rm -f /usr/local/bin/k3s-uninstall.sh
> exit
```

---

## 🧑‍💻 Example Flow

```bash
# Safe mode (default)
cli-help

# Allow destructive commands
cli-help --unsafe

# AI-assisted mode
cli-help --ai

# AI-assisted + Verbose debugging
cli-help --ai --verbose
```

---

## ⚠️ Disclaimer

This tool **executes real shell commands**.  
- By default, it **blocks `rm`** and similar destructive commands.  
- If you enable `--unsafe`, **you are responsible** for what runs.  

Use with care, especially in production environments.
