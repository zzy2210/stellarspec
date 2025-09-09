package reviewer

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"

    "github.com/fatih/color"
    "github.com/go-git/go-git/v5"
    "github.com/go-git/go-git/v5/plumbing/object"
    "github.com/sergi/go-diff/diffmatchpatch"
)

type gitDiff struct {
    // 文件
    FilePath string
    // 变更内容
    Content string
}

func (e *Engine) gitDiff() ([]gitDiff, error) {
    workPath, err := e.getWorkPath()
    if err != nil {
        return nil, err
    }

    repo, err := git.PlainOpen(workPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open repo: workPath=%s, err=%v", workPath, err)
    }

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
        return nil, fmt.Errorf("failed to get work tree: %v", err)
    }

    status, err := worktree.Status()
    if err != nil {
        return nil, fmt.Errorf("failed to get status: %v", err)
    }

    // 获取HEAD的tree
    headTree, err := commit.Tree()
    if err != nil {
        return nil, fmt.Errorf("failed to get HEAD tree: %v", err)
    }

    diffs := []gitDiff{}
    for file, fileStatus := range status {
        // 可选过滤：常见无关文件
        if file == "go.sum" || file == "go.mod" || strings.Contains(strings.ToLower(file), "readme") {
            continue
        }
        // 1. 未追踪文件：直接读取内容
        if fileStatus.Staging == git.Untracked || fileStatus.Worktree == git.Untracked {
            content, err := e.getFileContent(filepath.Join(workPath, file))
            if err != nil {
                color.Red("failed to get change path: path=%s, err=%v\n", file, err)
                continue
            }
            diffs = append(diffs, gitDiff{FilePath: file, Content: content})
            color.Yellow("Δ add: %s\n", filepath.Join(workPath, file))
        }
        // 2. 已修改文件：生成 diff
        if fileStatus.Staging == git.Modified || fileStatus.Worktree == git.Modified {
            diffContent, err := e.getModifiedFileDiff(repo, headTree, file, workPath)
            if err != nil {
                color.Red("failed to get diff for file: path=%s, err=%v\n", file, err)
                continue
            }
            diffs = append(diffs, gitDiff{FilePath: file, Content: diffContent})
            color.Yellow("Δ mod: %s\n", filepath.Join(workPath, file))
        }
        // 3. 已添加到暂存区的新文件
        if fileStatus.Staging == git.Added {
            content, err := e.getFileContent(filepath.Join(workPath, file))
            if err != nil {
                color.Red("failed to get file content: path=%s, err=%v\n", file, err)
                continue
            }
            diffs = append(diffs, gitDiff{FilePath: file, Content: content})
            color.Yellow("Δ staged: %s\n", filepath.Join(workPath, file))
        }
    }
    return diffs, nil
}

func (e *Engine) getWorkPath() (string, error) {
    var workPath string
    if e.cfg.ReviewPath == "" || e.cfg.ReviewPath == "." {
        wd, err := os.Getwd()
        if err != nil {
            return "", fmt.Errorf("get work path failed: %v", err)
        }
        workPath = wd
    } else {
        abs, err := filepath.Abs(e.cfg.ReviewPath)
        if err != nil {
            return "", fmt.Errorf("convert to abs path failed: %v", err)
        }
        workPath = abs
    }
    return workPath, nil
}

func (e *Engine) getFileContent(filePath string) (string, error) {
    content, err := os.ReadFile(filePath)
    if err != nil {
        return "", err
    }
    return string(content), nil
}

func (e *Engine) getModifiedFileDiff(repo *git.Repository, headTree *object.Tree, filePath, workPath string) (string, error) {
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

func (e *Engine) generateProfessionalDiff(filePath, oldContent, newContent string) string {
    dmp := diffmatchpatch.New()
    diffs := dmp.DiffMain(oldContent, newContent, false)
    return dmp.DiffPrettyText(diffs)
}
