package un_zip_file

import (
	"archive/zip"
	"deploy-tool-go/src/log"
	"deploy-tool-go/src/util"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

// 日志变量
var sugar *zap.SugaredLogger

func init() {
	sugar = log.Sugar
}

func UnZip(dst, src string) (err error) {
	// 规整路径
	cleanSrc := filepath.Clean(src)
	cleanDst := filepath.Clean(dst)
	// 输出目标必须是路径
	if res := !util.IsDir(cleanDst); res {
		panic(fmt.Sprintf("cleanDst: %v, !util.IsDir(cleanDst): %v \n 输出目标必须是路径", cleanDst, res))
	}
	// 解压对象必须是文件
	if res := !util.IsFile(cleanSrc); res {
		panic(fmt.Sprintf("cleanSrc: %v, !util.IsFile(cleanSrc): %v \n 解压对象必须是文件", cleanSrc, res))
	}

	// 打开压缩文件，这个 zip 包有个方便的 ReadCloser 类型
	// 这个里面有个方便的 OpenReader 函数，可以比 tar 的时候省去一个打开文件的步骤
	zr, err := zip.OpenReader(cleanSrc)
	defer func() {
		if err := zr.Close(); err != nil {
			sugar.Infof("zr.Close() fail. %v", err)
		}
	}()

	if err != nil {
		panic(fmt.Sprintf("zip.OpenReader() fail, %v", err))
	}

	// 如果解压后不是放在当前目录就按照保存目录去创建目录
	if cleanDst != "." {
		if err := os.MkdirAll(cleanDst, 0755); err != nil {
			return err
		}
	}

	// 遍历 zr ，将文件写入到磁盘
	for _, file := range zr.File {

		err := func() (err error) {
			cleanFileName := filepath.Clean(file.Name)
			path := filepath.Join(cleanDst, cleanFileName)
			fmt.Printf("file %s path %s \n", file.Name, path)
			// 如果是目录，就创建目录
			if file.FileInfo().IsDir() {
				if err := os.MkdirAll(path, file.Mode()); err != nil {
					return err
				}
				// 因为是目录，跳过当前循环，因为后面都是文件的处理
				return nil
			}

			// 获取到 Reader
			fr, err := file.Open()
			defer func() {
				if err := fr.Close(); err != nil {
					sugar.Infof("fr.Close() fail。%v", err)
				}
			}()
			if err != nil {
				return err
			}

			// 创建要写出的文件对应的 Write
			fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
			defer func() {
				if err := fw.Close(); err != nil {
					sugar.Infof("fw.Close() fail。%v", err)
				}
			}()
			if err != nil {
				return err
			}

			n, err := io.Copy(fw, fr)
			if err != nil {
				return err
			}

			// 将解压的结果输出
			sugar.Infof("成功解压 %s ，共写入了 %s KB的数据\n", path, util.ByteToMb(n, 2))

			return nil
		}()

		if err != nil {
			panic(fmt.Sprintf("unzip fail, %v", err))
		}
	}
	return nil
}
