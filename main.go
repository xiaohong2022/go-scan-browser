package main

import (
	"fmt"
	"math"
	"os"
)

func main() {
	// dir, err := os.ReadDir("C:\\")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// var toPrintList []string

	// for _, handle := range dir {
	// 	info, err := handle.Info()

	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}

	// 	toPrintList = append(toPrintList, fmt.Sprintf("%s %s", info.Name(), memorySizeToStr(int(info.Size()))))
	// }

	// for _, str := range toPrintList {
	// 	fmt.Println(str)
	// }

	/*_ := */
	paths := scanBrowserUserDataDirs()

	fmt.Print(`
0：退出
1：清理缓存文件（Cache、Code Cache、GPUCache）
选择操作：`)

	inp := 0
	fmt.Scanln(&inp)

	switch inp {
	case 0:
		break
	case 1:
		fmt.Print(`
本程序仅清除 文件夹 Cache、Code Cache 和 CPUCache 里的文件
清除后不可恢复
确认继续？(Y)：`)
		var inp2 string
		fmt.Scanln(&inp2)
		if inp2 == "y" || inp2 == "Y" {
			clearCacheDirs(paths)
		}
	}
}

func clearCacheDirs(paths []string) {
	fmt.Println("")
	fmt.Println("开始清除...")
	var free int64 = 0
	for _, path := range paths {
		pathsArray := [...]string{
			path + "Cache\\",
			path + "Code Cache\\",
			path + "CPUCache\\",
			path + "Default\\Cache\\",
			path + "Default\\Code Cache\\",
			path + "Default\\CPUCache\\",
		}

		for _, delPath := range pathsArray {
			// 文件夹存在就删，不存在跳过
			_, err := os.ReadDir(delPath)
			if err == nil {
				free += removeDirAndReturnFreeSize(delPath)
			}
		}
	}

	fmt.Printf("\n释放了 %s 空间", memorySizeToStr(free))

	var inp int
	fmt.Print("\n\n按回车键退出")
	fmt.Scanln(&inp)
}

func scanBrowserUserDataDirs() []string {
	var toPrintList []string
	var paths []string
	var _callback func(path string, list []os.DirEntry)
	callback := func(path string, list []os.DirEntry) {
		var isBrowserDirTicks int = 0
		var nameList = [...]string{
			"Local State",

			"Default",

			"Local Storage",
			"Service Worker",
			"Code Cache",
		}

		for _, handle := range list {
			for _, name := range nameList {
				if handle.Name() == name {
					isBrowserDirTicks++
					break
				}
			}
		}

		if isBrowserDirTicks >= 2 {
			toPrintList = append(toPrintList, fmt.Sprintf("%s %s", startPendingStr(memorySizeToStr(getDirSize(path)), 7), path))
			paths = append(paths, path)
		} else {
			for _, handle := range list {
				if handle.IsDir() {
					scanDir(path+handle.Name()+"\\", _callback)
				}
			}
		}

	}
	_callback = callback

	scanDir("C:\\Users\\", callback)

	fmt.Println("大小    路径")
	for _, str := range toPrintList {
		fmt.Println(str)
	}

	return paths
}

func scanDir(path string, callback func(path string, list []os.DirEntry)) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return
	}
	callback(path, dir)
}

func getDirSize(path string) int64 {
	var result int64 = 0
	var _callback func(path string, list []os.DirEntry)
	callback := func(path string, list []os.DirEntry) {
		for _, handle := range list {
			if handle.IsDir() {
				scanDir(path+handle.Name()+"\\", _callback)
			} else {
				info, err := handle.Info()

				if err != nil {
					continue
				}

				result += info.Size()
			}
		}

	}
	_callback = callback

	scanDir(path, callback)

	return result
}

func removeDirAndReturnFreeSize(path string) int64 {
	var result int64 = 0
	var _callback func(path string, list []os.DirEntry)
	callback := func(path string, list []os.DirEntry) {
		for _, handle := range list {
			if handle.IsDir() {
				scanDir(path+handle.Name()+"\\", _callback)
			} else {
				info, err := handle.Info()

				if err != nil {
					continue
				}

				size := info.Size()
				path_ := path + handle.Name()

				err = os.Remove(path_)

				if err != nil {
					println(err.Error())
					continue
				} else {
					println("removed " + path_)
				}

				result += size
			}
		}

	}
	_callback = callback

	scanDir(path, callback)

	return result
}

func memorySizeToStr(size int64) string {
	var (
		res  float32
		unit string
	)

	if size < 1024 {
		res, unit = float32(size), "B"
	} else if size < int64(math.Pow(1024, 2)) {
		res, unit = float32(size/1024), "KB"
	} else if size < int64(math.Pow(1024, 3)) {
		res, unit = float32(size/int64(math.Pow(1024, 2))), "MB"
	} else {
		res, unit = float32(size/int64(math.Pow(1024, 3))), "GB"
	}

	return fmt.Sprintf("%g%s", res, unit)
}

func startPendingStr(str string, length int) string {
	res := str
	if len(res) < length {
		for range length - len(res) {
			res = res + " "
		}
	}
	return res
}
