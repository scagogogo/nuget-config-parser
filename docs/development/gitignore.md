# .gitignore Configuration

This document explains the `.gitignore` configuration for the NuGet Config Parser project.

## Overview

The `.gitignore` file is carefully configured to ignore files that should not be tracked in version control while ensuring that important files are preserved.

## Categories of Ignored Files

### 🔧 Go-specific Files

- **Binaries**: `*.exe`, `*.dll`, `*.so`, `*.dylib`
- **Test files**: `*.test`
- **Coverage reports**: `coverage.out`, `coverage.html`, `coverage.xml`
- **Profiling data**: `*.prof`, `*.pprof`, `*.trace`
- **Benchmark results**: `*.bench`
- **Common binary names**: `main`, `nuget-config-parser`

### 📚 Documentation Files

- **Node.js dependencies**: `docs/node_modules/`
- **VitePress build output**: `docs/.vitepress/dist/`
- **VitePress cache**: `docs/.vitepress/cache/`
- **Temporary files**: `docs/.temp/`, `docs/.cache/`

### 💻 IDE and Editor Files

- **VS Code**: `.vscode/` (all settings)
- **IntelliJ IDEA**: `.idea/`
- **Vim**: `*.swp`, `*.swo`, `*~`
- **Emacs**: Various Emacs-specific files

### 🖥️ Operating System Files

- **macOS**: `.DS_Store`, `.AppleDouble`, `._*`
- **Windows**: `Thumbs.db`, `desktop.ini`, `*.lnk`
- **Linux**: `.directory`, `.Trash-*`, `.nfs*`

### 🎯 Project-specific Files

- **Example configs**: `examples/*/NuGet.Config`, `examples/*/*.Config`
- **Test artifacts**: `test-output/`, `test-results/`, `*.test.xml`
- **Temporary configs**: `temp-*.config`, `test-*.config`, `*.tmp.config`

## Important Files That ARE Tracked

These files are intentionally tracked and should not be ignored:

- **`go.mod`** and **`go.sum`**: Essential for reproducible builds
- **`docs/package.json`** and **`docs/package-lock.json`**: Documentation dependencies
- **Source code**: All `.go` files
- **Documentation**: All `.md` files
- **Configuration**: Project configuration files

## Validation

Use the provided script to validate the `.gitignore` configuration:

```bash
./scripts/check-gitignore.sh
```

This script will:
- Test that important file types are properly ignored
- Verify that essential files are still tracked
- Check for common configuration mistakes
- Provide helpful tips and warnings

## Common Commands

### Check if a file is ignored
```bash
git check-ignore -v <filename>
```

### See all ignored files
```bash
git status --ignored
```

### Force add an ignored file
```bash
git add -f <filename>
```

### Remove a tracked file and ignore it
```bash
git rm --cached <filename>
# Then add the pattern to .gitignore
```

## Troubleshooting

### File is tracked but should be ignored

If a file is already tracked by Git, adding it to `.gitignore` won't ignore it. You need to:

1. Remove it from tracking: `git rm --cached <filename>`
2. Add the pattern to `.gitignore`
3. Commit the changes

### File is ignored but should be tracked

If a file is being ignored but you want to track it:

1. Check what's ignoring it: `git check-ignore -v <filename>`
2. Either remove the pattern from `.gitignore` or use `git add -f <filename>`

### Global gitignore conflicts

Check your global gitignore configuration:

```bash
git config --get core.excludesfile
cat $(git config --get core.excludesfile)
```

Global rules might conflict with local ones.

## Best Practices

1. **Test changes**: Run `./scripts/check-gitignore.sh` after modifying `.gitignore`
2. **Be specific**: Use specific patterns rather than broad wildcards
3. **Document exceptions**: If you need to track a normally-ignored file, document why
4. **Regular review**: Periodically review and clean up the `.gitignore` file

## Pattern Examples

### Ignore all files with extension
```gitignore
*.log
*.tmp
```

### Ignore directory
```gitignore
build/
temp/
```

### Ignore files in any subdirectory
```gitignore
**/node_modules/
**/dist/
```

### Ignore specific file
```gitignore
config.local.json
secret.key
```

### Exception (don't ignore)
```gitignore
# Ignore all .env files
.env*
# But track the example
!.env.example
```

## Maintenance

The `.gitignore` file should be reviewed and updated when:

- Adding new build tools or dependencies
- Changing development environments
- Adding new file types to the project
- Team members report issues with ignored/tracked files

For questions or suggestions about the `.gitignore` configuration, please open an issue or discussion in the project repository.
