package reviewer

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	config "stellarspec/internal/model/conf"
	"strings"
	"sync"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	nodeOfModel  = "model"
	nodeOfPrompt = "prompt"
)

type ReviewEngine struct {
	ctx context.Context

	reviewPath string
	commitID   string
	chatModel  *openai.ChatModel

	// 文件锁
	mutex sync.Mutex
}

func NewReviewEngine(ctx context.Context, path string, commitID string) *ReviewEngine {
	return &ReviewEngine{
		ctx:        ctx,
		reviewPath: path,
		commitID:   commitID,
	}
}

func (e *ReviewEngine) CreateModel(conf *config.BaseConfig) {
	modelConf := &openai.ChatModelConfig{
		APIKey:  conf.Key,
		BaseURL: conf.APIServer,
		Model:   conf.Model,
	}
	chatModel, err := openai.NewChatModel(e.ctx, modelConf)
	if err != nil {
		return
	}
	e.chatModel = chatModel
}

func (e *ReviewEngine) Run() {
	diffs, err := e.gitDiff()
	if err != nil {
		fmt.Printf("get git diff failed: err= %v \n", err)
		return
	}

	var wg sync.WaitGroup
	maxWorkers := 10
	semaphore := make(chan struct{}, maxWorkers)
	for _, diff := range diffs {
		wg.Add(1)
		go func(d gitDiff) {
			defer wg.Done()
			semaphore <- struct{}{}        // 获取信号量
			defer func() { <-semaphore }() // 释放信号量
			e.reviewerSignleFile(d)
		}(diff)

	}
	wg.Wait()
}

type gitDiff struct {
	// 文件
	FilePath string
	// 变更内容
	Content string
}

func (e *ReviewEngine) gitDiff() ([]gitDiff, error) {
	workPath, err := e.getWorkPath()
	if err != nil {
		return nil, err
	}

	repo, err := git.PlainOpen(workPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repo: workPath = %s, err= %v,", workPath, err)
	}

	// 如果指定了 commitID，则审查指定 commit 的变更
	if e.commitID != "" {
		return e.getCommitDiff(repo)
	}

	// 否则审查工作区变更（原有逻辑）
	return e.getWorkingTreeDiff(repo, workPath)
}

func (e *ReviewEngine) getCommitDiff(repo *git.Repository) ([]gitDiff, error) {
	// 解析指定的 commit
	commitHash, err := repo.ResolveRevision(plumbing.Revision(e.commitID))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve commit %s: %v", e.commitID, err)
	}

	commit, err := repo.CommitObject(*commitHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit object %s: %v", e.commitID, err)
	}

	// 获取 commit 的 tree
	commitTree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get commit tree: %v", err)
	}

	var parentTree *object.Tree
	// 获取父 commit 的 tree（如果存在）
	if commit.NumParents() > 0 {
		parentCommit, err := commit.Parent(0)
		if err != nil {
			return nil, fmt.Errorf("failed to get parent commit: %v", err)
		}
		parentTree, err = parentCommit.Tree()
		if err != nil {
			return nil, fmt.Errorf("failed to get parent tree: %v", err)
		}
	}

	// 比较两个 tree，获取变更
	changes, err := object.DiffTree(parentTree, commitTree)
	if err != nil {
		return nil, fmt.Errorf("failed to diff trees: %v", err)
	}

	diffs := []gitDiff{}
	for _, change := range changes {
		// 跳过一些不需要审查的文件
		fileName := change.To.Name
		if fileName == "" {
			fileName = change.From.Name
		}
		if fileName == "go.sum" || fileName == "go.mod" || strings.Contains(fileName, "README") {
			continue
		}

		// 生成 patch
		patch, err := change.Patch()
		if err != nil {
			fmt.Printf("failed to get patch for file: path= %s, err= %v \n", fileName, err)
			continue
		}

		diffs = append(diffs, gitDiff{
			FilePath: fileName,
			Content:  patch.String(),
		})
	}

	return diffs, nil
}

func (e *ReviewEngine) getWorkingTreeDiff(repo *git.Repository, workPath string) ([]gitDiff, error) {
	// 获取HEAD commit
	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %v", err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get commit: %v", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get work tree, err= %v", err)
	}

	status, err := worktree.Status()
	if err != nil {
		return nil, fmt.Errorf("failed to get status, err= %v", err)
	}

	// 获取HEAD的tree
	headTree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD tree: %v", err)
	}

	diffs := []gitDiff{}

	for file, fileStatus := range status {
		if file == "go.sum" || file == "go.mod" || strings.Contains(file, "README") {
			continue
		}
		// 1. 新增文件（未跟踪）
		if fileStatus.Staging == git.Untracked || fileStatus.Worktree == git.Untracked {
			content, err := e.getFileContent(filepath.Join(workPath, file))
			if err != nil {
				fmt.Printf("failed to get change path: path= %s, err= %v \n", file, err)
				continue
			}
			diffs = append(diffs, gitDiff{
				FilePath: file,
				Content:  content,
			})
		}
		// 2. 已修改的文件（需要获取diff）
		if fileStatus.Staging == git.Modified || fileStatus.Worktree == git.Modified {
			diffContent, err := e.getModifiedFileDiff(repo, headTree, file, workPath)
			if err != nil {
				fmt.Printf("failed to get diff for file: path= %s, err= %v \n", file, err)
				continue
			}
			diffs = append(diffs, gitDiff{
				FilePath: file,
				Content:  diffContent,
			})
		}

		// 3. 已添加到暂存区的新文件
		if fileStatus.Staging == git.Added {
			content, err := e.getFileContent(filepath.Join(workPath, file))
			if err != nil {
				fmt.Printf("failed to get file content: path= %s, err= %v \n", file, err)
				continue
			}
			diffs = append(diffs, gitDiff{
				FilePath: file,
				Content:  content,
			})
		}
	}

	return diffs, nil
}

func (e *ReviewEngine) getWorkPath() (string, error) {
	var workPath string
	if e.reviewPath == "" || e.reviewPath == "." {
		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("get work path failed: err= %v", err)
		}
		workPath = wd

	} else {
		abs, err := filepath.Abs(e.reviewPath)
		if err != nil {
			return "", fmt.Errorf("convert to abs path failed: err= %v", err)
		}
		workPath = abs

	}
	return workPath, nil
}

func (e *ReviewEngine) getFileContent(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (e *ReviewEngine) getModifiedFileDiff(repo *git.Repository, headTree *object.Tree, filePath, workPath string) (string, error) {
	// 获取HEAD中的文件内容
	var oldContent string
	if entry, err := headTree.FindEntry(filePath); err == nil {
		blob, err := repo.BlobObject(entry.Hash)
		if err != nil {
			return "", fmt.Errorf("failed to get blob: %v", err)
		}

		reader, err := blob.Reader()
		if err != nil {
			return "", fmt.Errorf("failed to get blob reader: %v", err)
		}
		defer reader.Close()

		content, err := io.ReadAll(reader)
		if err != nil {
			return "", fmt.Errorf("failed to read blob content: %v", err)
		}
		oldContent = string(content)
	}

	// 获取当前工作区的文件内容
	newContent, err := e.getFileContent(filepath.Join(workPath, filePath))
	if err != nil {
		return "", fmt.Errorf("failed to get current file content: %v", err)
	}
	return e.generateProfessionalDiff(filePath, oldContent, newContent), nil

}

func (e *ReviewEngine) generateProfessionalDiff(filePath, oldContent, newContent string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(oldContent, newContent, false)
	return dmp.DiffPrettyText(diffs)
}

func (e *ReviewEngine) reviewerSignleFile(d gitDiff) {
	g := compose.NewGraph[map[string]any, *schema.Message]()
	ext := filepath.Ext(d.FilePath)

	systemTpl := fmt.Sprintf("你是一位  %s 研发专家，现在你将对用户给出的代码变更内容给出对应的code reviewer 结论。我需要你在结论中输出原有代码相关问题，你的评审建议，与修改方案.请将整体输出控制在200字内", ext)

	chatTpl := prompt.FromMessages(schema.FString,
		schema.SystemMessage(systemTpl),
		schema.MessagesPlaceholder("message_histories", true),
		schema.UserMessage("{user_query}"),
	)
	_ = g.AddChatTemplateNode(nodeOfPrompt, chatTpl)
	_ = g.AddChatModelNode(nodeOfModel, e.chatModel)
	_ = g.AddEdge(compose.START, nodeOfPrompt)
	_ = g.AddEdge(nodeOfPrompt, nodeOfModel)
	_ = g.AddEdge(nodeOfModel, compose.END)
	r, err := g.Compile(e.ctx, compose.WithMaxRunSteps(10))
	if err != nil {
		panic(err)
	}

	ret, err := r.Invoke(e.ctx, map[string]any{
		"message_histories": []*schema.Message{},
		"user_query":        d.Content,
	})
	if err != nil {
		fmt.Printf("invoke failed: err= %v \n", err)
		return
	}

	e.mutex.Lock()
	defer e.mutex.Unlock()
	if err := e.writeReviewToFile(d.FilePath, ret, ext); err != nil {
		fmt.Printf("write review faile: err= %v", err)
		return
	}

}

func (e *ReviewEngine) writeReviewToFile(path string, result any, language string) error {
	workDir, err := e.getWorkPath()
	if err != nil {
		return fmt.Errorf("failed to get work path: err= %v", err)
	}

	outputFile := filepath.Join(workDir, "code-review.md")
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
func (e *ReviewEngine) formatReviewResult(filePath string, result any, language string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	var content string

	// 如果结果是 *schema.Message 类型，提取内容
	if msg, ok := result.(*schema.Message); ok {
		content = msg.Content
	} else {
		content = fmt.Sprintf("%v", result)
	}

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
