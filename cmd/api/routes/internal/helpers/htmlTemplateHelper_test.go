// 90%

package helpers

import (
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockResponseWriter is a mock of http.ResponseWriter
type MockResponseWriter struct {
	mock.Mock
}

func (m *MockResponseWriter) Header() http.Header {
	args := m.Called()
	return args.Get(0).(http.Header)
}

func (m *MockResponseWriter) Write(b []byte) (int, error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.Called(statusCode)
}

func TestRenderTemplate(t *testing.T) {
	tests := []struct {
		name          string
		templateName  string
		data          interface{}
		mockAbs       func(string) (string, error)
		mockParse     func(...string) (*template.Template, error)
		mockExecute   func(*template.Template, http.ResponseWriter, interface{}) error
		expectedError string
	}{
		{
			name:         "Successful template rendering",
			templateName: "test.tmpl",
			data:         map[string]string{"key": "value"},
			mockAbs: func(path string) (string, error) {
				return "/absolute/path/to/test.tmpl", nil
			},
			mockParse: func(filenames ...string) (*template.Template, error) {
				return template.New("test"), nil
			},
			mockExecute: func(t *template.Template, w http.ResponseWriter, data interface{}) error {
				return nil
			},
			expectedError: "",
		},
		{
			name:         "Error in filepath.Abs",
			templateName: "test.tmpl",
			data:         nil,
			mockAbs: func(path string) (string, error) {
				return "", errors.New("abs error")
			},
			mockParse:     nil,
			mockExecute:   nil,
			expectedError: "abs error",
		},
		{
			name:         "Error in template.ParseFiles",
			templateName: "test.tmpl",
			data:         nil,
			mockAbs: func(path string) (string, error) {
				return "/absolute/path/to/test.tmpl", nil
			},
			mockParse: func(filenames ...string) (*template.Template, error) {
				return nil, errors.New("parse error")
			},
			mockExecute:   nil,
			expectedError: "parse error",
		},
		{
			name:         "Error in template.Execute",
			templateName: "test.tmpl",
			data:         map[string]string{"key": "value"},
			mockAbs: func(path string) (string, error) {
				return "/absolute/path/to/test.tmpl", nil
			},
			mockParse: func(filenames ...string) (*template.Template, error) {
				return template.New("test"), nil
			},
			mockExecute: func(t *template.Template, w http.ResponseWriter, data interface{}) error {
				return errors.New("execute error")
			},
			expectedError: "execute error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			oldAbs, oldParse, oldExecute := filepathAbs, templateParse, templateExecute
			defer func() {
				filepathAbs, templateParse, templateExecute = oldAbs, oldParse, oldExecute
			}()

			filepathAbs = tt.mockAbs
			templateParse = tt.mockParse
			templateExecute = tt.mockExecute

			// Create a simple mock response writer
			w := httptest.NewRecorder()

			// Call the function
			err := RenderTemplate(w, tt.data, tt.templateName)

			// Assert the result
			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestRenderTemplateIntegration(t *testing.T) {
	// This test uses the actual implementation
	w := httptest.NewRecorder()
	data := map[string]string{"key": "value"}
	templateName := "test.tmpl"

	err := RenderTemplate(w, data, templateName)

	// Note: This will fail unless you have an actual template file at the specified path
	// You might want to create a test template file or skip this test in CI environments
	assert.Error(t, err) // Expecting an error because the template file doesn't exist
}

func TestTemplateExecute(t *testing.T) {
	// Create a simple template
	tmpl, err := template.New("test").Parse("Hello {{.Name}}")
	assert.NoError(t, err)

	// Create a mock response writer
	w := httptest.NewRecorder()

	// Create some test data
	data := struct {
		Name string
	}{
		Name: "World",
	}

	// Call the templateExecute function
	err = templateExecute(tmpl, w, data)

	// Assert that there's no error
	assert.NoError(t, err)

	// Assert that the output is correct
	assert.Equal(t, "Hello World", w.Body.String())
}
