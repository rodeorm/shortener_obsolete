package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReturnShortKey(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		want    int
		wantErr bool
	}{
		{
			name:    "Контроль ошибок. Отрицательные значения",
			value:   -1,
			wantErr: true,
		},
		{
			name:    "Контроль ошибок. Ноль",
			value:   0,
			wantErr: true,
		},
		{
			name:    "Проверка длины ключа 5",
			value:   5,
			want:    5,
			wantErr: false,
		}, {
			name:    "Проверка длины ключа 10",
			value:   10,
			want:    10,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReturnShortKey(tt.value)
			if !tt.wantErr {
				require.NoError(t, err)
				assert.Equal(t, len(got), tt.want)
				return
			}
			assert.Error(t, err)
		})
	}
}
