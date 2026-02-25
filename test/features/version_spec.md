# Feature Spec: muxtree version

This document captures the exact behavior of `muxtree version` command for feature parity in the Go reimplementation.

## Test Case: Version output

**Input:**
```bash
muxtree version
```

**Expected Output:**
```
muxtree v1.0.0
```

**Exit Code:** 0

**Behavior:**
- Prints a single line with the version number
- Output is identical for `muxtree version`, `muxtree -v`, and `muxtree --version`

---

## Test Case: Alias - -v

**Input:**
```bash
muxtree -v
```

**Expected Output:**
```
(Same output as "Version output")
```

**Exit Code:** 0

---

## Test Case: Alias - --version

**Input:**
```bash
muxtree --version
```

**Expected Output:**
```
(Same output as "Version output")
```

**Exit Code:** 0
