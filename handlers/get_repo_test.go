package handlers_test

import (
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/handlers"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces/mock_usecases"
)

func TestGetRepository_ServeHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repositoryID := domain.RepositoryID{Owner: "owner", Name: "repo"}

	getRepository := mock_usecases.NewMockGetRepository(ctrl)
	getRepository.EXPECT().Do(gomock.Not(nil), repositoryID).
		Return(&usecases.GetRepositoryResponse{}, nil)

	h := handlers.NewRouter(handlers.Handlers{
		GetRepository: handlers.GetRepository{GetRepository: getRepository},
	})
	r := httptest.NewRequest("GET", "/owner/repo/status", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("StatusCode wants 200 but %v", resp.StatusCode)
	}
	contentType := resp.Header.Get("content-type")
	if w := "text/html"; contentType != w {
		t.Errorf("content-type wants %s but %s", w, contentType)
	}
}

func TestGetRepository_ServeHTTP_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repositoryID := domain.RepositoryID{Owner: "owner", Name: "repo"}

	for _, c := range []struct {
		getRepositoryError error
	}{
		{
			func() error {
				err := mock_usecases.NewMockGetRepositoryError(ctrl)
				err.EXPECT().NoSuchRepository().AnyTimes().Return(true)
				err.EXPECT().NoGradleVersion().AnyTimes()
				return err
			}(),
		}, {
			func() error {
				err := mock_usecases.NewMockGetRepositoryError(ctrl)
				err.EXPECT().NoSuchRepository().AnyTimes()
				err.EXPECT().NoGradleVersion().AnyTimes().Return(true)
				return err
			}(),
		},
	} {
		getRepository := mock_usecases.NewMockGetRepository(ctrl)
		getRepository.EXPECT().Do(gomock.Not(nil), repositoryID).
			Return(nil, c.getRepositoryError)

		h := handlers.NewRouter(handlers.Handlers{
			GetRepository: handlers.GetRepository{GetRepository: getRepository},
		})
		r := httptest.NewRequest("GET", "/owner/repo/status", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)

		resp := w.Result()
		if resp.StatusCode != 404 {
			t.Errorf("StatusCode wants 404 but %v", resp.StatusCode)
		}
		contentType := resp.Header.Get("content-type")
		if w := "text/html"; contentType != w {
			t.Errorf("content-type wants %s but %s", w, contentType)
		}
	}
}
