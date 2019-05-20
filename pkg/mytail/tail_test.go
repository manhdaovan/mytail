package mytail

import (
	"testing"
)

type outTest struct {
	content []string
}

func (out *outTest) Write(p []byte) (int, error) {
	out.content = append(out.content, string(p))
	return len(p), nil
}

func sameContent(content1 []string, content2 []string) bool {
	if len(content1) != len(content2) {
		return false
	}

	allLinesSame := true
	for idx, line := range content1 {
		if line != content2[idx] {
			allLinesSame = false
		}
	}

	return allLinesSame
}

func Test_tailFiles(t *testing.T) {
	fileSize135Bytes5Lines27BytesPerLine := "./test_data/FileSize135Bytes_5Lines_27BytesPerLine.txt"
	fileSize198Bytes3Lines32BytesAnd3Lines34BytesPerLine := "./test_data/FileSize198Bytes_3Lines32Bytes_And3Lines34BytesPerLine.txt"
	emptyFile := "./test_data/empty.txt"

	type args struct {
		filePaths      []string
		numLine        uint64
		defaultBufSize int64
	}
	tests := []struct {
		name    string
		args    args
		wantOut []string
		wantErr bool
	}{
		// 1. Tail part of single file
		{
			name: "tail part of 1 file && bufSize > filesize",
			args: args{
				filePaths:      []string{fileSize198Bytes3Lines32BytesAnd3Lines34BytesPerLine},
				numLine:        4,
				defaultBufSize: 300,
			},
			wantOut: []string{
				"2,20Lines_32And34BytesPerLine,2\n",
				"10,20Lines_32And34BytesPerLine,10\n",
				"11,20Lines_32And34BytesPerLine,11\n",
				"12,20Lines_32And34BytesPerLine,12\n",
			},
			wantErr: false,
		},
		{
			name: "tail part of 1 file && bufSize == filesize",
			args: args{
				filePaths:      []string{fileSize198Bytes3Lines32BytesAnd3Lines34BytesPerLine},
				numLine:        4,
				defaultBufSize: 198,
			},
			wantOut: []string{
				"2,20Lines_32And34BytesPerLine,2\n",
				"10,20Lines_32And34BytesPerLine,10\n",
				"11,20Lines_32And34BytesPerLine,11\n",
				"12,20Lines_32And34BytesPerLine,12\n",
			},
			wantErr: false,
		},
		{
			name: "tail part of 1 file && bufSize > 1line bytes",
			args: args{
				filePaths:      []string{fileSize198Bytes3Lines32BytesAnd3Lines34BytesPerLine},
				numLine:        4,
				defaultBufSize: 300,
			},
			wantOut: []string{
				"2,20Lines_32And34BytesPerLine,2\n",
				"10,20Lines_32And34BytesPerLine,10\n",
				"11,20Lines_32And34BytesPerLine,11\n",
				"12,20Lines_32And34BytesPerLine,12\n",
			},
			wantErr: false,
		},
		{
			name: "tail part of 1 file && bufSize = 1 short line",
			args: args{
				filePaths:      []string{fileSize198Bytes3Lines32BytesAnd3Lines34BytesPerLine},
				numLine:        3,
				defaultBufSize: 32,
			},
			wantOut: []string{
				"10,20Lines_32And34BytesPerLine,10\n",
				"11,20Lines_32And34BytesPerLine,11\n",
				"12,20Lines_32And34BytesPerLine,12\n",
			},
			wantErr: false,
		},
		{
			name: "tail part of 1 file && bufSize = 1 long line",
			args: args{
				filePaths:      []string{fileSize198Bytes3Lines32BytesAnd3Lines34BytesPerLine},
				numLine:        3,
				defaultBufSize: 34,
			},
			wantOut: []string{
				"10,20Lines_32And34BytesPerLine,10\n",
				"11,20Lines_32And34BytesPerLine,11\n",
				"12,20Lines_32And34BytesPerLine,12\n",
			},
			wantErr: false,
		},
		{
			name: "tail part of 1 file && bufSize in middle short and long line",
			args: args{
				filePaths:      []string{fileSize198Bytes3Lines32BytesAnd3Lines34BytesPerLine},
				numLine:        3,
				defaultBufSize: 33,
			},
			wantOut: []string{
				"10,20Lines_32And34BytesPerLine,10\n",
				"11,20Lines_32And34BytesPerLine,11\n",
				"12,20Lines_32And34BytesPerLine,12\n",
			},
			wantErr: false,
		},
		{
			name: "tail part of 1 file && bufSize < 1line",
			args: args{
				filePaths:      []string{fileSize198Bytes3Lines32BytesAnd3Lines34BytesPerLine},
				numLine:        3,
				defaultBufSize: 22,
			},
			wantOut: []string{
				"10,20Lines_32And34BytesPerLine,10\n",
				"11,20Lines_32And34BytesPerLine,11\n",
				"12,20Lines_32And34BytesPerLine,12\n",
			},
			wantErr: false,
		},
		// 2. Tail whole single file
		{
			name: "tail whole one file && tailing lines > number lines of file",
			args: args{
				filePaths:      []string{fileSize198Bytes3Lines32BytesAnd3Lines34BytesPerLine},
				numLine:        30,
				defaultBufSize: 22,
			},
			wantOut: []string{
				"0,20Lines_32And34BytesPerLine,0\n",
				"1,20Lines_32And34BytesPerLine,1\n",
				"2,20Lines_32And34BytesPerLine,2\n",
				"10,20Lines_32And34BytesPerLine,10\n",
				"11,20Lines_32And34BytesPerLine,11\n",
				"12,20Lines_32And34BytesPerLine,12\n",
			},
			wantErr: false,
		},
		{
			name: "tail whole one file && tailing lines == number lines of file",
			args: args{
				filePaths:      []string{fileSize198Bytes3Lines32BytesAnd3Lines34BytesPerLine},
				numLine:        6,
				defaultBufSize: 22,
			},
			wantOut: []string{
				"0,20Lines_32And34BytesPerLine,0\n",
				"1,20Lines_32And34BytesPerLine,1\n",
				"2,20Lines_32And34BytesPerLine,2\n",
				"10,20Lines_32And34BytesPerLine,10\n",
				"11,20Lines_32And34BytesPerLine,11\n",
				"12,20Lines_32And34BytesPerLine,12\n",
			},
			wantErr: false,
		},
		// 3. Tail whole multiple files
		{
			name: "tail whole multiple files",
			args: args{
				filePaths: []string{
					fileSize198Bytes3Lines32BytesAnd3Lines34BytesPerLine,
					fileSize135Bytes5Lines27BytesPerLine,
					emptyFile,
				},
				numLine:        100,
				defaultBufSize: 22,
			},
			wantOut: []string{
				"==> ./test_data/FileSize198Bytes_3Lines32Bytes_And3Lines34BytesPerLine.txt <==\n",
				"0,20Lines_32And34BytesPerLine,0\n",
				"1,20Lines_32And34BytesPerLine,1\n",
				"2,20Lines_32And34BytesPerLine,2\n",
				"10,20Lines_32And34BytesPerLine,10\n",
				"11,20Lines_32And34BytesPerLine,11\n",
				"12,20Lines_32And34BytesPerLine,12\n",
				"\n",
				"==> ./test_data/FileSize135Bytes_5Lines_27BytesPerLine.txt <==\n",
				"0,10Lines_27BytesPerLine,0\n",
				"1,10Lines_27BytesPerLine,1\n",
				"2,10Lines_27BytesPerLine,2\n",
				"3,10Lines_27BytesPerLine,3\n",
				"4,10Lines_27BytesPerLine,4\n",
				"\n",
				"==> ./test_data/empty.txt <==\n",
			},
			wantErr: false,
		},
		// 4. Tail part of multiple files
		{
			name: "tail part of multiple files",
			args: args{
				filePaths: []string{
					fileSize198Bytes3Lines32BytesAnd3Lines34BytesPerLine,
					fileSize135Bytes5Lines27BytesPerLine,
					emptyFile,
				},
				numLine:        5,
				defaultBufSize: 22,
			},
			wantOut: []string{
				"==> ./test_data/FileSize198Bytes_3Lines32Bytes_And3Lines34BytesPerLine.txt <==\n",
				"1,20Lines_32And34BytesPerLine,1\n",
				"2,20Lines_32And34BytesPerLine,2\n",
				"10,20Lines_32And34BytesPerLine,10\n",
				"11,20Lines_32And34BytesPerLine,11\n",
				"12,20Lines_32And34BytesPerLine,12\n",
				"\n",
				"==> ./test_data/FileSize135Bytes_5Lines_27BytesPerLine.txt <==\n",
				"0,10Lines_27BytesPerLine,0\n",
				"1,10Lines_27BytesPerLine,1\n",
				"2,10Lines_27BytesPerLine,2\n",
				"3,10Lines_27BytesPerLine,3\n",
				"4,10Lines_27BytesPerLine,4\n",
				"\n",
				"==> ./test_data/empty.txt <==\n",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		out := &outTest{
			content: make([]string, 0),
		}

		t.Run(tt.name, func(t *testing.T) {
			if err := tailFiles(tt.args.filePaths, tt.args.numLine, tt.args.defaultBufSize, out); (err != nil) != tt.wantErr {
				t.Errorf("tailFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut := out.content; !sameContent(gotOut, tt.wantOut) {
				t.Errorf("tailFiles() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
