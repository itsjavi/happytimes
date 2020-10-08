package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

func PathExists(dir string) bool {
	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		return false
	}

	return true
}

func IsDir(dir string) bool {
	dirStat, err := os.Stat(dir)

	if os.IsNotExist(err) {
		return false
	}

	return dirStat.IsDir()
}

func FileGetMD5Checksum(path string) string {
	f, err := os.Open(path)
	defer f.Close()

	Catch(err)

	h := md5.New()
	if _, err := io.Copy(h, f); IsError(err) {
		Catch(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func FileAppend(path, str string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, FilePerms)
	Catch(err)

	defer f.Close()

	if _, err = f.WriteString(str); IsError(err) {
		Catch(err)
	}
}

func FileFixDates(path string, creationDate time.Time, modificationDate time.Time) error {
	if !IsUnix {
		return nil
	}
	err := exec.Command("touch", "-t", creationDate.Format("200601021504.05"), path).Run()

	if IsError(err) {
		return err
	}

	err = exec.Command("touch", "-mt", modificationDate.Format("200601021504.05"), path).Run()

	return err
}

func FileCopy(src, dest string, keepAttributes bool) error {
	if keepAttributes == true && IsUnix { // windows does not support cp nor preserving attributes
		err := exec.Command("cp", "-pRP", src, dest).Run()

		return err
	}
	s, err := os.Open(src)
	if IsError(err) {
		return err
	}

	defer s.Close()
	d, err := os.Create(dest)
	if IsError(err) {
		return err
	}
	if _, err := io.Copy(d, s); IsError(err) {
		d.Close()
		return err
	}
	return d.Close()
}

func FileMove(src, dest string) error {
	err := os.Rename(src, dest)

	if IsError(err) {
		return err
	}

	return nil
}

func MakeDirIfNotExists(dir string) {
	if !PathExists(dir) {
		Catch(os.MkdirAll(dir, DirPerms))
	}
}
