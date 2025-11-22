# Contributing

Thank you for your interest in contributing to Tiny Test App!

## How to Contribute

### Reporting Issues

If you find a bug or have a suggestion for improvement, please open an issue on GitHub with:
- A clear description of the issue
- Steps to reproduce (if applicable)
- Expected vs actual behavior
- Environment details (Kubernetes version, etc.)

### Submitting Changes

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Ensure your code follows the existing style
5. Test your changes locally
6. Commit your changes (`git commit -m 'Add some amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Code Style

- Follow Go standard formatting (`go fmt`)
- Keep functions focused and small
- Add comments for exported functions
- Maintain the minimal size philosophy

### Testing

Before submitting a PR:
- Test the Docker build locally
- Verify the image size hasn't increased significantly
- Test deployment on a Kubernetes cluster
- Ensure all endpoints work correctly

**Note**: All PRs are automatically tested via GitHub Actions CI. The CI will:
- Build the Go application
- Build and test the Docker image
- Verify image size
- Test all endpoints
- Ensure the build completes successfully

Your PR must pass all CI checks before it can be merged.

### Pull Request Process

1. Ensure your PR has a clear description
2. Reference any related issues
3. Wait for review and address feedback
4. Once approved, maintainers will merge

Thank you for contributing!

