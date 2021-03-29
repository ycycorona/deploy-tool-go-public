package zip_file

import (
	"archive/zip"
	"deploy-tool-go/src/log"
	"deploy-tool-go/src/util"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

// 日志变量
var sugar *zap.SugaredLogger

func init() {
	sugar = log.Sugar
}

func Zip(dst, src string) (err error) {

	// src是否是文件夹
	var srcIsDir = true
	// 规整路径
	cleanSrc := filepath.Clean(src)
	cleanDst := filepath.Clean(dst)

	// 解压输出目标 必须是一个文件
	if util.IsDir(cleanDst) {
		panic(fmt.Sprintf(
			"cleanDst: \"%s\" \n action zip, can't regard dir as the dst, must specify a file path!!",
			cleanDst))
	}

	if util.IsFile(cleanSrc) {
		srcIsDir = false
	}

	// 创建准备写入的文件
	fw, err := os.Create(cleanDst)
	if err != nil {
		return err
	}
	defer func() {
		if err := fw.Close(); err != nil {
			sugar.Infof("fw.Close() fail %v", err)
		}
	}()

	// 通过 fw 来创建 zip.Write
	zw := zip.NewWriter(fw)
	defer func() {
		// 检测一下是否成功关闭
		if err := zw.Close(); err != nil {
			sugar.Infof("zip.NewWriter()fail。%v", err)
		}
	}()

	// 下面来将文件写入 zw ，因为有可能会有很多个目录及文件，所以递归处理
	return filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) (err error) {
		if errBack != nil {
			return errBack
		}
		// 规整路径
		// cleanSrc := filepath.Clean(src)

		path = filepath.Clean(path)

		// 通过文件信息，创建 zip 的文件信息
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}

		// 替换文件信息中的文件名
		// 如果是单文件，则正好取了文件名
		if srcIsDir {
			fh.Name = strings.Replace(path, cleanSrc, "", -1)
			// 清除路径的前导分隔符
			if len(fh.Name) != 0 {
				fh.Name = fh.Name[1:]
			}
		}
		// 替换为unix文件分隔符号
		fh.Name = filepath.ToSlash(fh.Name)
		// 这步开始没有加，会发现解压的时候说它不是个目录
		if fi.IsDir() {
			fh.Name += "/"
		} else {
			fh.Method = zip.Deflate
		}

		// 写入文件信息，并返回一个 Write 结构
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}

		// 检测，如果不是标准文件就只写入头信息，不写入文件数据到 w
		// 如目录，也没有数据需要写
		if !fh.Mode().IsRegular() {
			return nil
		}

		// 打开要压缩的文件
		fr, err := os.Open(path)
		defer func() {
			if err := fr.Close(); err != nil {
				sugar.Infof("fr.Close() fail %v", err)
			}
		}()
		if err != nil {
			return err
		}

		// 将打开的文件 Copy 到 w
		n, err := io.Copy(w, fr)
		if err != nil {
			return err
		}
		// 输出压缩的内容
		sugar.Infof("成功压缩文件： %s, 共写入了 %s KB的数据\n", path, util.ByteToMb(n, 2))

		return nil
	})
}
