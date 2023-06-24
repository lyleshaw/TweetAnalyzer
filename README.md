# TweetAnalyzer

## Description

该项目抓取 Twitter 推文内容，调用 Claude 或 GPT3.5 进行分析而对推主做出评价。

## Usage

### 1. 安装依赖

pre-env: go 1.18 + nodejs

- 前端

```bash
npm install
npm run dev
```

- 后端

```bash
go mod tidy
go run ./service/index.go
```

### 2. 配置

环境变量配置：

```text
ANTHROPIC_API_KEY=sk-ant-xxx
PORT=8080
API_KEY=sb-xxx
BASE_URL=https://api.openai-sb.com/v1
```
