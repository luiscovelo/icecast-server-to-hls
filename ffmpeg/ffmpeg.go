package ffmpeg

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"
)

func FFmpeg(ctx context.Context, reader *io.PipeReader) {
	args := []string{
		"-loglevel", "debug",
		"-i", "pipe:0",
		"-c:a", "aac",
		"-b:a", "128k",
		"-f", "hls",
		"-hls_time", "5",
		"-hls_list_size", "6",
		"-hls_start_number_source", "epoch",
		"-hls_segment_filename", "./output/segment_%v_%03d.ts",
		"-master_pl_name", "master.m3u8",
		"./output/playlist_%v.m3u8",
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	cmd.Stdin = reader
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Println("ffmpeg error", err)
	}
}
