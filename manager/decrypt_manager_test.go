package manager

import (
	"testing"
	"time"
)

func TestDecodeTimeKeySuccess(t *testing.T) {
	sn := "SN-TEST-0001"
	head := []byte{
		'A', 'A', // reserved
		'F', 'D', // second: 35
		'G', 'C', // minute: 26
		'I', 'A', // hour: 08
		'B', 'C', // day: 12
		'C', 'A', // month: 02
		'G', 'C', // year: 26
	}

	xorPayload := make([]byte, 14)
	for i := 0; i < 14; i++ {
		xorPayload[i] = head[i] ^ sn[i%len(sn)]
	}
	tail, err := calcUnlockKey(xorPayload, 0, 14)
	if err != nil {
		t.Fatalf("calc unlock key failed: %v", err)
	}

	key := string(append(head, tail...))
	result, err := DecryptManagerGlobal.DecodeTimeKey(sn, key)
	if err != nil {
		t.Fatalf("decode should succeed: %v", err)
	}
	if !result.Valid {
		t.Fatalf("expected valid result")
	}
	if result.DecodedYear != 26 || result.DecodedMonth != 2 || result.DecodedDay != 21 {
		t.Fatalf("unexpected date: %+v", result)
	}
	if result.DecodedHour != 8 || result.DecodedMinute != 26 || result.DecodedSecond != 35 {
		t.Fatalf("unexpected time: %+v", result)
	}
}

func TestDecodeTimeKeyInvalidSignature(t *testing.T) {
	sn := "SN-TEST-0002"
	key := "AAAAAAAAAAAAAABBBBBBBBBBBBBB"

	_, err := DecryptManagerGlobal.DecodeTimeKey(sn, key)
	if err == nil {
		t.Fatalf("expected signature verification to fail")
	}
}

func TestGenerateTimeKeyRoundTrip(t *testing.T) {
	sn := "SN-TEST-0009"
	target := time.Date(2026, time.March, 22, 13, 45, 36, 0, time.Local)

	gen, err := DecryptManagerGlobal.GenerateTimeKey(sn, target)
	if err != nil {
		t.Fatalf("generate should succeed: %v", err)
	}
	if len(gen.Key) != 28 {
		t.Fatalf("unexpected key length: %d", len(gen.Key))
	}

	decoded, err := DecryptManagerGlobal.DecodeTimeKey(sn, gen.Key)
	if err != nil {
		t.Fatalf("decode generated key should succeed: %v", err)
	}
	if decoded.FullYear != target.Year() ||
		decoded.DecodedMonth != int(target.Month()) ||
		decoded.DecodedDay != target.Day() ||
		decoded.DecodedHour != target.Hour() ||
		decoded.DecodedMinute != target.Minute() ||
		decoded.DecodedSecond != target.Second() {
		t.Fatalf("decoded value mismatch: got=%+v target=%v", decoded, target)
	}
}
