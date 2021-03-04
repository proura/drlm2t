package model

import (
	"archive/tar"
	"compress/gzip"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
)

type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

func downloadFile(filepath string, url string) error {

	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Print("\n")

	// Close the file without defer so it can happen before Rename()
	out.Close()

	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}
	return nil
}

func generateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)

	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func generateMAC() string {

	buf := make([]byte, 6)
	var mac net.HardwareAddr

	_, err := rand.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	buf[0] |= 2

	return append(mac, 0x00, 0x16, 0x3e, buf[3], buf[4], buf[5]).String()
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func execCopy(src, dst string) {
	_, err := exec.Command("/bin/bash", "-c", "[ -f "+dst+" ]").Output()
	if err != nil {
		_, err = exec.Command("/bin/bash", "-c", "[ ! -f "+dst+" ] && cp "+src+" "+dst+"").Output()
		if err != nil {
			log.Fatal("- Error ", err, " creating file ", dst)
		} else {
			log.Println("+ File", dst, "created")
		}
	} else {
		log.Println("- File", dst, "already exists")
	}
}

func dirExists(dir string) bool {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func RemoveDirEsp(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	os.Remove(dir)
	return nil
}

func removeDir(dir string) {
	os.RemoveAll(dir)
}

func createDir(dir string) {
	os.MkdirAll(dir, os.ModePerm)
}

func ExtractTarGz(dst, fileDst string) {

	gzipStream, err := os.Open(fileDst)
	if err != nil {
		log.Println("Error opening gzip file")
	}

	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		log.Fatal("- ExtractTarGz: NewReader failed")
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("- ExtractTarGz: Next() failed: %s", err.Error())
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(dst+"/"+header.Name, 0755); err != nil {
				log.Println("* ExtractTarGz: Destination " + dst + "/" + header.Name + " already exists")
			}
		case tar.TypeReg:
			outFile, err := os.Create(dst + "/" + header.Name)
			if err != nil {
				log.Fatalf("- ExtractTarGz: Create() failed: %s", err.Error())
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				log.Fatalf("- ExtractTarGz: Copy() failed: %s", err.Error())
			}
			outFile.Close()

		default:
			log.Fatalf("ExtractTarGz: uknown type: %d in %s", header.Typeflag, header.Name)
		}
	}
}
