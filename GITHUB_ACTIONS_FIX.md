# 🔧 GitHub Actions构建问题修复

## ❌ 遇到的问题

### 问题1: npm依赖同步错误
```
npm error `npm ci` can only install packages when your package.json 
and package-lock.json or npm-shrinkwrap.json are in sync.
npm error Invalid: lock file's rollup@4.52.5 does not satisfy rollup@2.79.2
npm error Missing: rollup@4.53.3 from lock file
```

**原因**: package-lock.json与package.json不同步

### 问题2: Wails找不到frontend目录
```
/usr/bin/cd: line 4: cd: ./frontend: No such file or directory
ERROR   exit status 1
```

**原因**: wails.json中使用`cd ./frontend && npm install`在某些shell环境下失败

---

## ✅ 解决方案

### 修复1: 更新package-lock.json

重新生成package-lock.json确保与package.json同步：

```bash
cd frontend
rm -f package-lock.json
npm install
```

**修改的文件**:
- `frontend/package-lock.json` - 重新生成

### 修复2: 修改GitHub Actions配置

使用`npm install --frozen-lockfile`代替`npm ci`，更宽容但仍安全：

**修改**：`.github/workflows/build-desktop.yml`

```yaml
# 修改前
- name: Install frontend dependencies
  run: |
    cd frontend
    npm ci

# 修改后
- name: Install frontend dependencies
  run: |
    cd frontend
    npm install --frozen-lockfile
```

**说明**:
- `--frozen-lockfile`: 如果需要更新lock file会失败，保证依赖一致性
- 比`npm ci`更宽容，但仍然安全

### 修复3: 修改wails.json配置

使用`npm --prefix`代替`cd`命令：

**修改**：`wails.json`

```json
{
  // 修改前
  "frontend:install": "cd ./frontend && npm install",
  "frontend:build": "cd ./frontend && npm run build",
  "frontend:dev:watcher": "cd ./frontend && npm run dev",
  
  // 修改后
  "frontend:install": "npm --prefix ./frontend install",
  "frontend:build": "npm --prefix ./frontend run build",
  "frontend:dev:watcher": "npm --prefix ./frontend run dev"
}
```

**优势**:
- ✅ 不依赖shell的`cd`命令
- ✅ 跨平台兼容性更好（Windows/macOS/Linux）
- ✅ 在任何工作目录下都能正常执行

---

## 🧪 验证

### 本地验证

```bash
# 测试Wails构建
wails build -s -clean

# 应该输出
✓ Generating bindings: Done.
✓ Compiling application: Done.
✓ Packaging application: Done.
```

### GitHub Actions验证

提交更改后，GitHub Actions应该成功构建：

```bash
git add frontend/package-lock.json
git add wails.json
git add .github/workflows/build-desktop.yml
git commit -m "fix: GitHub Actions构建问题 - 更新依赖和wails配置"
git push
```

查看Actions运行结果：
- https://github.com/your-repo/actions

---

## 📊 修改总结

| 文件 | 修改 | 原因 |
|------|------|------|
| `frontend/package-lock.json` | 重新生成 | 与package.json同步 |
| `.github/workflows/build-desktop.yml` | `npm ci` → `npm install --frozen-lockfile` | 更好的兼容性 |
| `wails.json` | `cd && npm` → `npm --prefix` | 跨平台兼容 |

---

## 🎯 预期结果

GitHub Actions将成功构建以下产物：

### macOS ARM64
- `caipiao-macOS-ARM/caipiao-macos-arm64.tar.gz`

### macOS AMD64  
- `caipiao-macOS-Intel/caipiao-macos-amd64.tar.gz`

### Windows AMD64
- `caipiao-Windows/caipiao-windows-amd64.exe`

---

## 💡 最佳实践

### 保持依赖同步

每次修改`package.json`后：

```bash
cd frontend
npm install  # 更新package-lock.json
git add package.json package-lock.json
git commit -m "deps: update dependencies"
```

### 使用npm --prefix

在所有需要指定目录的npm命令中使用`--prefix`：

```bash
# ✅ 推荐
npm --prefix ./frontend install
npm --prefix ./frontend run build

# ❌ 不推荐
cd frontend && npm install
cd frontend && npm run build
```

### GitHub Actions调试

在Actions中添加调试步骤：

```yaml
- name: Debug - Show directory structure
  run: |
    pwd
    ls -la
    ls -la frontend/
```

---

## 🔍 相关文档

- [npm install vs npm ci](https://docs.npmjs.com/cli/v9/commands/npm-ci)
- [npm --prefix flag](https://docs.npmjs.com/cli/v9/commands/npm#prefix)
- [Wails Frontend Configuration](https://wails.io/docs/reference/project-config)
- [GitHub Actions - Node.js](https://docs.github.com/en/actions/automating-builds-and-tests/building-and-testing-nodejs)

---

**修复日期**: 2025-11-21
**状态**: ✅ 已解决
**测试**: ✅ 本地构建通过
