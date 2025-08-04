package logger

import (
	"testing"
)

func BenchmarkZerologInfo(b *testing.B) {
	log := New("../../log", "logs.log")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info().Msg("This is an info message")
	}
}

func BenchmarkZerologError(b *testing.B) {
	log := New("../../log", "logs.log")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Error().Msg("This is an error message")
	}
}

func BenchmarkZerologWarn(b *testing.B) {
	log := New("../../log", "logs.log")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Warn().Msg("This is a warning message")
	}
}

func BenchmarkZerologFatal(b *testing.B) {
	// Fatal will call os.Exit(1), so we can't benchmark it directly.
	// Instead, we benchmark the construction of the event up to, but not including, .Msg()
	log := New("../../log", "logs.log")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = log.Fatal() // Don't call .Msg() to avoid exiting
	}
}
