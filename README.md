# Cyber Fortune Teller(TweetAnalyzer)

## Description

Grab tweets, call Claude or GPT3.5 for analysis and make comments on tweeters

<div align="center">
    &nbsp;
    <a href="https://t.me/tweet_analyzer"><img src="https://img.shields.io/badge/-Telegram-red?style=social&logo=telegram" height=25></a>
    &nbsp;
    <a href="https://twitter.com/L_x_x_x_x_x"><img src="https://img.shields.io/badge/-Twitter-red?style=social&logo=twitter" height=25></a>
    &nbsp;
</div>

### Built With

- [Next.js](https://nextjs.org/?ref=cal.com)
- [React.js](https://reactjs.org/?ref=cal.com)
- [Tailwind CSS](https://tailwindcss.com/?ref=cal.com)
- [Daisy UI](https://daisyui.com/)
- [Go](https://go.dev/)

## Getting Started

To get a local copy up and running, please follow these simple steps.

### Prerequisites

Here is what you need to be able to run Cal.com.

- Node.js (Version: >=18.x)
- Go (Version: >= 1.18)
- Pnpm ((recommended))

## Development

### Setup

1. Clone the repo into a public GitHub repository (or fork https://github.com/lyleshaw/TweetAnalyzer/fork).

    ```bash
    git clone https://github.com/lyleshaw/TweetAnalyzer.git
    ```
2. Go to the project folder
    ```bash
    cd TweetAnalyzer
    ```
3. Install packages with pnpm
4. Set up your .env file
   1. Duplicate .env.example to .env
      ```bash
      cp .env.example .env
      ```
5. Run (in development mode)
   - frontend
        ```bash
        pnpm  dev
        ```
   - rearend
        ```bash
        go mod tidy
        go run ./service/index.go
        ```

## Deployment

### Vercel

You can use Vercel to deploy the front end, just click the button below


[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Flyleshaw%2FTweetAnalyzer)

## Contribution

### Commit

All commit message should startwith ```feat: ```/```bugfix:```/```refactor: ```/```docs:```/```style```

### Branch

All branch should startwith ```feature/xxx``` etc.
