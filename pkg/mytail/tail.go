package mytail

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// Tail prints the last lines of the files to stdout
func Tail(filePaths []string, numLine uint64) error {
	defaultBufSize := int64(20 << 20) // 20MB
	out := os.Stdout
	err := tailFiles(filePaths, numLine, defaultBufSize, out)
	if err != nil {
		return err
	}

	return nil
}

func tailFiles(filePaths []string, numLine uint64, defaultBufSize int64, out io.Writer) error {
	resultChan := make(chan tailedFile, len(filePaths))
	errChan := make(chan error, len(filePaths))
	var wg sync.WaitGroup
	wg.Add(len(filePaths))
	for i, fp := range filePaths {
		go func(filePath string, printOrder int) {
			defer wg.Done()
			file, err := os.Open(filePath)
			if err != nil {
				errChan <- err
				return
			}
			defer func() {
				err := file.Close()
				if err != nil {
					errChan <- err
				}
			}()
			if tf, err := tail(file, numLine, defaultBufSize, printOrder); err != nil {
				errChan <- err
			} else {
				resultChan <- tf
			}
		}(fp, i)
	}
	wg.Wait()
	close(resultChan)
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	tfs := make([]tailedFile, len(filePaths))
	for tf := range resultChan {
		tfs[tf.order] = tf
	}

	return printTailedFiles(tfs, out)
}

func printTailedFiles(tfs []tailedFile, out io.Writer) error {
	printFileName := len(tfs) > 1
	for _, tf := range tfs {
		if printFileName {
			if tf.order > 0 {
				// From 2nd file,
				// print new line before print file name
				if _, err := fmt.Fprint(out, "\n"); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprint(out, "==> ", tf.name, " <==\n"); err != nil {
				return err
			}
		}

		for i := len(tf.lines) - 1; i >= 0; i-- {
			if _, err := fmt.Fprint(out, string(tf.lines[i])); err != nil {
				return err
			}
		}
	}

	return nil
}

func tail(file *os.File, numLine uint64, defaultBufSize int64, order int) (tailedFile, error) {
	tf := tailedFile{
		name:  file.Name(),
		order: order,
		lines: make([]line, 0),
	}

	fInfo, err := file.Stat()
	if err != nil {
		return tf, err
	}
	fSize := fInfo.Size()
	if fSize == 0 {
		return tf, nil
	}

	var bufSize = defaultBufSize
	var startLineIdx, endLineIdx int

	if fSize < defaultBufSize {
		bufSize = fSize
	}
	chunks := fSize / bufSize
	lastChunkSize := fSize % bufSize
	if lastChunkSize > 0 {
		chunks++
	}

	remainingNumBytes := int64(0)
	readBytes := make([]byte, 0)
	foundEOF := 0
	for c := int64(1); c <= chunks; c++ {
		if numLine == 0 {
			break
		}

		offset := fSize - bufSize*c
		if offset < 0 { // last chunk, and remaining bytes is less than bufSize
			offset = 0
			bufSize = lastChunkSize
		}

		readBytes = make([]byte, bufSize+remainingNumBytes)
		readNumBytes, err := file.ReadAt(readBytes, offset)
		if err != nil {
			tf.lines = []line{} // reset lines to free keeping memory
			return tf, err
		}

		endLineIdx = readNumBytes
		foundEOF = 0
		for i := readNumBytes - 1; i >= 0 && numLine > 0; i-- {
			if readBytes[i] == '\n' { // EOL
				foundEOF++
				if foundEOF == 2 { // found EOF for both start and end of line
					startLineIdx = i + 1
					tf.lines = append(tf.lines, readBytes[startLineIdx:endLineIdx])
					numLine--
					endLineIdx = startLineIdx
					foundEOF = 1
				}
			}
			startLineIdx = i
		}
		// keep tracking on remaining bytes in current chunk
		// then add to bufSize on reading next chunk
		remainingNumBytes = int64(endLineIdx)
	}

	if startLineIdx == 0 && numLine > 0 { // read the remaining bytes as first line of file
		tf.lines = append(tf.lines, readBytes[0:remainingNumBytes])
	}

	return tf, nil
}

type line []byte
type tailedFile struct {
	name  string
	order int
	lines []line
}
