package solutions

import (
	"fmt"
	"os"
	"sync"
)

func Tail(filePaths []string, numLine uint64) error {
	maxBufSize := int64(20 << 20) // 20MB
	if tailFiles, err := tailFiles(filePaths, numLine, maxBufSize); err != nil {
		return err
	} else {
		printTail(tailFiles)
		return nil
	}
}

func tailFiles(filePaths []string, numLine uint64, maxBufSize int64) (tailedFiles, error) {
	resultChan := make(chan tailedFile, len(filePaths))
	errChan := make(chan error, len(filePaths))
	var wg sync.WaitGroup
	wg.Add(len(filePaths))
	for i, fp := range filePaths {
		go func(filePath string, displayOrder int) {
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
			if tf, err := tail(file, numLine, maxBufSize, displayOrder); err != nil {
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
			return nil, err
		}
	}

	tfs := make([]tailedFile, len(filePaths))
	for tf := range resultChan {
		tfs[tf.order] = tf
	}

	return tfs, nil
}

func tail(file *os.File, numLine uint64, maxBufSize int64, order int) (tailedFile, error) {
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

	var bufSize = maxBufSize
	var startLineIdx, endLineIdx int

	if fSize < maxBufSize {
		bufSize = fSize
	}
	chunks := fSize / bufSize
	lastChunkSize := fSize % bufSize
	if lastChunkSize > 0 {
		chunks++
	}

	fmt.Println("chunks ---- ", chunks)
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
		if c == int64(1) { // set startLineIdx = endLineIdx in first chunk
			startLineIdx = endLineIdx
		}

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
		remainingNumBytes = int64(endLineIdx)
	}

	if startLineIdx == 0 && numLine > 0 { // read the remaining bytes as first line of file
		tf.lines = append(tf.lines, readBytes[0:remainingNumBytes])
	}

	return tf, nil
}

func printTail(tfs tailedFiles) {
	for _, tf := range tfs {
		for i := len(tf.lines) - 1; i >= 0; i-- {
			fmt.Print(string(tf.lines[i]))
		}
	}
}

type line []byte
type tailedFile struct {
	name  string
	order int
	lines []line
}
type tailedFiles []tailedFile
