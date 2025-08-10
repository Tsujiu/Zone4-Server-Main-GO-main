# 🎮 Zone4-Server-Main-GO

Main Game Server for **Zone4** (written in Go, previously disguised as MU Online server for safety).

> ⚠️ **Project in Development** – Not 100% complete yet. Help us finish it!

---

## 📖 About the Project
- Main server for the **Zone4** game, recreated in **Go**.
- Goal: **revive the original Zone4 experience** for the community.
- **Non-commercial** and nonprofit initiative.
- Open to contributions: code, testing, ideas, and improvements are welcome.
- **Hosting is ready** for testing! Contact via GitHub.

---

## 🆕 Important Update
- The **Channel Manager** is now **integrated directly into the server code**.
- No need to download or run an external manager.
- Server startup now handles channels, status, and processes internally.

---

## ✨ Features
- Core game logic and networking for Zone4 (in Go)
- Compatible with [Zone4 LandVerse Client](https://github.com/Tsujiu/Zone4-LandVerse-Client-and-PDB-and-Debug)
- Configuration files, startup scripts, and Docker support
- Uses **Git LFS** for large files
- Modular structure with native **Manager** integration

---

## ⚙️ Default Credentials & Configuration

**Microsoft SQL Server**
```
Host:     103.208.24.115
Port:     1433
User:     muvl
Password: muvl123
Database: game
```

**Redis**
```
Host:     103.208.24.171
Port:     6379
User:     admin
Password: MuRedisP@ssw0rd
```

**Environment Variables (.env)**
```env
# ===== Connections =====
REDIS_ADDR=103.208.24.171:6379
REDIS_USER=admin
REDIS_PASS=MuRedisP@ssw0rd

SQLSERVER_GAME=sqlserver://muvl:muvl123@103.208.24.115:1433?database=game
SQLSERVER_GAME_TEST=sqlserver://muvl:muvl123@localhost:1433?database=game
SQLSERVER_GAME_INVENTORY=sqlserver://muvl:muvl123@103.208.24.115:1433?database=game_inventory
SQLSERVER_GAME_INVENTORT_TEST=sqlserver://muvl:muvl123@localhost:1433?database=game_inventory

# ===== Encryption =====
AES_KEY=p*{Ilqw<8AT_@poI2Kq3D1uVcp`*@bRh
AES_IV=EB484700C94AB2CF81B8C05B324ED164

# ===== Ports =====
TCP_PORT=9090
UDP_PORT=9091
PPROF_PORT=6060
CHANNEL_PORT=9090
CHANNEL_PORTS=29998,29996,29995,29993,29994,29992
CHANNEL_INDEX=0

# ===== Other =====
WORKER_POOL=10
MAX_MATCHING_WORKER=10

DB_USER=muvl
DB_PASSWORD=muvl123
DB_NAME=game
DB_SERVER=103.208.24.115
DB_PORT=1433
MAX_USER_CHANNEL=500
```

---

## 🚀 Getting Started

**Requirements:**
- Go 1.19+
- Docker (optional, via `docker-compose`)
- Redis + Microsoft SQL Server (for persistence – see `.env`)

**Run locally:**
```bash
go run cmd/server/main.go
```

---

## 🤝 Contributing
- Open a **Pull Request**
- Create an **Issue**
- Join the development

---

## 📜 Changelog
> Summary of key changes. For full details, see [CHANGELOG.md](CHANGELOG.md).

### v1.0.0 – Initial Release
- Base structure with Go and Docker support
- TCP/UDP connection handling
- Modules: user, character, inventory, skills, gameplay, gangwar
- Item, monster, and `.att` map data
- **Native Manager integration**, removing external dependency

---

## 📄 License
Licensed under **MIT**.

> This project is not affiliated with the official publishers of Zone4. It is a nonprofit initiative by fans, for fans.
