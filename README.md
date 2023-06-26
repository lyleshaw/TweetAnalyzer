# Cyber Fortune Teller(TweetAnalyzer)

## Description

Grab tweets, call Claude or GPT3.5 for analysis and make comments on tweeters

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

## Contribution

### Commit

All commit message should startwith ```feat: ```/```bugfix:```/```refactor: ```/```docs:```/```style```

### Branch

All branch should startwith ```feature/xxx``` etc.
