# ❤️ FindMe

> Find your love with verified people and authentic video profiles - no filters, just real connections.

[![License: Commercial](https://img.shields.io/badge/License-Commercial-red.svg)](LICENSE)
[![React Native](https://img.shields.io/badge/React%20Native-0.73+-blue.svg)](https://reactnative.dev/)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://golang.org/)
[![Qdrant](https://img.shields.io/badge/Qdrant-Vector%20DB-purple.svg)](https://qdrant.tech/)
[![Status](https://img.shields.io/badge/Status-In%20Development-yellow.svg)]()

## 🌟 About

**FindMe** is a revolutionary AI-powered dating platform that brings authenticity back to online dating. Through verified video profiles and intelligent matching, we help you connect with real people looking for meaningful relationships.

### The FindMe Difference

- **🎥 Video-First Profiles**: No photos, no filters - just authentic video introductions
- **✅ Verified Users**: Every profile is verified to ensure genuine connections
- **🤖 AI-Powered Matching**: Our AI analyzes your video presentation and recommends 4 compatible matches weekly
- **📹 Video Call Connections**: Get to know your match through face-to-face video calls
- **⏰ Weekly Decision System**: At the end of each week, decide if you want to continue or explore new connections
- **🆓 Free & Open Source**: Built by the community, for the community

## 🏗️ Architecture

FindMe is built with a modern, scalable architecture:

```
┌─────────────────┐
│  Mobile App     │
│  (React Native  │
│   + TypeScript) │
└────────┬────────┘
         │
         │ REST API
         │
┌────────▼────────┐       ┌──────────────┐
│  Backend API    │◄──────┤   Qdrant     │
│  (Go + Gin)     │       │  Vector DB   │
└─────────────────┘       └──────────────┘
         │
         │
┌────────▼────────┐
│  Video Storage  │
│  & Processing   │
└─────────────────┘
```

### Tech Stack

**Mobile Application:**
- React Native with TypeScript
- React Navigation for routing
- WebRTC for video calls
- Redux Toolkit for state management

**Backend Services:**
- Go 1.21+ with Gin framework
- JWT authentication
- RESTful API design
- WebSocket support for real-time features

**Database & Storage:**
- Qdrant vector database for AI matching
- PostgreSQL for user data
- S3-compatible storage for videos
- Redis for caching and sessions

**AI & Machine Learning:**
- Vector embeddings for profile matching
- Video analysis for verification
- Recommendation engine

## 🚀 Quick Start

### Prerequisites

- Node.js 20+ and npm
- Go 1.21+
- Docker and Docker Compose
- React Native development environment
- iOS Simulator (Mac) or Android Emulator

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/findme.git
   cd findme
   ```

2. **Set up environment variables**
   ```bash
   cp .env.sample .env
   # Edit .env with your configuration
   ```

3. **Start the backend services**
   ```bash
   cd backend
   docker-compose up -d  # Starts Qdrant, PostgreSQL, Redis
   go run main.go
   ```

4. **Start the mobile app**
   ```bash
   cd mobile
   npm install
   npm run ios     # For iOS
   npm run android # For Android
   ```

For detailed setup instructions, see [SETUP.md](docs/SETUP.md).

## 📚 Documentation

- **[Setup Guide](docs/SETUP.md)** - Complete development environment setup
- **[Architecture](ARCHITECTURE.md)** - System design and technical architecture
- **[API Documentation](docs/API.md)** - Backend API reference
- **[Deployment Guide](docs/DEPLOYMENT.md)** - Production deployment instructions
- **[Migration Guide](docs/MIGRATION.md)** - Database migration and schema management
- **[Contributing](CONTRIBUTING.md)** - How to contribute to FindMe

## 🤝 Contributing

We welcome contributions from the community! FindMe is free and open source, and we believe in building great software together.

Please read our [Contributing Guidelines](CONTRIBUTING.md) to get started.

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Make your changes with emoji commits 🎨
4. Write tests for new features
5. Submit a pull request

## 🎯 Roadmap

- [x] Project initialization and documentation
- [ ] Core backend API with user authentication
- [ ] Qdrant integration for vector search
- [ ] Mobile app UI/UX implementation
- [ ] Video upload and verification system
- [ ] AI recommendation engine
- [ ] WebRTC video call integration
- [ ] Weekly matching cycle implementation
- [ ] Beta testing program
- [ ] Production deployment

## 📄 License

This project is licensed under a Commercial License - see the [LICENSE](LICENSE) file for details.

**Free to use for personal and commercial purposes** with attribution. The source code is open for learning, modification, and distribution while maintaining the commercial license terms.

## 🌐 Community

- **Website**: [findme.ai](https://findme.ai)
- **Issues**: [GitHub Issues](https://github.com/yourusername/findme/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/findme/discussions)

## 💖 Acknowledgments

Built with love by the FindMe community. Special thanks to all our contributors who believe in authentic connections.

---

**Made with ❤️ by the FindMe Team**
