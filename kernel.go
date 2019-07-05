package vos

import (
	"fmt"
	"strings"

	"github.com/spf13/afero"
)

func (v *OS) defaultInit(*OS) ([]string, error) {
	v.Debug("DefaultInit")
	return nil, ErrInitScriptNotFound
}

// 遍历挂载树，通过 **绝对路径** 寻找文件对象, **使用前后记得加/解锁**
func (v *OS) findFileObject(path string) (afero.Fs, string) {
	v.Debug(fmt.Sprintf("findFileObject path=%s", path))
	//v.fsLock.RLock()
	//defer v.fsLock.RUnlock()
	for _, fs := range v.fs {
		v.Debug(fmt.Sprintf("findFileObject-for point=%s", fs.Point))
		if strings.HasPrefix(path, fs.Point) {
			p := strings.Replace(path, fs.Point, "", 1)
			if p == "" {
				p = "/"
			}
			if exist, _ := afero.Exists(fs.Mount, p); exist {
				v.Debug(fmt.Sprintf("findFileObject-for path=%s match point=%s", path, fs.Point))
				return fs.Mount, p
			}
		}
	}
	v.Debug(fmt.Sprintf("findFileObject-for path=%s nomatch", path))
	return nil, ""
}
