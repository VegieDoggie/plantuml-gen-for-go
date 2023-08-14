package filepathx

import (
	"github.com/VegetableDoggies/plantuml-gen-for-go/utils/arraysx"
	"os"
	"path/filepath"
	"strings"
)

// WalkWithExcludes 检索路径及其下某后缀名的所有文件
// @param nodePath - 路径
// @param suffix - 文件后缀(文件扩展名是后缀的子集)
// @param excludes - 排除文件(夹)列表
// @param isSkipHidden - 是否跳过隐藏文件(夹)，即"."开头
func WalkWithExcludes(nodePath string, suffix string, excludes []string, isSkipHidden bool) (targetFiles []string, err error) {
	nodePath, err = filepath.Abs(nodePath)
	if err != nil {
		return nil, err
	}
	for i := range excludes {
		if excludes[i], err = filepath.Abs(excludes[i]); err != nil {
			return nil, err
		}
	}
	var stat os.FileInfo
	if stat, err = os.Stat(nodePath); err == nil {
		if stat.IsDir() {
			_ = filepath.WalkDir(nodePath, func(path string, d os.DirEntry, err error) error {
				if isSkipHidden && strings.HasPrefix(filepath.Base(path), ".") {
					return filepath.SkipDir
				}
				if l := len(excludes); l > 0 {
					for i := 0; i < l; i++ {
						if strings.HasPrefix(path, excludes[i]) {
							excludes = arraysx.RemoveByIndex(excludes, i)
							if d.IsDir() {
								return filepath.SkipDir
							}
							return nil
						}
					}
				}
				if strings.HasSuffix(d.Name(), suffix) {
					targetFiles = append(targetFiles, path)
				}
				return err
			})
		} else {
			if strings.HasSuffix(nodePath, suffix) {
				targetFiles = append(targetFiles, nodePath)
			}
		}
		return targetFiles, nil
	}
	return nil, err
}

func IsDir(path string) bool {
	if stat, err := os.Stat(path); err == nil {
		return stat.IsDir()
	}
	return false
}

// Sub1Files 返回rootDir目录子级的文件列表
func Sub1Files(rootDir string, files []string) (subFiles []string) {
	if strings.HasSuffix(rootDir, string(os.PathSeparator)) {
		rootDir = rootDir[:len(rootDir)-1]
	}
	for i := 0; i < len(files); {
		if rootDir == filepath.Dir(files[i]) {
			subFiles = append(subFiles, files[i])
		}
	}
	return subFiles
}

// Sub1Dirs 返回某目录子级的文件夹列表
func Sub1Dirs(dir string, isSkipHidden bool) (dirs []string) {
	des, _ := os.ReadDir(dir)
	if isSkipHidden {
		for i := range des {
			if des[i].IsDir() && des[i].Name()[0] != '.' {
				dirs = append(dirs, filepath.Join(dir, des[i].Name()))
			}
		}
	} else {
		for i := range des {
			if des[i].IsDir() {
				dirs = append(dirs, filepath.Join(dir, des[i].Name()))
			}
		}
	}
	return
}
