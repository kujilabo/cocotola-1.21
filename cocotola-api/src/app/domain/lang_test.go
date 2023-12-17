package domain

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	libdomain "github.com/kujilabo/redstart/lib/domain"
)

func TestNewLang2(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		args          string
		want          Lang2
		wantErr       bool
		wantErrDetail error
	}{
		{
			name:    "en",
			args:    "en",
			want:    Lang2EN,
			wantErr: false,
		},
		{
			name:    "ja",
			args:    "ja",
			want:    Lang2JA,
			wantErr: false,
		},
		{
			name:          "empty string",
			args:          "",
			wantErr:       true,
			wantErrDetail: libdomain.ErrInvalidArgument,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewLang2(tt.args)
			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				if tt.wantErrDetail != nil && !errors.Is(err, tt.wantErrDetail) {
					t.Errorf("NewLang2() err = %v, wantErrDetail %v", err, tt.wantErrDetail)
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLang2() = %v, want %v", got, tt.want)
			}
		})
	}
}
