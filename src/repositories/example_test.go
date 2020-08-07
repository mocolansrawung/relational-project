package repositories

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/evermos/boilerplate-go/infras"
)

func TestExampleRepository_Get(t *testing.T) {
	tests := []struct {
		name          string
		configureMock func(sqlmock.Sqlmock)
		want          string
		wantErr       bool
	}{
		{
			name: "return good",
			configureMock: func(s sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"status"}).AddRow("good")
				s.ExpectQuery("SELECT 'good' as status").WillReturnRows(rows)
			},
			want:    "good",
			wantErr: false,
		},
		{
			name: "must return error",
			configureMock: func(s sqlmock.Sqlmock) {
				s.ExpectQuery("SELECT 'good' as status").WillReturnError(errors.New("invalid query"))
			},
			want:    "",
			wantErr: true,
		},
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ExampleRepository{
				Db: infras.OpenMock(db),
			}
			tt.configureMock(mock)
			got, err := r.Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
