/*  CREDITS:

all functions within this file are taken from the following source:

Summerfield, Mark. Programming in Go : creating applications for the 21st century. Addison-Wesley 2012.

*/

package compression

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"config"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var DIR = config.DIR

// helper function for compression
func sanitizedName(filename string) string {
	if len(filename) > 1 && filename[1] == ':' &&
		runtime.GOOS == "windows" {
		filename = filename[2:]
	}
	filename = filepath.ToSlash(filename)
	filename = strings.TrimLeft(filename, "/.")
	return strings.Replace(filename, "../", "", -1)
}

func writeFileToTar(writer *tar.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		log.Println(err)
		return err
	}
	header := &tar.Header{
		Name:    sanitizedName(filename),
		Mode:    int64(stat.Mode()),
		Uid:     os.Getuid(),
		Gid:     os.Getgid(),
		Size:    stat.Size(),
		ModTime: stat.ModTime(),
	}
	if err = writer.WriteHeader(header); err != nil {
		log.Println(err)
		return err
	}
	_, err = io.Copy(writer, file)
	return err
}

func CreateTar(filename string, files []string) error {
	file, err := os.Create(filename)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()
	var fileWriter io.WriteCloser = file
	if strings.HasSuffix(filename, ".gz") {
		fileWriter = gzip.NewWriter(file)
		defer fileWriter.Close()
	}
	writer := tar.NewWriter(fileWriter)
	defer writer.Close()
	for _, name := range files {
		if err := writeFileToTar(writer, name); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func writeFileToZip(zipper *zip.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		log.Println(err)
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		log.Println(err)
		return err
	}
	header.Name = sanitizedName(filename)
	writer, err := zipper.CreateHeader(header)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = io.Copy(writer, file)
	return err
}

func CreateZip(filename string, files []string) error {
	file, err := os.Create(filename)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()
	zipper := zip.NewWriter(file)
	defer zipper.Close()
	for _, name := range files {
		fmt.Println("writing", name, "into", filename)
		if err := writeFileToZip(zipper, name); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

// unpack tar files

func UnpackTar(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return err
	} else {
		log.Println("func UnpackTar: opening", filename)
	}
	defer file.Close()
	var fileReader io.ReadCloser = file
	if strings.HasSuffix(filename, ".gz") {
		log.Println("func UnpackTar:", filename, "hasSuffix .gz!")
		if fileReader, err = gzip.NewReader(file); err != nil {
			log.Println(err)
			return err
		}
		defer fileReader.Close()
	}
	reader := tar.NewReader(fileReader)
	return unpackTarFiles(reader)
}

/*
func UnpackTar(sourcefile string) {
	file, err := os.Open(sourcefile)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer file.Close()

	var fileReader io.ReadCloser = file

	// just in case we are reading a tar.gz file, add a filter to handle gzipped file
	if strings.HasSuffix(sourcefile, ".gz") {
		if fileReader, err = gzip.NewReader(file); err != nil {

			fmt.Println(err)
			os.Exit(1)
		}
		defer fileReader.Close()
	}

	tarBallReader := tar.NewReader(fileReader)

	// Extracting tarred files

	for {
		header, err := tarBallReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			os.Exit(1)
		}

		// get the individual filename and extract to the current directory
		filename := header.Name

		switch header.Typeflag {
		case tar.TypeDir:
			// handle directory
			fmt.Println("Creating directory :", filename)
			err = os.MkdirAll(filename, os.FileMode(header.Mode)) // or use 0755 if you prefer

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		case tar.TypeReg:
			// handle normal file
			fmt.Println("Untarring :", filename)
			writer, err := os.Create(filename)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			io.Copy(writer, tarBallReader)

			err = os.Chmod(filename, os.FileMode(header.Mode))

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			writer.Close()
		default:
			fmt.Printf("Unable to untar type : %c in file %s", header.Typeflag, filename)
		}
	}
}
*/
//helper for unpackTar
func unpackTarFile(filename, tarFilename string, reader *tar.Reader) error {
	writer, err := os.Create(DIR + filename)
	log.Println("func unpackTarFile: trying to create", filename)
	if err != nil {
		log.Println(err)
		return err
	}
	defer writer.Close()
	if _, err = io.Copy(writer, reader); err != nil {
		log.Println(err)
		return err
	}
	if filename == tarFilename {
		log.Println("func unpackTarFile:", filename, "==", tarFilename, "!!")
		fmt.Println(filename)
	} else {
		fmt.Printf("%s [%s]\n", filename, tarFilename)
	}
	return nil
}

//helper for unpackTar
func unpackTarFiles(reader *tar.Reader) error {
	for {
		header, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				return nil // OK
			}
			return err
		}
		filename := sanitizedName(header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err = os.MkdirAll(filename, 0755); err != nil {
				log.Println("tried to create filename")
				log.Println(err)
				return err
			}
		case tar.TypeReg:
			if err = unpackTarFile(filename, header.Name, reader); err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return nil
}

// unzips a file from src to dst
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			f, err := os.OpenFile(
				path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()
			// kopiert aus reader in writer
			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
