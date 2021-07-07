package audio

import "encoding/binary"

var (
	tokenRiff       = [4]byte{'R', 'I', 'F', 'F'}
	tokenWaveFormat = [4]byte{'W', 'A', 'V', 'E'}
	tokenChunkFmt   = [4]byte{'f', 'm', 't', ' '}
	tokenData       = [4]byte{'d', 'a', 't', 'a'}
)

// EncodeToWav creates a pcm encoded wav file from raw pcm bytes
func EncodeToWav(rawPcmBytes []byte, samplingRate uint32, bitDepth uint16, numChannels uint16) ([]byte, error) {
	numDataLength := len(rawPcmBytes)
	// numFrames := numDataLength / 2

	// append wav header of 44 bytes
	waveAudioBytes := make([]byte, 44 + numDataLength)
	// write RIFF and size
	copy(waveAudioBytes, tokenRiff[:])
	binary.LittleEndian.PutUint32(waveAudioBytes[4:], uint32(36 + numDataLength))

	// write format
	copy(waveAudioBytes[8:], tokenWaveFormat[:])

	// write fmt chunk
	copy(waveAudioBytes[12:], tokenChunkFmt[:])
	binary.LittleEndian.PutUint32(waveAudioBytes[16:], 16)
	binary.LittleEndian.PutUint16(waveAudioBytes[20:], 1)
	binary.LittleEndian.PutUint16(waveAudioBytes[22:], numChannels)
	binary.LittleEndian.PutUint32(waveAudioBytes[24:], samplingRate)
	binary.LittleEndian.PutUint32(waveAudioBytes[28:], uint32(numChannels) * samplingRate * uint32(bitDepth) / 8) // bytes per sec
	binary.LittleEndian.PutUint16(waveAudioBytes[32:], (bitDepth / 8) * numChannels) // bytes per block
	binary.LittleEndian.PutUint16(waveAudioBytes[34:], bitDepth) // bytes per block

	// write data chunk
	copy(waveAudioBytes[36:], tokenData[:])
	binary.LittleEndian.PutUint32(waveAudioBytes[40:], uint32(36 + numDataLength))
	copy(waveAudioBytes[44:], rawPcmBytes)

	return waveAudioBytes, nil
}
