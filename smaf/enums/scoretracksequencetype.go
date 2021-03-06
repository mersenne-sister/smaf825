package enums

import (
	"encoding/json"
	"fmt"
)

type ScoreTrackSequenceType int

const (
	// Sequence Data は1つの連続したシーケンス・データである。Seek Point や Phrase List はシーケンス中の意味のある位置を外部から参照する目的で利用する。
	ScoreTrackSequenceType_StreamSequence ScoreTrackSequenceType = iota
	// Sequence Data は複数のフレーズデータを連続で表記したものである。Phrase List は外部から個別フレーズを認識する為に用いる。
	ScoreTrackSequenceType_Subsequence
)

func (t ScoreTrackSequenceType) String() string {
	s := "undefined"
	switch t {
	case ScoreTrackSequenceType_StreamSequence:
		s = "StreamSequence"
	case ScoreTrackSequenceType_Subsequence:
		s = "Subsequence"
	}
	return fmt.Sprintf("%s(0x%02X)", s, int(t))
}

func (t ScoreTrackSequenceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}
