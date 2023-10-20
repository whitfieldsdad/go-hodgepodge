package main

import (
	"bytes"
	"context"
	"net/url"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"github.com/pkg/errors"
)

func PlayMP3File(path string, ctx context.Context) error {
	return PlayMP3FileWithCallback(path, ctx, nil, nil)
}

func PlayMP3FileWithCallback(
	path string,
	ctx context.Context,
	startCallback chan interface{},
	stopCallback chan interface{}) error {

	u, err := url.ParseRequestURI(path)
	if err == nil && u.Scheme != "" {
		log.Infof("Playing MP3 from remote location: %s\n", path)
		return playRemoteMP3File(path, ctx, startCallback, stopCallback)
	} else {
		log.Infof("Playing MP3 from local location: %s\n", path)
		return playLocalMP3File(path, ctx, startCallback, stopCallback)
	}
}

func playLocalMP3File(
	path string,
	ctx context.Context,
	startCallback chan interface{},
	stopCallback chan interface{}) error {

	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return PlayMP3BytesWithCallback(b, ctx, startCallback, stopCallback)
}

func playRemoteMP3File(
	path string,
	ctx context.Context,
	startCallback chan interface{},
	stopCallback chan interface{}) error {

	b, err := DownloadFileIntoMemory(path)
	if err != nil {
		return errors.Wrap(err, "download failed")
	}
	err = PlayMP3BytesWithCallback(b, ctx, startCallback, stopCallback)
	if err != nil {
		return errors.Wrap(err, "failed to play MP3")
	}
	return nil
}

func PlayMP3Bytes(fileBytes []byte, ctx context.Context) error {
	return PlayMP3BytesWithCallback(fileBytes, ctx, nil, nil)
}

func PlayMP3BytesWithCallback(
	fileBytes []byte,
	ctx context.Context,
	startCallback chan interface{},
	stopCallback chan interface{}) error {

	// Decode the MP3 file.
	fileBytesReader := bytes.NewReader(fileBytes)
	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		return errors.Wrap(err, "MP3 decode failed")
	}
	samplingRate := 44100
	numOfChannels := 2
	audioBitDepth := 2

	otoCtx, readyChan, err := oto.NewContext(samplingRate, numOfChannels, audioBitDepth)
	if err != nil {
		return errors.Wrap(err, "failed to build context for MP3 player")
	}
	<-readyChan

	// Play the MP3 file.
	player := otoCtx.NewPlayer(decodedMp3)
	player.Play()
	defer player.Close()

	// Fire an optional callback as soon as the sound starts.
	if startCallback != nil {
		startCallback <- nil
	}

	done := false
	for player.IsPlaying() && !done {
		select {
		case <-ctx.Done():
			done = true
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
	log.Infof("A sound has finished playing")
	if stopCallback != nil {
		stopCallback <- nil
	}

	// Close the MP3 player.
	err = player.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close MP3 player")
	}
	return nil
}
