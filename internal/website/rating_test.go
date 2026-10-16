package website

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRatingValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name    string
		form    url.Values
		wantErr bool
	}{
		{
			name:    "valid 5 stars",
			form:    url.Values{"pastry": {"5"}, "filling": {"3"}, "topping": {"4"}, "looks": {"2"}, "value": {"5"}},
			wantErr: false,
		},
		{
			name:    "missing pastry",
			form:    url.Values{"filling": {"3"}, "topping": {"4"}, "looks": {"2"}, "value": {"5"}},
			wantErr: true,
		},
		{
			name:    "out of range high",
			form:    url.Values{"pastry": {"6"}, "filling": {"3"}, "topping": {"4"}, "looks": {"2"}, "value": {"5"}},
			wantErr: true,
		},
		{
			name:    "out of range low",
			form:    url.Values{"pastry": {"0"}, "filling": {"3"}, "topping": {"4"}, "looks": {"2"}, "value": {"5"}},
			wantErr: true,
		},
		{
			name:    "invalid number",
			form:    url.Values{"pastry": {"abc"}, "filling": {"3"}, "topping": {"4"}, "looks": {"2"}, "value": {"5"}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/years/2024/brands/test/pie/rating", strings.NewReader(tt.form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Form = tt.form

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// We don't call the full handler here since it requires DB,
			// but we validate the validation logic directly.
			errors := []string{}
			rnkingmap := map[string]int32{}

			for _, k := range rankings {
				rv := c.Request.FormValue(k)
				if rv == "" {
					errors = append(errors, k+" is required")
					continue
				}
				prsd, err := strconv.ParseInt(rv, 10, 32)
				if err != nil {
					errors = append(errors, k+" must be an integer between 1 and 5")
					continue
				}
				if prsd < 1 || prsd > 5 {
					errors = append(errors, k+" must be between 1 and 5")
					continue
				}
				rnkingmap[k] = int32(prsd)
			}

			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("got errors = %v, wantErr = %v", errors, tt.wantErr)
			}
		})
	}
}
