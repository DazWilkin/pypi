package json

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"gotest.tools/assert"
)

var c *Client

var tests = []struct {
	Name     string
	Version  string
	Filename string
	URL      string
}{
	{
		Name:     "google-api-core",
		Version:  "1.14.2",
		Filename: "google_api_core-1.14.2-py2.py3-none-any.whl",
		URL:      "https://files.pythonhosted.org/packages/71/e5/7059475b3013a3c75abe35015c5761735ab224eb1b129fee7c8e376e7805/google_api_core-1.14.2-py2.py3-none-any.whl",
	},
	{
		Name:     "grpcio",
		Version:  "1.23.0",
		Filename: "grpcio-1.23.0-cp37-cp37m-manylinux1_x86_64.whl",
		URL:      "https://files.pythonhosted.org/packages/e5/27/1f908ebb99c8d48a5ba4eb9d7997f5633b920d98fe712f67aaa0663f1307/grpcio-1.23.0-cp37-cp37m-manylinux1_x86_64.whl",
	},
	{
		Name:     "opencensus",
		Version:  "0.7.2",
		Filename: "opencensus-0.7.2-py2.py3-none-any.whl",
		URL:      "https://files.pythonhosted.org/packages/b6/13/c37904d9f77c320dfd39b6d5c6bd95947c59e19baa3c6c4f4df6418d1433/opencensus-0.7.2-py2.py3-none-any.whl"},
	{
		Name:     "prometheus_client",
		Version:  "0.7.1",
		Filename: "prometheus_client-0.7.1.tar.gz",
		URL:      "https://files.pythonhosted.org/packages/b3/23/41a5a24b502d35a4ad50a5bb7202a5e1d9a0364d0c12f56db3dbf7aca76d/prometheus_client-0.7.1.tar.gz"},
	{
		Name:     "protobuf",
		Version:  "3.9.1",
		Filename: "protobuf-3.9.1.tar.gz",
		URL:      "https://files.pythonhosted.org/packages/6d/54/12c5c92ffab546538ea5b544c6afbfcce333fd47e99c1198e24a8efdef1f/protobuf-3.9.1.tar.gz"},
}

func TestMain(m *testing.M) {
	c = NewClient(&http.Client{})
	os.Exit(m.Run())
}
func TestPackage(t *testing.T) {
	for _, test := range tests {
		var err error
		var r Response
		var pp Packages

		want := Package{
			Filename: test.Filename,
			URL:      test.URL,
		}

		t.Run("Release", func(t *testing.T) {
			r, err = c.Release(test.Name, test.Version)
			if err != nil {
				t.Error(err)
			}
		})
		t.Run("Packages", func(t *testing.T) {
			pp, err = r.Packages(test.Version)
			if err != nil {
				t.Error(err)
			}
		})
		t.Run("Package{Filename:\"\",URL:\"\"}", func(t *testing.T) {
			_, err := pp.Package(Package{
				Filename: "",
				URL:      "",
			})
			assert.Error(t, err, "undefined criteria will never match")
		})
		t.Run(fmt.Sprintf("Package{Filename:\"%s\",URL:\"\"}", test.Filename), func(t *testing.T) {
			got, err := pp.Package(Package{
				Filename: test.Filename,
				URL:      "",
			})
			if err != nil {
				t.Error(err)
			}
			assert.DeepEqual(t, got, want)
		})
		t.Run(fmt.Sprintf("Package{Filename:\"\",URL:\"%s\"}", test.URL), func(t *testing.T) {
			got, err := pp.Package(Package{
				Filename: "",
				URL:      test.URL,
			})
			if err != nil {
				t.Error(err)
			}
			assert.DeepEqual(t, got, want)
		})
		t.Run(fmt.Sprintf("Package{Filename:\"%s\",URL:\"%s\"}", test.Filename, test.URL), func(t *testing.T) {
			got, err := pp.Package(Package{
				Filename: test.Filename,
				URL:      test.URL,
			})
			if err != nil {
				t.Error(err)
			}
			assert.DeepEqual(t, got, want)
		})
	}
}
