package reviewer

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/cloudwego/eino/schema"
)

func (e *Engine) writeReviewToFile(path string, result any, language string) error {
    workDir, err := e.getWorkPath()
    if err != nil {
        return fmt.Errorf("failed to get work path: %v", err)
    }

    output := e.cfg.OutputFile
    if output == "" {
        output = "code-review.md"
    }
    outputFile := filepath.Join(workDir, output)
    file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to open output file: %v", err)
    }
    defer file.Close()

    // 格式化内容
    content := e.formatReviewResult(path, result, language)

    // 写入内容
    if _, err := file.WriteString(content); err != nil {
        return fmt.Errorf("failed to write to file: %v", err)
    }
    return nil
}

// 格式化审查结果
func (e *Engine) formatReviewResult(filePath string, result any, language string) string {
    timestamp := time.Now().Format("2006-01-02 15:04:05")

    var content string
    // 如果结果是 *schema.Message 类型，提取内容
    if msg, ok := result.(*schema.Message); ok {
        content = msg.Content
    } else {
        content = fmt.Sprintf("%v", result)
    }

    // 根据语言设置选择模板
    if e.cfg.Language == "en" {
        return fmt.Sprintf(`
## Code Review Report

**File Path**: %s  
**File Type**: %s  
**Review Time**: %s

### Review Result

%s

---

`, filePath, language, timestamp, content)
    } else {
        // 默认中文模板
        return fmt.Sprintf(`
## 文件审查报告

**文件路径**: %s  
**文件类型**: %s  
**审查时间**: %s

### 审查结果

%s

---

`, filePath, language, timestamp, content)
    }
}

// 获取文件语言类型的辅助函数
func getFileLanguage(filePath string) string {
    ext := strings.ToLower(filepath.Ext(filePath))
    languageMap := map[string]string{
        ".go":    "Go",
        ".js":    "JavaScript",
        ".jsx":   "JavaScript",
        ".ts":    "TypeScript",
        ".tsx":   "TypeScript",
        ".py":    "Python",
        ".java":  "Java",
        ".cpp":   "C++",
        ".cc":    "C++",
        ".cxx":   "C++",
        ".c":     "C",
        ".rs":    "Rust",
        ".php":   "PHP",
        ".rb":    "Ruby",
        ".swift": "Swift",
        ".kt":    "Kotlin",
        ".scala": "Scala",
        ".sh":    "Shell",
        ".sql":   "SQL",
        ".html":  "HTML",
        ".htm":   "HTML",
        ".css":   "CSS",
        ".yaml":  "YAML",
        ".yml":   "YAML",
        ".json":  "JSON",
        ".xml":   "XML",
        ".md":    "Markdown",
    }
    if lang, exists := languageMap[ext]; exists {
        return lang
    }
    return "Unknown"
}
