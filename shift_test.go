package kowalski

import (
	"reflect"
	"testing"
)

func TestCaesarShift(t *testing.T) {
	tests := []struct {
		name  string
		query string
		count int
		want  string
	}{
		{"positive", "Btusb jodmjobou, tfe opo pcmjhbou", 1, "Cuvtc kpenkpcpv, ugf pqp qdnkicpv"},
		{"26", "Btusb jodmjobou, tfe opo pcmjhbou", 26, "Btusb jodmjobou, tfe opo pcmjhbou"},
		{"40", "Btusb jodmjobou, tfe opo pcmjhbou", 40, "Phigp xcraxcpci, hts cdc dqaxvpci"},
		{"negative", "Btusb jodmjobou, tfe opo pcmjhbou", -5, "Wopnw ejyhejwjp, oaz jkj kxhecwjp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CaesarShift(tt.query, uint8(tt.count)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CaesarShift() = %v, want %v", got, tt.want)
			}
		})
	}
}
