# Contributing to FindMe ❤️

First off, thank you for considering contributing to FindMe! It's people like you that make FindMe such a great tool for helping people find authentic connections.

## 📜 Code of Conduct

### Our Pledge

We are committed to providing a welcoming and inspiring community for all. Please be respectful and constructive in your interactions.

### Our Standards

**Positive behavior includes:**
- Using welcoming and inclusive language
- Being respectful of differing viewpoints and experiences
- Gracefully accepting constructive criticism
- Focusing on what is best for the community
- Showing empathy towards other community members

**Unacceptable behavior includes:**
- Harassment, trolling, or discriminatory comments
- Publishing others' private information without permission
- Other conduct which could reasonably be considered inappropriate

## 🚀 Getting Started

### Prerequisites

Before contributing, make sure you have:

- Read the [README.md](README.md)
- Set up the development environment following [docs/SETUP.md](docs/SETUP.md)
- Familiarized yourself with the [architecture](ARCHITECTURE.md)
- Checked existing [issues](https://github.com/yourusername/findme/issues) and [pull requests](https://github.com/yourusername/findme/pulls)

### First Time Contributors

If you're new to open source, check out issues labeled with `good-first-issue` or `help-wanted`. These are great starting points!

## 🔄 Development Workflow

### 1. Fork & Clone

```bash
# Fork the repository on GitHub, then:
git clone https://github.com/YOUR_USERNAME/findme.git
cd findme
git remote add upstream https://github.com/original/findme.git
```

### 2. Create a Branch

Always create a new branch for your work:

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
# or
git checkout -b docs/your-documentation-update
```

**Branch naming conventions:**
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation changes
- `refactor/` - Code refactoring
- `test/` - Adding or updating tests
- `chore/` - Maintenance tasks

### 3. Make Your Changes

Write clean, maintainable code following our style guides (see below).

### 4. Commit Your Changes

We use emoji commits to make our history more readable and fun! 🎨

```bash
git add .
git commit -m "✨ Add video profile upload feature"
```

**Commit message format:**
```
<emoji> <Short description (50 chars or less)>

<Detailed description if necessary>

<Reference to issue if applicable>
```

**Common commit emojis:**
- ✨ `:sparkles:` - New feature
- 🐛 `:bug:` - Bug fix
- 📝 `:memo:` - Documentation
- 🎨 `:art:` - Code style/formatting
- ♻️ `:recycle:` - Refactoring
- ⚡ `:zap:` - Performance improvement
- ✅ `:white_check_mark:` - Tests
- 🔒 `:lock:` - Security fix
- ⬆️ `:arrow_up:` - Upgrade dependencies
- ⬇️ `:arrow_down:` - Downgrade dependencies
- 🔧 `:wrench:` - Configuration changes
- 🌐 `:globe_with_meridians:` - Internationalization
- 🚀 `:rocket:` - Deployment

### 5. Keep Your Branch Updated

```bash
git fetch upstream
git rebase upstream/main
```

### 6. Push Your Changes

```bash
git push origin feature/your-feature-name
```

### 7. Submit a Pull Request

Go to GitHub and create a pull request from your branch to the main repository.

**Pull Request Guidelines:**
- Fill out the PR template completely
- Link related issues
- Add screenshots/videos for UI changes
- Ensure all tests pass
- Request review from maintainers

## 💻 Code Style Guidelines

### TypeScript / React Native

**General Principles:**
- Use TypeScript for all new code
- Follow functional programming patterns
- Use React hooks (avoid class components)
- Keep components small and focused

**Naming Conventions:**
```typescript
// Components: PascalCase
export const VideoProfile: React.FC<Props> = () => { ... }

// Hooks: camelCase starting with "use"
export const useVideoUpload = () => { ... }

// Constants: UPPER_SNAKE_CASE
const MAX_VIDEO_SIZE = 50 * 1024 * 1024;

// Functions/Variables: camelCase
const handleVideoUpload = async () => { ... }
```

**File Structure:**
```
src/
├── components/       # Reusable UI components
├── screens/         # Screen-level components
├── hooks/           # Custom React hooks
├── services/        # API and business logic
├── types/           # TypeScript type definitions
├── utils/           # Utility functions
└── constants/       # Application constants
```

**Code Examples:**

```typescript
// ✅ Good
interface VideoProfileProps {
  userId: string;
  onUploadComplete: (videoUrl: string) => void;
}

export const VideoProfile: React.FC<VideoProfileProps> = ({
  userId,
  onUploadComplete,
}) => {
  const [isUploading, setIsUploading] = useState(false);
  
  const handleUpload = useCallback(async (file: File) => {
    try {
      setIsUploading(true);
      const url = await uploadVideo(file);
      onUploadComplete(url);
    } catch (error) {
      console.error('Upload failed:', error);
    } finally {
      setIsUploading(false);
    }
  }, [onUploadComplete]);
  
  return <View>{/* UI */}</View>;
};

// ❌ Bad - Using any, class component
class VideoProfile extends React.Component<any, any> {
  // Avoid this pattern
}
```

### Go / Backend

**General Principles:**
- Follow standard Go conventions
- Use `gofmt` and `golint`
- Write idiomatic Go code
- Handle errors explicitly

**Project Structure:**
```
backend/
├── cmd/
│   └── api/         # Application entry point
├── internal/
│   ├── api/         # HTTP handlers (Gin)
│   ├── auth/        # Authentication logic
│   ├── models/      # Data models
│   ├── repository/  # Database layer
│   ├── service/     # Business logic
│   └── util/        # Utilities
├── migrations/      # Database migrations
└── pkg/            # Public packages
```

**Naming Conventions:**
```go
// Exported: PascalCase
type UserProfile struct { ... }

// Unexported: camelCase
func validateVideo(url string) error { ... }

// Interfaces: end with "er" when possible
type VideoUploader interface { ... }

// Constants: PascalCase or UPPER_CASE
const MaxVideoSizeMB = 50
```

**Code Examples:**

```go
// ✅ Good
func (s *VideoService) UploadVideo(ctx context.Context, userID string, video io.Reader) (string, error) {
    if err := s.validateUser(userID); err != nil {
        return "", fmt.Errorf("user validation failed: %w", err)
    }
    
    url, err := s.storage.Upload(ctx, video)
    if err != nil {
        return "", fmt.Errorf("upload failed: %w", err)
    }
    
    return url, nil
}

// ❌ Bad - Ignoring errors, no context
func (s *VideoService) UploadVideo(userID string, video io.Reader) string {
    s.storage.Upload(video) // No error handling
    return "url"
}
```

**Gin Router Example:**
```go
func SetupRouter() *gin.Engine {
    r := gin.Default()
    
    // Group routes
    api := r.Group("/api/v1")
    {
        api.POST("/videos", authMiddleware(), handleVideoUpload)
        api.GET("/matches", authMiddleware(), handleGetMatches)
    }
    
    return r
}
```

## 🧪 Testing Requirements

### Mobile App Tests

**Required for all new features:**
- Unit tests for utilities and hooks
- Component tests with React Testing Library
- Integration tests for critical flows

```bash
cd mobile
npm test
npm run test:coverage
```

**Coverage Requirements:**
- Minimum 80% coverage for new code
- 100% coverage for utility functions

**Example:**
```typescript
describe('useVideoUpload', () => {
  it('should upload video successfully', async () => {
    const { result } = renderHook(() => useVideoUpload());
    
    await act(async () => {
      await result.current.upload(mockFile);
    });
    
    expect(result.current.isUploading).toBe(false);
    expect(result.current.uploadUrl).toBeTruthy();
  });
});
```

### Backend Tests

**Required for all new features:**
- Unit tests for services and utilities
- Integration tests for API endpoints
- Database migration tests

```bash
cd backend
go test ./...
go test -cover ./...
```

**Coverage Requirements:**
- Minimum 75% coverage for new code
- 100% coverage for critical paths (auth, payments)

**Example:**
```go
func TestVideoService_Upload(t *testing.T) {
    service := NewVideoService(mockStorage, mockDB)
    
    video := bytes.NewReader([]byte("test video"))
    url, err := service.UploadVideo(context.Background(), "user123", video)
    
    assert.NoError(t, err)
    assert.NotEmpty(t, url)
}
```

## 📝 Documentation Requirements

**For all contributions:**
- Update relevant documentation
- Add inline comments for complex logic
- Update API.md for new endpoints
- Add examples for new features

**Documentation should include:**
- What the feature does
- Why it was implemented
- How to use it
- Any limitations or edge cases

## 🔍 Code Review Process

**What reviewers look for:**
1. **Functionality** - Does it work as intended?
2. **Tests** - Are there adequate tests?
3. **Documentation** - Is it well documented?
4. **Style** - Does it follow our guidelines?
5. **Performance** - Is it efficient?
6. **Security** - Are there security concerns?

**Response Time:**
- Maintainers aim to review PRs within 48 hours
- Complex PRs may take longer
- Feel free to ping after 3 days if no response

## 🐛 Reporting Bugs

**Before submitting a bug report:**
- Check if the bug has already been reported
- Ensure you're using the latest version
- Collect relevant information (logs, screenshots, etc.)

**Bug report should include:**
- Clear, descriptive title
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, versions, etc.)
- Screenshots or error logs

**Use the bug report template on GitHub Issues**

## 💡 Suggesting Features

**Feature requests should include:**
- Clear use case and problem statement
- Proposed solution
- Alternative solutions considered
- Mockups or examples (if applicable)

**Use the feature request template on GitHub Issues**

## 📜 License

By contributing to FindMe, you agree that your contributions will be licensed under the same [Commercial License](LICENSE) as the project.

## 🙏 Recognition

Contributors will be recognized in:
- Our [Contributors](https://github.com/yourusername/findme/graphs/contributors) page
- Release notes for significant contributions
- Special mentions in our documentation

## 📧 Questions?

- Open a [Discussion](https://github.com/yourusername/findme/discussions)
- Join our community chat
- Email: contributors@findme.ai

---

**Thank you for contributing to FindMe! Together, we're building something special. ❤️**
