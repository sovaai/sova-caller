package constants

import "os"

var Enviroment = map[string]string{
	"PORT": "4000",
}

func EnvVariable(key string) string {

	ENV_VAR, exists := os.LookupEnv(key)

	if !exists {
		ENV_VAR = Enviroment[key]
		return ENV_VAR
	}

	return ENV_VAR
}

const ASR_FILE_NAME = "asr.wav"
const TTS_FILE_NAME = "tts.wav"
const AUDIO_DIR = "audio"
const SESSIONS_DIR = "sessions"
const DEFAULT_SAMPLE_RATE = 16000
const DEFAULT_TTS_VOICE = "natasha"
