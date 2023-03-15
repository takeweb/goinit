package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	var dir string
	var module string
	var file string

	// コマンドライン引数取得
	flag.StringVar(&module, "m", "test", "module to make")
	flag.StringVar(&dir, "d", GetCurrentDir(), "base dir to make")
	flag.StringVar(&file, "f", config.DefFilename, "filename to make")
	flag.Parse()

	// ターゲットディレクトリ作成
	targetDir := filepath.Join(dir, GetPathFromModule(module))
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		error(err)
		os.Exit(1)
	}
	info("module:", module)

	// テンプレートファイルコピー
	fromFile := filepath.Join(GetHomeDir(), config.TemplateDir, config.DefFilename+".go")
	toFile := filepath.Join(targetDir, file+".go")
	CopyFile(fromFile, toFile)

	// Makeファイルコピー
	fromFile2 := filepath.Join(GetHomeDir(), config.TemplateDir, "Makefile")
	toFile2 := filepath.Join(targetDir, "Makefile")
	CopyFile(fromFile2, toFile2)

	// ターゲットディレクトリへ移動
	if err := os.Chdir(targetDir); err != nil {
		error(err)
		os.Exit(1)
	}

	// 外部コマンド実行(go mod init <module>)
	if err := exec.Command("go", "mod", "init", module).Run(); err != nil {
		error(err)
		os.Exit(1)
	}

	// 一階層上のディレクトリへ移動
	if err := os.Chdir(dir); err != nil {
		error(err)
		os.Exit(1)
	}
	info("dir:", dir)

	// go.workの存在チェック
	f := "go.work"
	if _, err := os.Stat(f); err == nil {
		// 外部コマンド実行(go work use <module>)
		if err := exec.Command("go", "work", "use", module).Run(); err != nil {
			error(err)
			os.Exit(1)
		}
	} else {
		// 外部コマンド実行(go work init <module>)
		if err := exec.Command("go", "work", "init", module).Run(); err != nil {
			error(err)
			os.Exit(1)
		}
	}
}
